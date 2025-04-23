package types

type ServerError struct {
	Message  string
	Location string
	MainErr  error
}

func (e ServerError) Error() string {
	return e.Message + " - " + e.MainErr.Error()
}

func NewServerError(message string, location string, mainErr error) *ServerError {
	return &ServerError{
		Message:  message,
		Location: location,
		MainErr:  mainErr,
	}
}
