package service



type Service interface {

	SimpleData(username string) string

	AdminData(username string)  string

	// HealthCheck check service health status
	HealthCheck() bool
}

type CommentService struct {

}

func (s *CommentService) SimpleData(username string) string {
	return "hello " + username + " ,simple data, with simple authority"
}

func (s *CommentService) AdminData(username string) string {
	return "hello " + username + " ,admin data, with admin authority"

}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (s *CommentService) HealthCheck() bool {
	return true
}

func NewCommentService() *CommentService {
	return &CommentService{}
}

