package service


type Service interface {


	// HealthCheck check service health status
	HealthCheck() bool

	// calculateService
	Calculate(a int , b int ) int


}


type CalculateServiceImpl struct {
}

func NewCalculateServiceImpl() Service  {
	return &CalculateServiceImpl{}
}

func (*CalculateServiceImpl) Calculate(a int , b int ) int {
	return a + b
}



// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (*CalculateServiceImpl) HealthCheck() bool {
	return true
}

