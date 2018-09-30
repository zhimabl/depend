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

func TestSendDingtalkBotMsg(t *testing.T) {
	webhook := "https://oapi.dingtalk.com/robot/send?access_token=901b25b1c69135266c15d8eadfecc836bb50097036bd563de5059a9fa4ffe332"
	err := depend.SendDingtalkBotMsg(nil, webhook, "测试订单监控", "这是内容呢", "text")
	if err != nil {
		t.Errorf("send dingtalk msg by bot fail:%v", err)
	}
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
