package depend

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
)

func SendEmail(ctx context.Context, subject, content, tos string) error {
	senderConf := randServiceConf(senderService)
	if senderConf == nil {
		logger.Errorln(senderService, "config is empty")
		return errors.New("sender service config is empty")
	}

	httpClient := http.Client{
		Timeout: 60 * time.Second,
	}

	data := url.Values{
		"subject":   {subject},
		"content":   {content},
		"from_name": {fromName()},
		"tos":       {tos},
		"timestamp": {strconv.FormatInt(time.Now().Unix(), 10)},
		"from":      {from},
	}

	data.Set("sign", goutils.GenSign(data, getServiceSecret(senderService)))

	apiUrl := "http://" + senderConf.httpAddr + "/email/send"
	resp, err := httpClient.PostForm(apiUrl, data)
	if err != nil {
		logger.Errorf("url:%q, data:%v, error: %s", apiUrl, data, err)
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ServiceWarning 技术服务报警
func ServiceWarning(ctx context.Context, emailInfo map[string]string, smsContent string) error {
	senderConf := randServiceConf(senderService)
	if senderConf == nil {
		logger.Errorln(senderService, "config is empty")
		return errors.New("sender service config is empty")
	}

	httpClient := http.Client{
		Timeout: 60 * time.Second,
	}

	data := url.Values{
		"subject":   {emailInfo["subject"] + "-服务【" + from + "】"},
		"content":   {emailInfo["content"]},
		"from_name": {fromName()},
		"gename":    {"tech_services"},
		"timestamp": {strconv.FormatInt(time.Now().Unix(), 10)},
		"from":      {from},
	}

	if smsContent != "" {
		data.Set("send_sms", "1")
		data.Set("sms_content", smsContent)
	}

	data.Set("sign", goutils.GenSign(data, getServiceSecret(senderService)))

	apiUrl := "http://" + senderConf.httpAddr + "/email/group"
	resp, err := httpClient.PostForm(apiUrl, data)
	if err != nil {
		logger.Errorf("url:%q, data:%v, error: %s", apiUrl, data, err)
		return err
	}
	defer resp.Body.Close()

	return nil
}

func fromName() string {
	if isPro {
		return "线上服务报警"
	}

	return "开发/测试服务报警"
}
