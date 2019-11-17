package service


type Service interface {
	// 健康检查
	HealthCheck() bool
	// 算术相加ice
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

// 用于检查服务的健康状态，这里仅仅返回true
func (*CalculateServiceImpl) HealthCheck() bool {
	return true
}

