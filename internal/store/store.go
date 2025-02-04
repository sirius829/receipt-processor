package store

import (
	"receipt-processor/internal/models"

	"github.com/google/uuid"
)

type Store interface {
	Save(id string, receipt models.Receipt)
	Get(id string) (models.Receipt, bool)
}

type memoryStore struct {
	data map[string]models.Receipt
}

func NewMemoryStore() Store {
	return &memoryStore{data: make(map[string]models.Receipt)}
}

func (m *memoryStore) Save(id string, receipt models.Receipt) {
	m.data[id] = receipt
}

func (m *memoryStore) Get(id string) (models.Receipt, bool) {
	receipt, ok := m.data[id]
	return receipt, ok
}

func GenerateID() string {
	return uuid.New().String()
}
