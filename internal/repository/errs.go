package repository

type (
	NotFound struct {
		What string
	}

	InternalServerError struct {
		Cause string
	}
)

func (n NotFound) Error() string {
	return n.What
}

func (n InternalServerError) Error() string {
	return n.Cause
}
