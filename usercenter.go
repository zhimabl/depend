package depend

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

// RecordUserOrder 记录用户完成的订单
func RecordUserOrder(ctx context.Context, uid, orderId, storeId int) error {
	usercenterConf := randServiceConf(usercenterService)
	if usercenterConf == nil {
		logger.Errorln(usercenterService, "config is empty")
		return errors.New("usercenter service config is empty")
	}

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	data := url.Values{
		"order_id":  {strconv.Itoa(orderId)},
		"store_id":  {strconv.Itoa(storeId)},
		"timestamp": {strconv.FormatInt(time.Now().Unix(), 10)},
		"from":      {from},
	}

	data.Set("sign", goutils.GenSign(data, getServiceSecret(usercenterService)))

	apiUrl := "http://" + usercenterConf.httpAddr + "/user/" + strconv.Itoa(uid)
	resp, err := putForm(httpClient, apiUrl, data)
	if err != nil {
		logger.Errorf("url:%q, data:%v, error: %v", apiUrl, data, err)
		return err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("url:%q, data:%v, result: %s, error: %v", apiUrl, data, result, err)
		return err
	}

	return nil
}
