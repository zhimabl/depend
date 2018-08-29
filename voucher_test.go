package depend_test

import (
	"testing"

	"github.com/zhimabl/depend"
)

func init() {
	depend.LoadServicesConf()
}

func TestReturnVoucher(t *testing.T) {
	err := depend.ReturnVoucher(nil, 1, 1)
	if err != nil {
		t.Errorf("return voucher error:%#v", err)
		return
	}
}
