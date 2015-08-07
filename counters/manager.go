package counters

import (
	"counts/counters/abstract"
	"counts/counters/immutable"
	"errors"

	"github.com/hashicorp/golang-lru"
)

/*
Manager is responsible for manipulating the counters and syncing to disk
*/
type Manager struct {
	cache *lru.Cache
}

var manager *Manager

/*
CreateDomain ...
*/
func (m *Manager) CreateDomain(domainID string, domainType string) error {
	//TODO: spit errir uf domainType is invalid
	//FIXME: no hardcoding of immutable here

	info := abstract.Info{ID: domainID, Type: domainType}
	switch domainType {
	case "immutable":
		m.cache.Add(info.ID, immutable.NewDomain(info))
	}
	return nil
}

/*
DeleteDomain ...
*/
func (m *Manager) DeleteDomain(domainID string) error {
	return nil
}

/*
GetDomains ...
*/
func (m *Manager) GetDomains() ([]string, error) {
	// TODO: Remove dummy data and implement proper result
	values := manager.cache.Keys()
	domains := make([]string, len(values), len(values))
	for i, v := range values {
		domains[i] = v.(string)
	}
	return domains, nil
}

/*
AddToDomain ...
*/
func (m *Manager) AddToDomain(domainID string, values []string) error {
	var val, ok = m.cache.Get(domainID)
	if ok == false {
		return errors.New("No such domain: " + domainID)
	}
	var counter abstract.Counter
	counter = val.(abstract.Counter)

	bytes := make([][]byte, len(values), len(values))
	for i, value := range values {
		bytes[i] = []byte(value)
	}
	counter.AddMultiple(bytes)
	return nil
}

/*
DeleteFromDomain ...
*/
func (m *Manager) DeleteFromDomain(domainID string, values []string) error {
	return nil
}

/*
GetCountForDomain ...
*/
func (m *Manager) GetCountForDomain(domainID string) (uint, error) {

	var val, ok = m.cache.Get(domainID)
	if ok == false {
		return 0, errors.New("No such domain: " + domainID)
	}
	var counter abstract.Counter
	counter = val.(abstract.Counter)
	count := counter.GetCount()
	return count, nil
}

/*
GetManager returns a singleton Manager
*/
func GetManager() *Manager {
	if manager == nil {
		cache, _ := lru.New(100)
		manager = &Manager{cache}
	}
	return manager
}