package depend

import (
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strings"

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

	beanVal := reflect.ValueOf(bean)
	if beanVal.Kind() == reflect.Ptr {
		beanVal = beanVal.Elem()
	}

	length := len(changeVals)
	fields := make([]string, length)
	fieldNames := make([]string, length)
	oldValues := make([]interface{}, length)
	newValues := make([]interface{}, length)

	i := 0
	for field, val := range changeVals {
		fields[i] = field

		comment := tableSchemaMap[field]
		if comment == "" {
			comment = field
		}
		fieldNames[i] = comment

		oldValues[i] = beanVal.FieldByName(field)
		newValues[i] = val

		i++
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
	priVal := val.FieldByName(fieldName)

	opRecord := map[string]interface{}{
		"talbe_name": tableInfo.Name,
		"pri_key":    priKey,
		"pri_val":    priVal,
		"op_user":    opUser,
	}

	return opRecord, nil
}
