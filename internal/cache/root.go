package cache

import (
	"encoding/json"
	"fmt"
	"os"
)

const db = "internal/db/cache.json"
func Get(key string) (*string, error) {
	cacheData, err := os.ReadFile(db)
	if err != nil {
		return nil, fmt.Errorf("error reading cache file: %v", err)
	}

	var fileCache map[string]string
	if err := json.Unmarshal(cacheData, &fileCache); err != nil {
		return nil, fmt.Errorf("error parsing cache file: %v", err)
	}

	value, ok := fileCache[key]
	if !ok {
		return nil, nil
	}
	return &value, nil
}

func Set(key string, value string) (error) {
	cacheData, err := os.ReadFile(db)
	if err != nil {
		return fmt.Errorf("error reading cache file: %v", err)
	}

	var fileCache map[string]string
	if err := json.Unmarshal(cacheData, &fileCache); err != nil {
		return fmt.Errorf("error parsing cache file: %v", err)
	}
	fileCache[key] = value
	cacheJson, err := json.Marshal(fileCache)
	if err != nil {
		return fmt.Errorf("error writing cache: %v", err)
	}

	if err := os.WriteFile("internal/db/cache.json", cacheJson, 0644); err != nil {
		return fmt.Errorf("error saving cache to db: %v", err)
	}
	return nil
}

func ClearAll() {
	os.WriteFile(db, []byte("{}"), 0644)
}