package service

import (
	"fmt"
	conf "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/config"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-app/config"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-app/model"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-app/service/srv_err"
	"log"
	"time"
)

// Service Define a service interface
type Service interface {
	// HealthCheck check service health status
	HealthCheck() bool
	SecInfo(productId int) (date map[string]interface{})
	SecKill(req *model.SecRequest) (map[string]interface{}, int, error)
	SecInfoList() ([]map[string]interface{}, int, error)
}

//UserService implement Service interface
type SkAppService struct {
}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (s SkAppService) HealthCheck() bool {
	return true
}

type ServiceMiddleware func(Service) Service

func (s SkAppService) SecInfo(productId int) (date map[string]interface{}) {
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()

	v, ok := conf.SecKill.SecProductInfoMap[productId]
	if !ok {
		return nil
	}

	data := make(map[string]interface{})
	data["product_id"] = productId
	data["start_time"] = v.StartTime
	data["end_time"] = v.EndTime
	data["status"] = v.Status

	return data
}

func (s SkAppService) SecKill(req *model.SecRequest) (map[string]interface{}, int, error) {
	//对Map加锁处理
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()
	var code int
	//err := srv_limit.UserCheck(req)
	//if err != nil {
	//	code = srv_err.ErrUserCheckAuthFailed
	//	log.Printf("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
	//	return nil, code, err
	//}
	//
	//err = srv_limit.AntiSpam(req)
	//if err != nil {
	//	code = srv_err.ErrUserServiceBusy
	//	log.Printf("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
	//	return nil, code, err
	//}

	data, code, err := SecInfoById(req.ProductId)
	if err != nil {
		log.Printf("userId[%d] secInfoById Id failed, req[%v]", req.UserId, req)
		return nil, code, err
	}

	userKey := fmt.Sprintf("%d_%d", req.UserId, req.ProductId)
	fmt.Println("userKey : ", userKey)
	config.SkAppContext.UserConnMap[userKey] = req.ResultChan
	//将请求送入通道并推入到redis队列当中
	config.SkAppContext.SecReqChan <- req
	log.Printf("userId [%d] [%d]", time.Duration(conf.SecKill.AppWaitResultTimeout), conf.SecKill.AppWaitResultTimeout)

	ticker := time.NewTicker(time.Millisecond * time.Duration(conf.SecKill.AppWaitResultTimeout))

	defer func() {
		ticker.Stop()
		config.SkAppContext.UserConnMapLock.Lock()
		delete(config.SkAppContext.UserConnMap, userKey)
		config.SkAppContext.UserConnMapLock.Unlock()
	}()

	select {
	case <-ticker.C:
		code = srv_err.ErrProcessTimeout
		err = fmt.Errorf("request timeout")
		return nil, code, err
	case <-req.CloseNotify:
		code = srv_err.ErrClientClosed
		err = fmt.Errorf("client already closed")
		return nil, code, err
	case result := <-req.ResultChan:
		code = result.Code
		if code != 1002 {
			return data, code, srv_err.GetErrMsg(code)
		}
		log.Printf("secKill success")
		data["product_id"] = result.ProductId
		data["token"] = result.Token
		data["user_id"] = result.UserId
		return data, code, nil
	}
}

func NewSecRequest() *model.SecRequest {
	secRequest := &model.SecRequest{
		ResultChan: make(chan *model.SecResult, 1),
	}
	return secRequest
}

func (s SkAppService) SecInfoList() ([]map[string]interface{}, int, error) {
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()

	var data []map[string]interface{}
	for _, v := range conf.SecKill.SecProductInfoMap {
		item, _, err := SecInfoById(v.ProductId)
		if err != nil {
			log.Printf("get sec info, err : %v", err)
			continue
		}
		data = append(data, item)
	}
	return data, 0, nil
}

func SecInfoById(productId int) (map[string]interface{}, int, error) {
	//对Map加锁处理
	config.SkAppContext.RWSecProductLock.RLock()
	defer config.SkAppContext.RWSecProductLock.RUnlock()

	var code int
	v, ok := conf.SecKill.SecProductInfoMap[productId]
	if !ok {
		return nil, srv_err.ErrNotFoundProductId, fmt.Errorf("not found product_id:%d", productId)
	}

	start := false      //秒杀活动是否开始
	end := false        //秒杀活动是否结束
	status := "success" //状态

	nowTime := time.Now().Unix()
	log.Printf("now time is ", nowTime)
	//秒杀活动没有开始
	if nowTime-v.StartTime < 0 {
		start = false
		end = false
		status = "second kill not start"
		code = srv_err.ErrActiveNotStart
	}

	//秒杀活动已经开始
	if nowTime-v.StartTime > 0 {
		start = true
	}

	//秒杀活动已经结束
	if nowTime-v.EndTime > 0 {
		start = false
		end = true
		status = "second kill is already end"
		code = srv_err.ErrActiveAlreadyEnd
	}

	//商品已经被停止或售磬
	if v.Status == config.ProductStatusForceSaleOut || v.Status == config.ProductStatusSaleOut {
		start = false
		end = false
		status = "product is sale out"
		code = srv_err.ErrActiveSaleOut
	}

	//组装数据
	data := map[string]interface{}{
		"product_id": productId,
		"start":      start,
		"end":        end,
		"status":     status,
	}
	return data, code, nil
}
