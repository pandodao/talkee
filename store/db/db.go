package db

type Scanner interface {
	Scan(dest ...interface{}) error
	StructScan(dest interface{}) error
	Next() bool
	Close() error
}
