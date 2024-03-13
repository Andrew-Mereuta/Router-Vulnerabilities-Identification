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

func TestEngineIDParse(t *testing.T) {
	e_nums, err := load_enterprise_numbers("data/enterprise-numbers")
	if err != nil {
		t.Fatalf("Got an error: %v", err)
	}
	engID, err := decode_engine_id("800007c703748ef82a4900", e_nums)
	if err != nil {
		t.Fatalf("Got an error: %v", err)
	}
	// check vendor
	if engID.vendorInfo().number != 1991 {
		t.Fatalf("IANA PEN not decoded correctly, expected: 1991, but got: %d.\n", engID.vendorInfo().number)
	}
	if engID.vendorInfo().name != "Brocade Communication Systems, Inc. (formerly 'Foundry Networks, Inc.')" {
		t.Fatalf("IANA PEN not decoded correctly, expected: \"Brocade Communication Systems...\", but got: %v.\n", engID.vendorInfo().name)
	}
	if engID.vendorInfo().contact != "Scott Kipp" {
		t.Fatalf("IANA PEN not decoded correctly, expected: \"Scott Kipp\", but got: %v.\n", engID.vendorInfo().contact)
	}
	// NOTE: the data I got uses `&` instead of `@`
	if engID.vendorInfo().email != "skipp&brocade.com" {
		t.Fatalf("IANA PEN not decoded correctly, expected: \"skipp@brocade.com\", but got: %v.\n", engID.vendorInfo().email)
	}
	// TODO: check rest of the data
}

func TestLoadMacBlock(t *testing.T) {
	_, err := load_mac_blocks("data/mac-vendors-export.json")
	if err != nil {
		t.Fatalf("Got an error: %v", err)
	}
}
