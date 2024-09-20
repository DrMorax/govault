package main

import (
	"fmt"

	"github.com/DrMorax/govault"
)

func main() {
	// Create a cache with a memory limit of 3 MB
	// Usage of New: New[keyType(comparable values), valueType(anything, e.g. string, int, struct, []byte, etc...)]
	cache := govault.New[string, string](3) // 3 MB limit

	// Add some entries (assuming small values here for simplicity)
	cache.Set("a", oneMBText)
	cache.Set("b", oneMBText)
	cache.Set("c", oneMBText)

	_, found1 := cache.Get("a")
	_, found2 := cache.Get("b")
	_, found3 := cache.Get("c")
	// Check entries
	fmt.Println("a:", found1)
	fmt.Println("b:", found2)
	fmt.Println("c:", found3)

	// Check if least recently used entry ("a") has been evicted
	cache.Set("d", oneMBText)
	_, found4 := cache.Get("d")
	fmt.Println("d:", found4)

	_, found5 := cache.Get("a")
	fmt.Println("a:", found5)
}
