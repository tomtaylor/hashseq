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

	hashid "github.com/speps/go-hashids"
)

var HashData = &hashid.HashIDData{
	Alphabet:  hashid.DefaultAlphabet,
	MinLength: 4,
	Salt:      "",
}

// Set the salt to use for ID obfuscation
func SetSalt(salt string) {
	HashData.Salt = salt
}

type Id struct {
	Int64 int64
}

func (id *Id) Int() int {
	return int(id.Int64)
}

// Return the hashid as an obfuscated string
func (id *Id) String() string {
	str, err := hashid.NewWithData(HashData).Encode([]int{id.Int()})
	if err != nil {
		return ""
	}
	return str
}

// Returns the hashid as a string fulfilling the json.Marshaller interface
func (id Id) MarshalJSON() ([]byte, error) {
	str, err := hashid.NewWithData(HashData).Encode([]int{id.Int()})
	if err != nil {
		return nil, err
	}
	return json.Marshal(str)
}

// Unmarshal a string and decode into an integer
func (id *Id) UnmarshalJSON(data []byte) error {
	decoded, err := Decode(data)
	if err != nil {
		return err
	}
	id.Int64 = int64(decoded)
	return nil
}

// Decode a hashid byte into an Id, setting its integer
func Decode(hashid []byte) (id int64, err error) {
	return DecodeString(string(hashid))
}

func MustDecodeString(hashid string) int64 {
	i, err := DecodeString(hashid)
	if err != nil {
		panic(err.Error())
	}
	return int64(i)
}

// Decode a hashid string into an Id, setting its integer
func DecodeString(h string) (id int64, err error) {
	ints := hashid.NewWithData(HashData).Decode(h)
	if len(ints) != 1 {
		err = fmt.Errorf("Unexpected hashid value")
		return
	}
	return int64(ints[0]), nil
}

// Database scanning
func (id *Id) Scan(value interface{}) (err error) {
	var data int64

	// If the first four bytes of this are 0000
	switch value.(type) {
	// Same as []byte
	case int64:
		data = value.(int64)
	case nil:
		return
	default:
		return fmt.Errorf("Invalid format: can't convert %T into id.Id", value)
	}

	id.Int64 = data
	return nil
}

// This is called when saving the ID to a database
func (id Id) Value() (driver.Value, error) {
	return id.Int64, nil
}
