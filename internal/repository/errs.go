package repository

type (
	NotFound struct {
		What string
	}

	InternalServerError struct {
		Cause string
	}

	Unauthorized struct {
		Cause string
	}

	BadRequest struct {
		Cause string
	}
)

func (n NotFound) Error() string {
	return n.What
}

func (n InternalServerError) Error() string {
	return n.Cause
}

func (n Unauthorized) Error() string {
	return n.Cause
}
func (n BadRequest) Error() string {
	return n.Cause
}
