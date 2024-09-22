package govault_test

import (
	"fmt"

	"github.com/drmorax/govault"
)

func ExampleNew() {

	cache := govault.New[string, []int](10)

	fmt.Printf("%d bytes which is 10 Megabytes", cache.MaxSize)
	// Output: 10485760 bytes which is 10 Megabytes
}

func ExampleCache_Set() {
	cache := govault.New[string, []int](10)

	value := []int{1, 2, 3, 4, 5}

	cache.Set("key-1", value)

	output, found := cache.Get("key-1")

	if !found {
		fmt.Println("Key doesn't exist")
	}

	fmt.Println(output)
	// Output: [1 2 3 4 5]
}

func ExampleCache_Get() {
	cache := govault.New[string, []int](10)

	cache.Set("key-1", []int{1, 2, 3, 4, 5})

	output, found := cache.Get("key-1")

	if !found {
		fmt.Println("Key doesn't exist")
	}

	fmt.Println(output)
	// Output: [1 2 3 4 5]
}
func ExampleCache_Delete() {
	cache := govault.New[string, []int](10)

	cache.Set("key-1", []int{1, 2, 3, 4, 5})

	output, found := cache.Get("key-1")
	fmt.Print("\nkey-1: ", output, found)

	cache.Delete("key-1")

	output, found = cache.Get("key-1")
	fmt.Print("\nkey-1: ", output, found)
	// Output:
	//key-1: [1 2 3 4 5] true
	//key-1: [] false
}
