package main

import (
	"fmt"
	"log"

	"github.com/alowayed/go-univers/pkg/vers"
)

func main() {
	// Test the simple, clean VERS API
	fmt.Println("=== VERS API ===")
	
	// Simple case
	contains, err := vers.Contains("vers:maven/>=1.0.0|<=2.0.0", "1.5.0")
	if err != nil {
		log.Fatal("Error with vers.Contains:", err)
	}
	fmt.Printf("vers.Contains('vers:maven/>=1.0.0|<=2.0.0', '1.5.0') = %v\n", contains)

	// Complex case from the issue
	contains2, err := vers.Contains("vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7", "1.1.0")
	if err != nil {
		log.Fatal("Error with complex range:", err)
	}
	fmt.Printf("Complex range contains '1.1.0' = %v\n", contains2)

	// Test version outside the intervals
	contains3, err := vers.Contains("vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7", "5.0.0")
	if err != nil {
		log.Fatal("Error with complex range:", err)
	}
	fmt.Printf("Complex range contains '5.0.0' = %v\n", contains3)
}