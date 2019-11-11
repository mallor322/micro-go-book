package service

import (
	"encoding/json"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
	"github.com/keets2012/Micro-Go-Pracrise/ch10-resiliency/use-caculate/config"
	"github.com/keets2012/Micro-Go-Pracrise/common/discover"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Service interface {


	// HealthCheck check service health status
	HealthCheck() bool

	// calculateService
	UseCalculate(a int , b int ) (int, error)


}

var ErrServiceInstance  = errors.New("no instances are working")

type UseCalculateServiceImpl struct {
	discoveryClient discover.DiscoveryClient

}

func NewCalculateServiceImpl(client discover.DiscoveryClient) Service  {
	return &UseCalculateServiceImpl{
		client,
	}
}

type CalculateResponse struct {
	Result int `json:"result"`
	Error string `json:"error"`

}

func (service *UseCalculateServiceImpl) UseCalculate(a int , b int ) (int, error){

	serviceName := "Calculate"

	var addResult int

	err := hystrix.Do("Calculate.calculate", func() error {

		instances := service.discoveryClient.DiscoverServices(serviceName, config.Logger)

		if instances == nil || len(instances) == 0 {
			config.Logger.Println("No Calculate instances are working!")
			return ErrServiceInstance
		}

		// 随机选取一个服务实例进行计算
		rand.Seed(time.Now().UnixNano())
		selectInstance := instances[rand.Intn(len(instances))].(*api.AgentService)

		requestUrl := url.URL{
			Scheme:   "http",
			Host:     selectInstance.Address + ":" + strconv.Itoa(selectInstance.Port),
			Path:     "/calculate",
			RawQuery: "a=" + strconv.Itoa(a) + "&b=" + strconv.Itoa(b),
		}

		resp, err := http.Get(requestUrl.String())
		if err != nil {
			return err
		}
		result := &CalculateResponse{}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil{
			return err
		}else if result.Error != ""{
			return errors.New(result.Error)
		}

		addResult = result.Result

		return nil

	}, func(e error) error {
		return e
	})

	return addResult, err

}



// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (*UseCalculateServiceImpl) HealthCheck() bool {
	return true
}

