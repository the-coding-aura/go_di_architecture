package module

import (
	"errors"
	"go_di_architecture/internal/domain/models/module"
	"strconv"
	"strings"
	"sync"
)

type ModuleRepository struct {
	data            map[int]*module.Module
	mu              sync.Mutex
	autoIncrementID int
}

func NewModuleRepository() *ModuleRepository {
	return &ModuleRepository{
		data:            make(map[int]*module.Module),
		autoIncrementID: 1,
	}
}

func (r *ModuleRepository) CreateModule(m *module.Module) (*module.Module, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Simulate auto-increment ID
	m.ID = r.autoIncrementID
	r.autoIncrementID++

	r.data[m.ID] = m
	return m, nil
}

func (r *ModuleRepository) IsModuleNameExists(name string, excludeId int) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, mod := range r.data {
		if strings.EqualFold(mod.Name, name) && id != excludeId {
			return true, nil
		}
	}
	return false, nil
}

func (r *ModuleRepository) GetModuleById(id string) (*module.Module, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	moduleID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	m, exists := r.data[moduleID]
	if !exists {
		return nil, nil
	}
	return m, nil
}
