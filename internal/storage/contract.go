package storage

type saver interface {
	Save(cache map[string]string) error
}

type reader interface {
	read() (map[string]string, error)
}
