package depend

import (
	"net/url"
	"strconv"

	"golang.org/x/net/context"
)

// ChangeCredits 变更用户积分
func ChangeCredits(ctx context.Context, uid, creditNum int, remark string) error {
	data := url.Values{
		"uid":        {strconv.Itoa(uid)},
		"credit_num": {strconv.Itoa(creditNum)},
		"remark":     {remark},
	}
	_, err := callService(zhimauserService, "POST", "/user/change_credits", data)
	if err != nil {
		return err
	}

	return nil
}

// RefundCredits 退还用户积分
func RefundCredits(ctx context.Context, uid, creditNum int, remark string) error {
	data := url.Values{
		"uid":        {strconv.Itoa(uid)},
		"credit_num": {strconv.Itoa(creditNum)},
		"remark":     {remark},
	}
	_, err := callService(zhimauserService, "POST", "/user/refund_credits", data)
	if err != nil {
		return err
	}

	return nil
}
