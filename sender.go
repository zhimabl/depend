package depend

import (
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

func fromName() string {
	if isPro {
		return "线上服务报警"
	}

	return "开发/测试服务报警"
}
