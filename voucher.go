package depend

import (
	"net/url"
	"strconv"

	"golang.org/x/net/context"
)

// SendVoucher 给用户发送优惠券
func SendVoucher(uid, typ, state int) error {
	data := url.Values{
		"uid":   {strconv.Itoa(uid)},
		"typ":   {strconv.Itoa(typ)},
		"state": {strconv.Itoa(state)},
		"async": {"true"},
	}
	_, err := callService(couponService, "POST", "/voucher/send", data)
	if err != nil {
		return err
	}

	return nil
}

// ReturnVoucher 返回优惠券：一般指订单取消后
func ReturnVoucher(ctx context.Context, voucherID int) error {
	data := url.Values{
		"voucher_id": {strconv.Itoa(voucherID)},
		"async":      {"true"},
	}
	_, err := callService(couponService, "POST", "/voucher/go_back", data)
	if err != nil {
		return err
	}

	return nil
}
