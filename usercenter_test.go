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

func TestReadUserDevice(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{6, "71ef26febdfea41140b88008d3d257e9"},
		{36514, ""},
	}

	for _, tt := range tests {
		result := depend.ReadUserDevice(nil, tt.input)
		if result != nil {
			if actual := result.Get("client_id").MustString(); actual != tt.expected {
				t.Errorf("unexport(%d) = %q, want %q", tt.input, actual, tt.expected)
			}
		} else {
			t.Errorf("unexport(%d) = nil want %s", tt.input, tt.expected)
		}
	}
}
