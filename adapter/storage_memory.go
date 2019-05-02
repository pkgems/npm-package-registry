package adapter

// StorageMemory provides implementation of Storage interface
// for storing data in memory
type StorageMemory struct {
	data map[string][]byte
}

// NewStorageMemory is used to initialize new StorageMemory
func NewStorageMemory() *StorageMemory {
	return &StorageMemory{
		data: make(map[string][]byte),
	}
}

func (storage StorageMemory) WriteTarball(name, version string, data []byte) error {
	storage.data[name+version] = data

	return nil
}

func (storage StorageMemory) ReadTarball(name, version string) ([]byte, error) {
	return storage.data[name+version], nil
}
