package helpers

type ErrRecordNotFound struct {
	Err error
}

type ErrDuplicateEmail struct {
	Err error
}

func (m *ErrRecordNotFound) Error() string {
	return "record not found in DB"
}

func (m *ErrDuplicateEmail) Error() string {
	return "record with email already exists in DB"
}
