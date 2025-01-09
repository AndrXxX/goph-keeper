package storageadapters

type convertor[T1 any, T2 any] interface {
	Convert(e *T1) *T2
}
