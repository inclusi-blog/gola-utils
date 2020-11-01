package golang_error

type InternalServerError struct {
}

func (InternalServerError) Error() string {
	return "something went wrong"
}
