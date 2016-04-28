package depend_test

import (
	"testing"

	"github.com/zhimabl/depend"
)

func init() {
	depend.LoadServicesConf()
}

func TestRecordUpdateOp(t *testing.T) {
	// bean := &model.StoreGoodsClass{}
	// db.MasterDB.Id(29831).Get(bean)

	// change := structs.New(bean).Map()
	// change["StcName"] = "购2"
	// change["StcSort"] = 99

	// depend.RecordUpdateOp(context.Background(), db.MasterDB, bean, change, "xuxinhua")
	// t.Fatal("error")
}
