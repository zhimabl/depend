package depend

import (
	"errors"
	"net/url"

	"golang.org/x/net/context"
)

func SendEmail(ctx context.Context, subject, content, tos string) error {
	data := url.Values{
		"subject":   {subject},
		"content":   {content},
		"from_name": {fromName()},
		"tos":       {tos},
	}
	_, err := callService(senderService, "POST", "/email/send", data)
	if err != nil {
		return err
	}

	return nil
}

// ServiceWarning 技术服务报警
func ServiceWarning(ctx context.Context, emailInfo map[string]string, smsContent string) error {
	data := url.Values{
		"subject":   {emailInfo["subject"] + "-服务【" + from + "】"},
		"content":   {emailInfo["content"]},
		"from_name": {fromName()},
		"gename":    {"tech_services"},
	}
	if smsContent != "" {
		data.Set("send_sms", "1")
		data.Set("sms_content", smsContent)
	}
	_, err := callService(senderService, "POST", "/email/group", data)
	if err != nil {
		return err
	}

	return nil
}

// SendSms 给用户发送短信，smsTypes 默认发送通知。支持的值：00-验证码;01-通知;02-营销
func SendSms(ctx context.Context, mobile, content string, smsTypes ...string) error {
	if mobile == "" || content == "" {
		return errors.New("mobile or content is empty!")
	}

	data := url.Values{
		"mobile":  {mobile},
		"content": {content},
	}

	if len(smsTypes) > 0 {
		data.Set("sms_type", smsTypes[0])
	}

	_, err := callService(senderService, "POST", "/sms/unicast", data)
	if err != nil {
		return err
	}

	return nil
}

// SendDingtalkMsg 给员工发送钉钉消息
func SendDingtalkMsg(ctx context.Context, userid, agentid, content string, msgtypes ...string) error {
	if userid == "" || agentid == "" || content == "" {
		return errors.New("userid or agentid or content is empty!")
	}

	data := url.Values{
		"touser":  {userid},
		"agentid": {agentid},
		"content": {content},
	}

	if len(msgtypes) > 0 {
		data.Set("msgtype", msgtypes[0])
	}

	_, err := callService(senderService, "POST", "/dingtalk/send", data)
	if err != nil {
		return err
	}

	return nil
}

// SendAppMsg App push
// destType: driver 司机端；supply 供货宝；shopper 芝麻掌柜
func SendAppMsg(ctx context.Context, userid, content, iosMsg, destType string) error {
	if userid == "" || content == "" {
		return errors.New("userid or content or content is empty!")
	}

	data := url.Values{
		"userid":    {userid},
		"content":   {content},
		"ios_msg":   {iosMsg},
		"dest_type": {destType},
	}

	_, err := callService(senderService, "POST", "/dingtalk/send", data)
	if err != nil {
		return err
	}

	return nil
}

func fromName() string {
	if isPro {
		return "线上服务报警"
	}

	return "开发/测试服务报警"
}
