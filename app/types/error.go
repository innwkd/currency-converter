package types

type ConverterError string

func (c ConverterError) Error() string { return string(c) }

func IsConverterError(err error) bool {
	_, ok := err.(ConverterError)
	return ok
}
