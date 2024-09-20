package main

import (
	"fmt"

	"github.com/drmorax/govault"
)

func main() {
	// Create a cache with a memory limit of 1 MB
	cache := govault.New // 1 MB limit

	// Add some entries (assuming small values here for simplicity)
	cache.Set("a", "value_a")
	cache.Set("b", "value_b")
	cache.Set("c", "value_c")

	// Check entries
	fmt.Println("Found a:", cache.Get("a"))
	fmt.Println("Found b:", cache.Get("b"))

	// Add a large entry to trigger eviction (assuming large size for example)
	cache.Set("d", "this_is_a_large_value")

	// Check if least recently used entry ("a") has been evicted
	_, found := cache.Get("a")
	fmt.Println("Found a after eviction:", found)
}
