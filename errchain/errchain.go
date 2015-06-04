package errchain

import "fmt"

// ErrorChain chainable errors
type ErrorChain interface {
	Error() string
	Cause() ErrorChain
}

type errChain struct {
	message string
	cause   ErrorChain
}

func (err errChain) Error() string {
	if err.cause == nil {
		return err.message + "\n"
	}
	return err.message + "\n    " + err.cause.Error() + "\n"
}

func (err errChain) Cause() ErrorChain {
	return err.cause
}

type errWrapper struct {
	error
}

func (err errWrapper) Cause() ErrorChain {
	return nil
}

// Wrap wraps an error with a message
// returns nil if err is nil
func Wrap(err error, message string) ErrorChain {
	if err == nil {
		return nil
	}

	if chain, ok := err.(ErrorChain); ok {
		return errChain{
			message: message,
			cause:   chain,
		}
	}
	return errWrapper{err}
}

// Wrapf wraps an error with a formatted message
// returns nil if err is nil
func Wrapf(err error, format string, args ...interface{}) ErrorChain {
	if err == nil {
		return nil
	}
	return Wrap(err, fmt.Sprintf(format, args))
}
