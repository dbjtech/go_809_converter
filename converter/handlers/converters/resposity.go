package converters

type ConverterInterface interface {
	Convert() ([]byte, error)
}
