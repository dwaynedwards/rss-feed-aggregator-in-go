package rf

type DB interface {
	Open() error
	Close() error
}
