package storage

import (
	"fmt"
	"log/slog"
	"sync"
)

type KeyValueStoreOperations interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

type KeyValueStore struct {
	data   map[string]string
	mu     sync.RWMutex
	logger slog.Logger
}

func NewKeyValueStore(logger slog.Logger) *KeyValueStore {
	return &KeyValueStore{
		data:   make(map[string]string),
		logger: logger,
	}
}

func (kvs *KeyValueStore) Get(key string) (string, error) {
	kvs.mu.RLock()
	defer kvs.mu.RUnlock()

	kvs.logger.Info("Get Request for Key", "key", key)
	value, exists := kvs.data[key]
	if !exists {
		kvs.logger.Error("The Key doesn't exist", "key", key)
		return "", fmt.Errorf("The Key %s doesn't exists", key)
	}
	return value, nil
}

func (kvs *KeyValueStore) Set(key string, value string) error {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	kvs.logger.Info("Set Request", "key", key, "value", value)
	kvs.data[key] = value
	_, exist := kvs.data[key]
	if !exist {
		kvs.logger.Error("Failed to set key", "key", key)
		return fmt.Errorf("Failed to set the key %s", key)
	}
	kvs.logger.Info("Successfully set the key", "key", key)
	return nil
}

func (kvs *KeyValueStore) Delete(key string) error {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	kvs.logger.Info("Delete Request for key: ", "key", key)
	_, exist := kvs.data[key]
	if !exist {
		kvs.logger.Error("Failed to delete key", "key", key)
		return fmt.Errorf("Key doesn't exist. Failed to delete the key %s", key)
	}

	delete(kvs.data, key)

	_, existAfterDelete := kvs.data[key]

	if existAfterDelete {
		kvs.logger.Error("Failed to delete the key:", "key", key)
		return fmt.Errorf("Failed to delete the key %s", key)
	}

	kvs.logger.Info("Key deleted successfully")
	return nil
}
