package depend

import (
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

const (
	OpTypeAdd = iota
	OpTypeUpdate
	OpTypeDel
)

type TableSchema struct {
	Field   string `xorm:"Field"`
	Comment string `xorm:"Comment"`
}

// RecordAddOp 记录新增操作
func RecordAddOp(ctx context.Context, engine *xorm.Engine, bean interface{}, opUser string) error {
	opRecord, err := genOpRecord(engine, bean, opUser)
	if err != nil {
		logger.Errorln("RecordAddOp error:", err.Error())
		return err
	}

	return RecordOperate(ctx, OpTypeAdd, opRecord)
}

// RecordUpdateOp 记录修改操作
func RecordUpdateOp(ctx context.Context, engine *xorm.Engine, bean interface{}, changeVals map[string]interface{}, opUser string) error {
	opRecord, err := genOpRecord(engine, bean, opUser)
	if err != nil {
		logger.Errorln("RecordUpdateOp error:", err)
		return err
	}

	tableInfo := engine.TableInfo(bean)

	tableSchemas := make([]*TableSchema, 0)
	err = engine.Sql("show full fields from " + tableInfo.Name).Find(&tableSchemas)
	if err != nil {
		logger.Errorln("RecordUpdateOp show full fields error:", err)
		return err
	}

	tableSchemaMap := make(map[string]string, len(tableSchemas))
	for _, tableSchema := range tableSchemas {
		tableSchemaMap[tableSchema.Field] = tableSchema.Comment
	}

	length := len(changeVals)
	fields := make([]string, 0, length)
	fieldNames := make([]string, 0, length)
	oldValues := make([]interface{}, 0, length)
	newValues := make([]interface{}, 0, length)

	s := structs.New(bean)

	for field, newVal := range changeVals {
		var oldVal interface{}
		sField, ok := s.FieldOk(goutils.CamelName(field))
		if ok {
			oldVal = sField.Value()
		} else {
			sFields := s.Fields()
			for _, f := range sFields {
				if f.Tag("xorm") == field {
					oldVal = f.Value()
					break
				}
			}
		}

		// 排除时间的变更
		if _, ok := oldVal.(time.Time); ok {
			continue
		}
		if equal(oldVal, newVal) {
			continue
		}

		fields = append(fields, field)

		comment := tableSchemaMap[field]
		if comment == "" {
			comment = field
		}
		fieldNames = append(fieldNames, comment)

		oldValues = append(oldValues, oldVal)
		newValues = append(newValues, newVal)
	}

	if len(fields) == 0 {
		logger.Infoln("table_name", tableInfo.Name, "not any change!")
		return nil
	}

	opRecord["fields"] = strings.Join(fields, ",")
	buf, err := json.Marshal(fieldNames)
	if err != nil {
		logger.Errorln("RecordUpdateOp json encoding fields error:", err)
		return err
	}
	opRecord["field_names"] = string(buf)

	buf, err = json.Marshal(oldValues)
	if err != nil {
		logger.Errorln("RecordUpdateOp json encoding old_values error:", err)
		return err
	}
	opRecord["old_values"] = string(buf)

	buf, err = json.Marshal(newValues)
	if err != nil {
		logger.Errorln("RecordUpdateOp json encoding new_values error:", err)
		return err
	}
	opRecord["new_values"] = string(buf)

	return RecordOperate(ctx, OpTypeUpdate, opRecord)
}

// RecordDelOp 记录删除操作
func RecordDelOp(ctx context.Context, engine *xorm.Engine, bean interface{}, opUser string) error {
	opRecord, err := genOpRecord(engine, bean, opUser)
	if err != nil {
		logger.Errorln("RecordDelOp error:", err.Error())
		return err
	}

	return RecordOperate(ctx, OpTypeDel, opRecord)
}

// RecordOperate 记录操作
func RecordOperate(ctx context.Context, opType int, opRecord map[string]interface{}) error {
	opRecord["op_type"] = opType

	buf, err := json.Marshal(opRecord)
	if err != nil {
		ProcessNotInterruptErr(ctx, goutils.ConvertString(opRecord["pri_key"]), opRecord["pri_val"], err)
		return err
	}
	data := url.Values{
		"data":  {string(buf)},
		"async": {"true"},
	}
	_, err = callService(opRecordService, "POST", "/op_record", data)
	if err != nil {
		return err
	}

	return nil
}

func genOpRecord(engine *xorm.Engine, bean interface{}, opUser string) (map[string]interface{}, error) {
	val := reflect.ValueOf(bean)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, errors.New("bean must be struct or pointer of struct!")
	}

	tableInfo := engine.TableInfo(bean)
	priKeys := tableInfo.PrimaryKeys
	if len(priKeys) != 1 {
		return nil, errors.New("no pk or composite pk, shoudle call RecordOperate!")
	}
	priKey := priKeys[0]

	fieldName := tableInfo.GetColumn(priKey).FieldName
	var priVal interface{}
	priValue := val.FieldByName(fieldName)
	switch priValue.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		priVal = priValue.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
		priVal = priValue.Uint()
	default:
		return nil, errors.New("pk value must be integer!")
	}

	opRecord := map[string]interface{}{
		"table_name": tableInfo.Name,
		"pri_key":    priKey,
		"pri_val":    priVal,
		"op_user":    opUser,
	}

	return opRecord, nil
}

func equal(x, y interface{}) bool {
	if x == y {
		return true
	}

	if reflect.DeepEqual(x, y) {
		return true
	}

	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("equal panic:", err)
		}
	}()

	xVal := reflect.ValueOf(x)
	switch xVal.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		var y1 int64
		x1 := xVal.Int()
		if strY, ok := y.(string); ok {
			y1 = goutils.MustInt64(strY)
		} else if fY, ok := y.(float64); ok {
			y1 = int64(fY)
		} else {
			y1 = reflect.ValueOf(y).Int()
		}

		return x1 == y1
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
		var y1 uint64
		x1 := xVal.Uint()
		if strY, ok := y.(string); ok {
			y1 = uint64(goutils.MustInt64(strY))
		} else if fY, ok := y.(float64); ok {
			y1 = uint64(fY)
		} else {
			y1 = reflect.ValueOf(y).Uint()
		}

		return x1 == y1
	case reflect.Float32, reflect.Float64:
		var y1 float64
		x1 := xVal.Float()
		if strY, ok := y.(string); ok {
			y1 = goutils.MustFloat(strY)
		} else {
			y1 = reflect.ValueOf(y).Float()
		}

		return x1 == y1
	}

	return false
}
