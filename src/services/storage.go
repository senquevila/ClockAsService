package services

// StorageService defines basic operations for a storage-backed service
type StorageService interface {
	CreateTable() error
	Create(data interface{}) error
	Remove(id string) error
	List() ([]interface{}, error)
	FindByID(id string) (interface{}, error)
}