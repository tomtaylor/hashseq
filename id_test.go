package hashseq_test

import (
	"encoding/json"
	"testing"

	"github.com/tomtaylor/hashseq"
)

func init() {
	hashseq.SetSalt("testing")
}

func TestEncode(t *testing.T) {
	i := &hashseq.ID{Int64: 101}

	s := i.String()
	if s != "5exA" {
		t.Errorf("unexpected string: %s", s)
	}
}

func TestBadEncode(t *testing.T) {
	i := &hashseq.ID{Int64: -1}
	s := i.String()
	if s != "" {
		t.Errorf("expecting encode to be blank")
	}
}

func TestJSONEncode(t *testing.T) {
	type blob struct {
		Key hashseq.ID `json:"key"`
	}

	foo := &blob{Key: hashseq.ID{Int64: 101}}
	b, err := json.Marshal(foo)
	if err != nil {
		t.Error(err)
	}

	s := string(b)
	if s != `{"key":"5exA"}` {
		t.Errorf("encoded JSON did not match, got: %s", s)
	}
}

func TestJSONDecode(t *testing.T) {
	type testStruct struct {
		Key hashseq.ID `json:"key"`
	}

	s := testStruct{}

	b := []byte(`{"key":"5exA"}`)
	err := json.Unmarshal(b, &s)
	if err != nil {
		t.Error(err)
	}

	if s.Key.Int64 != 101 {
		t.Errorf("key was not expected, got %v", s.Key.Int64)
	}
}
func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		i := &hashseq.ID{Int64: int64(i)}
		_ = i.String()
	}
}
