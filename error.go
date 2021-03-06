package depend

import (
	"fmt"
	"runtime"

	"github.com/polaris1119/logger"
	"golang.org/x/net/context"
)

// ProcessNotInterruptErr 处理一些非中断型错误，一般不会出现。
// 这些错误发生时，程序依然往下执行，这里只是记录错误，同时报警
func ProcessNotInterruptErr(ctx context.Context, idName string, idVal interface{}, err error) {
	if err == nil {
		return
	}

	errMsg := fmt.Sprintf("%q=%v has error:%#v", idName, idVal, err)

	_, file, line, ok := runtime.Caller(1)
	if ok {
		errMsg = fmt.Sprintf("%s in file(%q) on line(%d)", errMsg, file, line)
	}
	recordLog(ctx, errMsg)

	emailInfo := map[string]string{
		"subject": "非中断型错误发生了，请留意！",
		"content": errMsg,
	}
	go ServiceWarning(ctx, emailInfo, fmt.Sprintf("非中断型错误发生了，请留意。%s=%q", idName, idVal))
}

func recordLog(ctx context.Context, errMsg string) {
	if ctx != nil {
		_logger, ok := ctx.Value("logger").(*logger.Logger)
		if ok {
			_logger.Errorln(errMsg)
			return
		}
	}

	logger.Errorln(errMsg)
}
