package localstorage

import (
	"fmt"
	"sync"
)

// Example cache data
//"data": {
//	"21b7d460-f0ed-4bdb-ad6f-c65aac59f481": {
//		"Key": "21b7d460-f0ed-4bdb-ad6f-c65aac59f481"
//	},
//	"21bd27dc-e566-42ec-a517-e44e2c196cfb": {
//		"Key": "21bd27dc-e566-42ec-a517-e44e2c196cfb"
//	},
//}

type LocalStorage struct {
	*cache
}

var localStorage *LocalStorage

type cache struct {
	mu   sync.RWMutex
	data map[string]StorageItem
}

type StorageItem struct {
	Key string
}

func New() *LocalStorage {
	items := make(map[string]StorageItem)
	newCache := &cache{
		data: items,
	}
	return &LocalStorage{cache: newCache}
}

func (l *LocalStorage) AddToStorage(key string) error {
	if !l.Exists(key) {
		l.mu.Lock()
		newItem := struct {
			Key string
		}{Key: key}
		l.data[key] = newItem
		l.mu.Unlock()
		return nil
	}
	return fmt.Errorf("id: %v already exists", key)
}

func (l *LocalStorage) GetFromStorage(k string) (StorageItem, error) {
	var item StorageItem
	if !l.Exists(k) {
		return item, fmt.Errorf("item with key: %v does not exist", k)
	}
	return l.data[k], nil
}

func (l *LocalStorage) GetAllFromStorage() map[string]StorageItem {
	return l.data
}

func (l *LocalStorage) Remove(k string) {
	l.mu.Lock()
	delete(localStorage.data, k)
	l.mu.Unlock()
}

func (l *LocalStorage) Exists(k string) bool {
	if _, found := l.data[k]; found {
		return found
	}
	return false
}
