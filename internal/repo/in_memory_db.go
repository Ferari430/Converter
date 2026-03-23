package repo

import "sync"

type InMemoryDatabase struct {
	mu      sync.RWMutex
	images  map[string]string
	mdFiles map[string]string
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		images:  make(map[string]string),
		mdFiles: make(map[string]string),
	}
}

func (db *InMemoryDatabase) AddImage(filename, path string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.images[filename] = path
}

func (db *InMemoryDatabase) GetImages() map[string]string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	result := make(map[string]string)
	for k, v := range db.images {
		result[k] = v
	}
	return result
}

func (db *InMemoryDatabase) ClearImages() {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.images = make(map[string]string)
}

// MD File methods
func (db *InMemoryDatabase) AddMDFile(filename, path string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.mdFiles[filename] = path
}

func (db *InMemoryDatabase) GetMDFiles() map[string]string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	// Return a copy to prevent external modifications
	result := make(map[string]string)
	for k, v := range db.mdFiles {
		result[k] = v
	}
	return result
}

func (db *InMemoryDatabase) ClearMDFiles() {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.mdFiles = make(map[string]string)
}

func (db *InMemoryDatabase) Clear() {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.images = make(map[string]string)
	db.mdFiles = make(map[string]string)
}
