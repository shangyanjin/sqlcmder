package main

import (
	"fmt"
	"os"
	"time"
)

// Test error display functionality
func main() {
	fmt.Println("Testing error display...")

	// Simulate ShowError being called
	testShowError()

	fmt.Println("Test completed - check if process exits cleanly")
	time.Sleep(2 * time.Second)
}

func testShowError() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIC caught: %v\n", r)
			os.Exit(1)
		}
	}()

	fmt.Println("Calling ShowError simulation...")

	// This simulates what happens when ShowError is called
	var cmdLine interface{} = nil
	if cmdLine != nil {
		fmt.Println("Command line is available")
	} else {
		fmt.Println("ERROR: globalCommandLine is nil!")
		fmt.Println("This is the bug - ShowError tries to use nil globalCommandLine")
	}

	fmt.Println("ShowError simulation completed")
}
