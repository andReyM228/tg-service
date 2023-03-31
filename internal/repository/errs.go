package repository

type (
	NotFound struct {
		NotFound string
	}

	InternalServerError struct {
		InternalServerError string
	}
)

func (n NotFound) Error() string {
	return "Not found: " + n.NotFound
}

func (n InternalServerError) Error() string {
	return "Internal server error: " + n.InternalServerError
}
