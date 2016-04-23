package depend

import (
	"net/url"
	"strconv"

	"github.com/bitly/go-simplejson"

	"golang.org/x/net/context"
)

// RecordUserOrder 记录用户完成的订单
func RecordUserOrder(ctx context.Context, uid, orderId, storeId int) error {
	data := url.Values{
		"order_id": {strconv.Itoa(orderId)},
		"store_id": {strconv.Itoa(storeId)},
	}
	_, err := callService(usercenterService, "PUT", "/user/"+strconv.Itoa(uid), data)
	if err != nil {
		return err
	}

	return nil
}

// ReadUserDevice 获取用户设备信息，用于推送，比如 client_id，client_source
func ReadUserDevice(ctx context.Context, uid int) *simplejson.Json {
	data := url.Values{"uid": {strconv.Itoa(uid)}}
	result, err := callService(usercenterService, "GET", "/device/client", data)
	if err != nil {
		return nil
	}

	return result.Get("user_device")
}
