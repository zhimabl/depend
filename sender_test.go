package depend_test

import (
	"testing"

	"github.com/zhimabl/depend"
)

func init() {
	depend.LoadServicesConf()
}

func TestSendEmail(t *testing.T) {
	// subject, content, tos := "usercenter subject", "content", "xuxinhua@zhimadj.com"
	// err := depend.SendEmail(nil, subject, content, tos)
	// if err != nil {
	// 	t.Errorf("send to %q fail:%v", tos, err)
	// }
}

func TestServiceWarning(t *testing.T) {
	// emailInfo := map[string]string{
	// 	"subject": "user center group test",
	// 	"content": "test",
	// }
	// smsContent := ""
	// err := depend.ServiceWarning(nil, emailInfo, smsContent)
	// if err != nil {
	// 	t.Errorf("send service warning fail:%v", err)
	// }
}
