package main

import (
	"fmt"
	"testing"
)

func TestLoadEnterprise(t *testing.T) {
	e_nums, err := load_enterprise_numbers("data/enterprise-numbers")
	if err != nil {
		t.Fatalf("Got an error: %v", err)
	}
	if len(e_nums) != 61612 {
		t.Fatalf("Didn't get enough records: %d", len(e_nums))
	}
	// assert that the loaded list is sorted by their id (for fast lookup)
	for i, num := range e_nums {
		if i != num.number {
			fmt.Printf("Entry %d did not match id %d", i, num.number)
			t.Fail()
		}
	}
}
