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

func fromName() string {
	if isPro {
		return "线上服务报警"
	}

	return "开发/测试服务报警"
}
