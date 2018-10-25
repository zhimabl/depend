package depend

import (
	"net/url"
	"strconv"

	"golang.org/x/net/context"
)

const (
	ProjectCms = iota
	ProjectZhimazg
	ProjectSupply
	ProjectDriver
)

// CancelExpress 取消物流单
func CancelExpress(ctx context.Context, expressSn int64, project, cancelType int, cancelReason string) error {
	data := url.Values{
		"express_sn":    {strconv.FormatInt(expressSn, 10)},
		"cancel_type":   {strconv.Itoa(cancelType)},
		"cancel_reason": {cancelReason},
		"project":       {strconv.Itoa(project)},
	}
	_, err := callService(orderService, "POST", "/order/cancel", data)
	if err != nil {
		return err
	}

	return nil
}

// CancelOrder 取消订单，如果有物流单，一起取消；
// 主要为了方便，可以一次性全部取消或者还没有物流单时候的取消
func CancelOrder(ctx context.Context, orderSn int64, project, cancelType int, cancelReason string) error {
	data := url.Values{
		"order_sn":      {strconv.FormatInt(orderSn, 10)},
		"cancel_type":   {strconv.Itoa(cancelType)},
		"cancel_reason": {cancelReason},
		"project":       {strconv.Itoa(project)},
	}
	_, err := callService(orderService, "POST", "/order/cancel", data)
	if err != nil {
		return err
	}

	return nil
}
