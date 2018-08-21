// This package provides ID masking using hash ids.
//
// Hash IDs are a string represnetation of numerical incrementing IDs,
// obfuscating the integer value. For more information see
// http://hashids.org/go/
package hashseq

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	hashids "github.com/speps/go-hashids"
)

var (
	globalHashID *hashids.HashID
)

func init() {
	setGlobalHashID("")
}

func setGlobalHashID(salt string) {
	hashID, err := hashids.NewWithData(&hashids.HashIDData{
		Alphabet:  hashids.DefaultAlphabet,
		MinLength: 4,
		Salt:      salt,
	})

	if err != nil {
		panic(err)
	}

	globalHashID = hashID
}

// SetSalt sets the salt used for ID obfuscation
func SetSalt(salt string) {
	setGlobalHashID(salt)
}

// ID is a struct containing the hashid, stored as an int64
type ID struct {
	Int64 int64
}

// Int returns the ID as an int
func (id *ID) Int() int {
	return int(id.Int64)
}

// String returns the ID as an obfuscated string
func (id *ID) String() string {
	str, err := globalHashID.Encode([]int{id.Int()})
	if err != nil {
		return ""
	}
	return str
}

// MarshalJSON returns the ID as a string fulfilling the json.Marshaller interface
func (id ID) MarshalJSON() ([]byte, error) {
	str, err := globalHashID.Encode([]int{id.Int()})
	if err != nil {
		return nil, err
	}
	return json.Marshal(str)
}

// UnmarshalJSON turns a string into an ID
func (id *ID) UnmarshalJSON(data []byte) error {
	input := string(data)
	input = strings.Trim(input, `"`)

	decoded, err := DecodeString(input)
	if err != nil {
		return err
	}
	id.Int64 = int64(decoded)
	return nil
}

// Decode a hashid byte into an ID, by setting its integer
func Decode(hashid []byte) (id int64, err error) {
	return DecodeString(string(hashid))
}

// MustDecodeString panics if decode failed
func MustDecodeString(hashid string) int64 {
	i, err := DecodeString(hashid)
	if err != nil {
		panic(err.Error())
	}
	return int64(i)
}

// DecodeString a hashid string into an ID, by setting its integer
func DecodeString(h string) (id int64, err error) {
	ints := globalHashID.Decode(h)
	if len(ints) != 1 {
		err = fmt.Errorf("Unexpected hashid value")
		return
	}
	return int64(ints[0]), nil
}

// Scan implements the driver.Scanner interface for converting from database
func (id *ID) Scan(value interface{}) (err error) {
	var data int64

	// If the first four bytes of this are 0000
	switch value.(type) {
	// Same as []byte
	case int64:
		data = value.(int64)
	case nil:
		return
	default:
		return fmt.Errorf("Invalid format: can't convert %T into id.Int64", value)
	}

	id.Int64 = data
	return nil
}

// Value implements the driver.Valuer interface, for converting to a database
func (id ID) Value() (driver.Value, error) {
	return id.Int64, nil
}
