package adapter

// Database provides base interace for database adapters
type Database interface {
	SetPackage(name string, data string) error
	GetPackage(name string) (string, error)
}

// Storage provides base interace for storage adapters
type Storage interface {
	WriteTarball(name, version string, data []byte) error
	ReadTarball(name, version string) ([]byte, error)
}
