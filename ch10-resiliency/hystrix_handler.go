package main

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

type HystrixHandler struct {

	// 记录hystrix是否已配置
	hystrixs      map[string]bool
	hystrixsMutex *sync.Mutex

	consulClient *api.Client
	logger       log.Logger
}

func NewHystrixHandler(consulClient *api.Client, logger log.Logger) *HystrixHandler {

	return &HystrixHandler{
		consulClient:  consulClient,
		logger:        logger,
		hystrixs:      make(map[string]bool),
		hystrixsMutex: &sync.Mutex{},
	}

}

func (hystrixHandler *HystrixHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	reqPath := req.URL.Path
	if reqPath == "" {
		return
	}
	//按照分隔符'/'对路径进行分解，获取服务名称serviceName
	pathArray := strings.Split(reqPath, "/")
	serviceName := pathArray[1]

	if serviceName == "" {
		// 路径不存在
		rw.WriteHeader(404)
		return
	}

	if _, ok := hystrixHandler.hystrixs[serviceName]; !ok {
		hystrixHandler.hystrixsMutex.Lock()
		if _, ok := hystrixHandler.hystrixs[serviceName]; !ok {
			//把serviceName作为 hystrix 命令命名
			hystrix.ConfigureCommand(serviceName, hystrix.CommandConfig{
				// 进行 hystrix 命令自定义
			})
			hystrixHandler.hystrixs[serviceName] = true
		}
		hystrixHandler.hystrixsMutex.Unlock()
	}

	err := hystrix.Do(serviceName, func() error {

		//调用consul api查询serviceName的服务实例列表
		result, _, err := hystrixHandler.consulClient.Catalog().Service(serviceName, "", nil)
		if err != nil {
			hystrixHandler.logger.Log("ReverseProxy failed", "query service instace error", err.Error())
			return errors.New("query service instace error")
		}

		if len(result) == 0 {
			hystrixHandler.logger.Log("ReverseProxy failed", "no such service instance", serviceName)
			return errors.New("no such service instance " + serviceName)
		}

		//创建Director
		director := func(req *http.Request) {

			//重新组织请求路径，去掉服务名称部分
			destPath := strings.Join(pathArray[2:], "/")

			//随机选择一个服务实例
			tgt := result[rand.Int()%len(result)]
			hystrixHandler.logger.Log("service id", tgt.ServiceID)

			//设置代理服务地址信息
			req.URL.Scheme = "http"
			req.URL.Host = fmt.Sprintf("%s:%d", tgt.ServiceAddress, tgt.ServicePort)
			req.URL.Path = "/" + destPath
		}

		var proxyError error

		// 返回代理异常，用于记录 hystrix.Do 执行失败
		errorHandler := func(ew http.ResponseWriter, er *http.Request, err error) {
			proxyError = err
		}

		proxy := &httputil.ReverseProxy{
			Director:     director,
			ErrorHandler: errorHandler,
		}

		proxy.ServeHTTP(rw, req)

		// 将执行异常反馈 hystrix
		return proxyError

	}, func(e error) error {
		hystrixHandler.logger.Log("proxy error ", e)
		return errors.New("fallback excute")
	})

	// hystrix.Do 返回执行异常
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(err.Error()))
	}

}
