package depend

import (
	"net/url"
	"strconv"
)

// SendVoucher 给用户发送优惠券
func SendVoucher(uid, typ, state int) error {
	data := url.Values{
		"uid":   {strconv.Itoa(uid)},
		"typ":   {strconv.Itoa(typ)},
		"state": {strconv.Itoa(state)},
		"async": {"true"},
	}
	_, err := callService(voucherService, "POST", "/voucher/send", data)
	if err != nil {
		return err
	}

	return nil
}
