package contracts

type Cache interface {
	SetVal(key string, val any) error
	SetHashVal(key, id string, val any) error
	GetHashVal(key, id string) (string, error)
	GetVal(key string) (string, error)
}
