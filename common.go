package depend

import (
	"encoding/json"
	"math/rand"
	"sync"

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

var services = []string{senderService}

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
