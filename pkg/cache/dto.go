package cache

type CacheError struct {
	Message  string
	Location string
}

func (e CacheError) Error() string {
	return e.Message
}

func NewCacheError(msg string, location string) *CacheError {
	return &CacheError{
		Message:  msg,
		Location: location,
	}
}
