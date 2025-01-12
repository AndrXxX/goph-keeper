package filestorage

type hashGenerator interface {
	Generate(data []byte) string
}
