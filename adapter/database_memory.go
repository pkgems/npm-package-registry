package adapter

// DatabaseMemory provides implementation of Database interface
// for storing data in memory
type DatabaseMemory struct {
	data map[string]string
}

// OptionsDatabaseMemory parameterizes DatabaseMemory
type OptionsDatabaseMemory struct{}

// NewDatabaseMemory is used to initialize new DatabaseMemory
func NewDatabaseMemory() *DatabaseMemory {
	return &DatabaseMemory{
		data: make(map[string]string),
	}
}

func (database DatabaseMemory) SetPackage(name string, data string) error {
	database.data[name] = data

	return nil
}

func (database DatabaseMemory) GetPackage(name string) (string, error) {
	return database.data[name], nil
}
