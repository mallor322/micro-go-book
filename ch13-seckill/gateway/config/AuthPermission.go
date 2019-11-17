package config

var (
	AuthPermitConfig AuthPermitAll
)
//Http配置
type AuthPermitAll struct {
	PermitAll []interface{}
}

