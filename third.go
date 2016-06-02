package depend

import (
	"net/url"

	"golang.org/x/net/context"
)

// ThirdConfirm 第三方订单确认
func ThirdOrderConfirm(ctx context.Context, orderId string) error {
	data := url.Values{
		"order_id": {orderId},
	}
	_, err := callService(thirdService, "POST", "/order/confirm", data)
	if err != nil {
		return err
	}

	return nil
}

// ThirdConfirm 配送中
func ThirdOrderDelivering(ctx context.Context, orderId string) error {
	data := url.Values{
		"order_id": {orderId},
	}
	_, err := callService(thirdService, "POST", "/order/delivering", data)
	if err != nil {
		return err
	}

	return nil
}

// ThirdOrderFinish 完成
func ThirdOrderFinish(ctx context.Context, orderId string) error {
	data := url.Values{
		"order_id": {orderId},
	}
	_, err := callService(thirdService, "POST", "/order/finish", data)
	if err != nil {
		return err
	}

	return nil
}

// ThirdOrderCancel 取消
func ThirdOrderCancel(ctx context.Context, orderId string) error {
	data := url.Values{
		"order_id": {orderId},
	}
	_, err := callService(thirdService, "POST", "/order/cancel", data)
	if err != nil {
		return err
	}

	return nil
}
