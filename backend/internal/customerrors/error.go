package customerrors

import "errors"

const (
	ItemAlreadyExist       int = 4001
	ItemNotExist           int = 4041
	DatabaseExecutionError int = 5001
	CommonErr              int = 5002
	InvalidAction          int = 4002
)

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrItemAlreadyExist = errors.New("item already exist")
)

type CustomError struct {
	ErrUser string
	ErrLog  error
	ErrCode int
}

func (c *CustomError) Error() string {
	return c.ErrUser
}

func (c *CustomError) GetStatusCode() int {
	return c.ErrCode / 10
}

func NewError(userErr string, logErr error, code int) error {
	return &CustomError{
		ErrUser: userErr,
		ErrLog:  logErr,
		ErrCode: code,
	}
}
