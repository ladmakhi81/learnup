package env

type EnvProvider interface {
	LoadLearnUp() (*EnvConfig, error)
}
