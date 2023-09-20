package actions

type MultiError interface {
	Unwrap() []error
}
