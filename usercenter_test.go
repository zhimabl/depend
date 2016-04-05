package depend_test

import (
	"testing"

	"github.com/zhimabl/depend"
)

func init() {
	depend.LoadServicesConf()
}

func TestRecordFirstOrderId(t *testing.T) {
	// err := depend.RecordFirstOrderId(nil, 6, 76, 118)
	// if err != nil {
	// 	t.Errorf("record first order id error:%#v", err)
	// 	return
	// }

}
