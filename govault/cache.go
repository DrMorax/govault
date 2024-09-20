package govault

import (
	"container/list"
	"reflect"
	"sync"
	"unsafe"
)

// Cache is a generic in-memory cache with a memory limit (measured in bytes).
type Cache[Key comparable, Value any] struct {
	mutex     sync.Mutex
	store     map[Key]*list.Element
	evictList *list.List // List to track access order for LRU
	maxSize   int64      // Max memory size in bytes
	currSize  int64      // Current memory usage in bytes
}

// entry holds both the key and value, and the memory size of the value.
type entry[Key comparable, Value any] struct {
	key   Key
	value Value
	size  int64 // Estimated memory size in bytes
}

// New creates a new cache instance with a memory limit *measured in MegaBytes*.
func New[Key comparable, Value any](maxMB int64) *Cache[Key, Value] {
	if maxMB <= 0 {
		panic("maxMB mutexst be greater than zero")
	}

	return &Cache[Key, Value]{
		store:     make(map[Key]*list.Element),
		evictList: list.New(),
		maxSize:   maxMB * 1024 * 1024, // Convert MB to bytes
	}
}

// Set adds or updates a key-value pair in the cache.
// If the cache exceeds the memory limit, it evicts the least recently used item.
func (c *Cache[Key, Value]) Set(key Key, value Value) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Calculate the size of the key and value in bytes
	entrySize := c.calculateEntrySize(key, value)

	// Check if the key already exists
	if elem, exists := c.store[key]; exists {
		// Update the value, adjust the size, and move the item to the front of the eviction list
		oldSize := elem.Value.(*entry[Key, Value]).size
		c.currSize -= oldSize   // Subtract the old size
		c.currSize += entrySize // Add the new size

		elem.Value.(*entry[Key, Value]).value = value
		elem.Value.(*entry[Key, Value]).size = entrySize
		c.evictList.MoveToFront(elem)
	} else {
		// Add new entry
		ent := &entry[Key, Value]{key: key, value: value, size: entrySize}
		elem := c.evictList.PushFront(ent)
		c.store[key] = elem
		c.currSize += entrySize
	}

	// If the cache exceeds the max memory size, evict the least recently used items
	for c.currSize > c.maxSize {
		c.evict()
	}
}

// Get retrieves a value from the cache by key and updates its LRU status.
func (c *Cache[Key, Value]) Get(key Key) (Value, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, exists := c.store[key]; exists {
		// Move the accessed element to the front of the eviction list
		c.evictList.MoveToFront(elem)
		return elem.Value.(*entry[Key, Value]).value, true
	}

	var zero Value
	return zero, false
}

// evict removes the least recently used (LRU) item from the cache.
func (c *Cache[Key, Value]) evict() {
	// Find the least recently used item, which is at the back of the list
	elem := c.evictList.Back()
	if elem == nil {
		return
	}

	// Remove the item from both the list and the map
	ent := elem.Value.(*entry[Key, Value])
	c.evictList.Remove(elem)
	delete(c.store, ent.key)

	// Adjust the current memory size
	c.currSize -= ent.size
}

// Delete removes a key from the cache.
func (c *Cache[Key, Value]) Delete(key Key) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, exists := c.store[key]; exists {
		c.evictList.Remove(elem)
		ent := elem.Value.(*entry[Key, Value])
		delete(c.store, key)
		c.currSize -= ent.size
	}
}

// calculateEntrySize estimates the memory size of a key-value pair in bytes.
// This version handles structs, maps, slices, and other composite types.
func (c *Cache[Key, Value]) calculateEntrySize(key Key, value Value) int64 {
	keySize := c.calculateSize(reflect.ValueOf(key))
	valueSize := c.calculateSize(reflect.ValueOf(value))
	return keySize + valueSize
}

// calculateSize recursively calculates the size of any Go type.
func (c *Cache[Key, Value]) calculateSize(v reflect.Value) int64 {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		// Handle pointer or interface types
		if v.IsNil() {
			return int64(unsafe.Sizeof(v.Pointer()))
		}
		return int64(unsafe.Sizeof(v.Pointer())) + c.calculateSize(v.Elem())

	case reflect.String:
		// String: Size of the string header + actual length of the string content
		return int64(unsafe.Sizeof(v.String())) + int64(len(v.String()))

	case reflect.Slice:
		// Slice: Size of the slice header + size of each element
		size := int64(unsafe.Sizeof(v.Pointer())) + int64(v.Cap())*int64(unsafe.Sizeof(v.Index(0).Interface()))
		for i := 0; i < v.Len(); i++ {
			size += c.calculateSize(v.Index(i))
		}
		return size

	case reflect.Map:
		// Map: Size of the map header + size of each key and value
		size := int64(unsafe.Sizeof(v.Pointer()))
		for _, key := range v.MapKeys() {
			size += c.calculateSize(key) + c.calculateSize(v.MapIndex(key))
		}
		return size

	case reflect.Struct:
		// Struct: Sum the size of each field
		size := int64(0)
		for i := 0; i < v.NumField(); i++ {
			size += c.calculateSize(v.Field(i))
		}
		return size

	case reflect.Array:
		// Array: Sum the size of each element
		size := int64(0)
		for i := 0; i < v.Len(); i++ {
			size += c.calculateSize(v.Index(i))
		}
		return size

	default:
		// For basic types (int, bool, etc.), just use unsafe.Sizeof
		return int64(unsafe.Sizeof(v.Interface()))
	}
}
