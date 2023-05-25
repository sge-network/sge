package convertors

type IConverter[T any, U any] interface {
	Convert(T) (U, error)
}
