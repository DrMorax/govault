package main

import (
	"fmt"

	"github.com/drmorax/govault"
)

func main() {
	// Create a cache with a memory limit of 3 MB
	cache := govault.New[string, string](3) // 3 MB limit

	// Add some entries (assuming small values here for simplicity)
	cache.Set("a", oneMBText)
	cache.Set("b", oneMBText)
	cache.Set("c", oneMBText)
	cache.Set("d", oneMBText)

	_, found1 := cache.Get("a")
	_, found2 := cache.Get("b")
	_, found3 := cache.Get("c")
	_, found4 := cache.Get("d")
	// Check entries
	fmt.Println("a:", found1)
	fmt.Println("b:", found2)
	fmt.Println("c:", found3)
	fmt.Println("d:", found4)

	// Check if least recently used entry ("a") has been evicted
	_, found5 := cache.Get("a")
	fmt.Println("a:", found5)
}
