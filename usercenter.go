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

// ReadUser 获取用户信息
func ReadUser(ctx context.Context, uid int) *simplejson.Json {
	result, err := callService(usercenterService, "GET", "/user/"+strconv.Itoa(uid), url.Values{})
	if err != nil {
		return nil
	}

	return result.Get("user")
}

// IsChannelUser 判断手机号是否是渠道用户（预注册过）
func IsChannelUser(ctx context.Context, mobile string) bool {
	result, err := callService(usercenterService, "GET", "/channel/"+mobile, url.Values{})
	if err != nil {
		return false
	}

	return result.Get("channel_user").MustBool(false)
}

// IsChannelUser 判断UID是否是渠道用户（预注册过）
func IsChannelUserByUid(ctx context.Context, uid int) bool {
	data := url.Values{"uid": {strconv.Itoa(uid)}}
	result, err := callService(usercenterService, "GET", "/channel/user", data)
	if err != nil {
		return false
	}

	return result.Get("channel_user").MustBool(false)
}
