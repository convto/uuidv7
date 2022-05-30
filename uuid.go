// uuidv7 implementation proposed in draft https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-03
package uuidv7

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

// UUID is an array type to represent the value of a UUID, as defined in RFC-4122.
type UUID [16]byte

var Nil = UUID{}

// UUID versions.
const (
	_  byte = iota
	_  // Version 1
	_  // Version 2
	_  // Version 3
	_  // Version 4
	_  // Version 5
	_  // Version 6
	V7 // Version 7 (k-sortable timestamp, and random data) [peabody draft]
	_  // Version 8
)

// UUID layout variants.
const (
	VariantNCS byte = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

func New() (UUID, error) {
	var uuid UUID

	var tms uint64
	now := time.Now()
	tms += uint64(now.Unix())*1e3
	tms += uint64(now.Nanosecond())/1e6
	binary.BigEndian.PutUint64(uuid[:8], tms<<16)

	_, err := io.ReadFull(rand.Reader, uuid[6:])
	if err != nil {
		return Nil, err
	}

	uuid[6] = (uuid[6] & 0x0f) | V7             // Version 7
	uuid[8] = (uuid[8] & 0x3f) | VariantRFC4122 // Variant is 10

	return uuid, nil
}

// Parse UUID strings that are formatted as defined in RFC-4122 (section 3):
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func Parse(s string) (UUID, error) {
	var uuid UUID
	if len(s) != 36 {
		return uuid, fmt.Errorf("invalid UUID length: %d", len(s))
	}
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return uuid, errors.New("invalid UUID format")
	}
	// Processed two hex characters at a time
	for i, x := range [16]int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34} {
		v := (s[x] << 4) | s[x+1]
		uuid[i] = v
	}
	return uuid, nil
}

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (uuid UUID) String() string {
	var buf [36]byte
	hex.Encode(buf[:8], uuid[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], uuid[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], uuid[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], uuid[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], uuid[10:])
	return string(buf[:])
}