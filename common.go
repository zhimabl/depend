package depend

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/polaris1119/goutils"
	"github.com/polaris1119/logger"
	"github.com/polaris1119/nosql"
)

var (
	// 调用方的名称，如 usercenter
	from string
	// 是否是生成环境
	isPro bool
)

// Init 调用方需要先调用该方法，初始化 depend
func Init(_from string, _isPro bool) {
	from = _from
	isPro = _isPro
}

type serviceConf struct {
	httpAddr  string
	secret    string
	failTimes int
}

const (
	senderService     = "sender_service"
	usercenterService = "usercenter_service"
	opRecordService   = "op_record_service"
	voucherService    = "voucher_service"
)

var services = []string{senderService, usercenterService, opRecordService, voucherService}

var (
	servicesMap = make(map[string][]*serviceConf, len(services))
	locker      sync.RWMutex
)

// LoadServicesConf 装载所有服务配置
func LoadServicesConf() {
	tmpServicesMap := make(map[string][]*serviceConf, len(services))

	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("LoadServicesConf panic:", err)
		}
	}()

	redisClient := nosql.NewRedisClientWithSection("redis.conf")
	defer redisClient.Close()

	redisClient.NoPrefix = true

	for _, service := range services {
		key := service + ":service_addr_list"
		dataMap, err := redisClient.HGETALL(key)
		if err != nil {
			logger.Errorf("load services conf hgetall %q error:%v", key, err)
			continue
		}

		serviceConfSlice := make([]*serviceConf, 0, len(dataMap))
		for addr, val := range dataMap {
			valMap := make(map[string]interface{}, 2)
			err = json.Unmarshal([]byte(val), &valMap)
			if err != nil {
				logger.Errorf("json decode error:%v", err)
				continue
			}

			objServiceConf := &serviceConf{
				httpAddr:  addr,
				secret:    valMap["secret"].(string),
				failTimes: int(valMap["fail_times"].(float64)),
			}
			serviceConfSlice = append(serviceConfSlice, objServiceConf)
		}

		if len(serviceConfSlice) > 0 {
			tmpServicesMap[service] = serviceConfSlice
		}
	}

	locker.Lock()
	servicesMap = tmpServicesMap
	locker.Unlock()
}

type resultStruct struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data *simplejson.Json `json:"data"`
}

var ServiceFailErr = errors.New("remote service return code is not 0")

// callService 调用远程服务
// TODO: 重试指数退避算法；错误自动报警？
func callService(serviceName, method, uri string, data url.Values) (*simplejson.Json, error) {
	serviceConf := randServiceConf(serviceName)
	if serviceConf == nil {
		logger.Errorln(serviceName, "config is empty")
		return nil, errors.New(serviceName + " service config is empty")
	}

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	data.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	data.Set("from", from)
	data.Set("sign", goutils.GenSign(data, getServiceSecret(serviceName)))

	apiUrl := "http://" + serviceConf.httpAddr + uri

	var (
		resp *http.Response
		err  error
	)
	switch strings.ToUpper(method) {
	case "PUT":
		resp, err = putForm(httpClient, apiUrl, data)
	case "POST":
		resp, err = httpClient.PostForm(apiUrl, data)
	case "GET":
		resp, err = httpClient.Get(apiUrl + "?" + data.Encode())
	}
	if err != nil {
		logger.Errorf("url:%q, data:%v, error: %v", apiUrl, data, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		logger.Errorf("url:%q, data:%v, status: %s(%d)", apiUrl, data, resp.Status, resp.StatusCode)
		return nil, errors.New("status code is not ok")
	}

	resultBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("url:%q, data:%v, result: %s, error: %v", apiUrl, data, resultBuf, err)
		return nil, err
	}

	result := &resultStruct{}
	err = json.Unmarshal(resultBuf, result)
	if err != nil {
		return nil, err
	}

	// 业务失败
	if result.Code != 0 {
		logger.Errorf("url:%q, data:%v, result: %s, code not 0", apiUrl, data, resultBuf)
		return nil, ServiceFailErr
	}

	return result.Data, nil
}

func getServiceConfSlice(serviceName string) []*serviceConf {
	locker.RLock()
	defer locker.RUnlock()

	return servicesMap[serviceName]
}

func randServiceConf(serviceName string) *serviceConf {
	locker.RLock()
	defer locker.RUnlock()

	serviceConfSlice := servicesMap[serviceName]
	if len(serviceConfSlice) > 0 {
		return serviceConfSlice[rand.Intn(len(serviceConfSlice))]
	}

	return nil
}

func getServiceSecret(serviceName string) string {
	locker.RLock()
	defer locker.RUnlock()

	serviceConfSlice := servicesMap[serviceName]
	if len(serviceConfSlice) > 0 {
		return serviceConfSlice[0].secret
	}
	return ""
}

func putForm(client *http.Client, apiUrl string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest("PUT", apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return client.Do(req)
}
