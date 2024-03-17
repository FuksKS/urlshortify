package storage

type saver interface {
	Save(cache map[string]string) error
}

type reader interface {
	Read() (map[string]string, error)
}
