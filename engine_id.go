package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	format_masc uint64 = 0x8000_0000
	vendor_masc uint64 = 0x7fff_ffff
	id_masc     uint64 = 0xff_ffff_ffff_ffff
)

type vendor struct {
	number  int
	name    string
	contact string
	email   string
}

type DeviceID interface {
	description() string
}

type mac struct {
	mac    string
	vendor string
}

func (this *mac) description() string {
	return fmt.Sprintf("Mac Adress Engine ID:\n\tMAC: %v\n\tVendor: %v", this.mac, this.vendor)
}

type ipv4 struct {
	ip string
}

func (this *ipv4) description() string {
	return fmt.Sprintf("IPv4 Engine ID:\n\tIP: %v", this.ip)
}

type ipv6 struct {
	ip string
}

func (this *ipv6) description() string {
	return fmt.Sprintf("IPv6 Engine ID:\n\tIP: %v", this.ip)
}

type text struct {
	id string
}

func (this *text) description() string {
	return fmt.Sprintf("Text Engine ID:\n\tID: %v", this.id)
}

type bytes struct {
	id []int8
}

func (this *bytes) description() string {
	return fmt.Sprintf("Octet Engine ID:\n\t0x%v", this.id)
}

type EngineId interface {
	vendorInfo() *vendor
	customFormat() bool
	deviceID() DeviceID
}

type CustomFormat struct {
	vendor vendor
	id     text
}

func (this *CustomFormat) vendorInfo() *vendor {
	return &this.vendor
}

func (this CustomFormat) customFormat() bool  { return true }
func (this *CustomFormat) deviceID() DeviceID { return &this.id }

type ProperFormat struct {
	vendor vendor
	iana   int
	id     DeviceID
}

func load_enterprise_numbers(path string) ([]vendor, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	recs, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	ens := []vendor{}
	for _, rec := range recs[1:] {
		num, err := strconv.Atoi(rec[0])
		if err != nil {
			return nil, err
		}
		ens = append(ens, vendor{
			num,
			rec[1],
			rec[2],
			rec[3],
		})
	}
	return ens, nil
}

type MacBlock struct {
	prefix     string
	name       string
	blockType  string
	lastUpdate string
	private    bool
}

func load_mac_blocks(path string) ([]MacBlock, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	type TmpBlock struct {
		// names are ones that reflect the json we downloaded
		MacPrefix  string
		VendorName string
		BlockType  string
		LastUpdate string
		private    bool
	}
	dec := json.NewDecoder(file)
	tmp := []TmpBlock{}
	dec_err := dec.Decode(&tmp)
	if dec_err != nil {
		return nil, dec_err
	}
	ret := []MacBlock{}
	for _, block := range tmp {
		shortened := strings.ReplaceAll(block.MacPrefix, ":", "")
		ret = append(ret, MacBlock{
			shortened,
			block.VendorName,
			block.BlockType,
			block.LastUpdate,
			block.private,
		})
	}
	return ret, nil
}

func decode_engine_id(id string, ens []vendor) (EngineId, error) {
	if len(id) < 12 {
		return nil, New(fmt.Sprintf("Incorrect ID string length (expected at least 12, but got %d)", len(id)))
	}
	vendor_id, err := strconv.ParseUint(id[:8], 16, 64)
	if err != nil {
		return nil, err
	}
	id_int, err := strconv.ParseUint(id[8:], 16, 64)
	if err != nil {
		return nil, err
	}
	return decode_engine_id_num(vendor_id, id_int, ens)
}

func decode_engine_id_num(vendor_id uint64, id uint64, ens []vendor) (EngineId, error) {
	vendor_obj := ens[vendor_masc&vendor_id]
	// if id & format_masc != 0 {
	// 	// engineID follows format
	// } else {
	// 	return CustomFormat{vendor_obj, bytes{[0]}}
	// }
	ret := CustomFormat{vendor_obj, text{strconv.FormatUint(id, 16)}}
	return &ret, nil
}
