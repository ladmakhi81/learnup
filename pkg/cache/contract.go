package cache

type Cache interface {
	SetHashVal(key, id string, val any) error
	GetHashVal(key, id string) (string, error)
}
