package bertlv

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	errIndefiniteLength = errors.New("intefinite length is not supported")
	errInvalidLength    = errors.New("invalid length")
)

// TagValue pair
type TagValue struct {
	T int
	V []byte
}

// Bytes encodes TagValue into a byte slice
func (tv TagValue) Bytes() []byte {
	result := append(tv.Tag(), tv.Len()...)
	return append(result, tv.V...)
}

func (tv TagValue) String() string {
	b, err := json.Marshal(tv)
	if err != nil {
		return err.Error()
	}

	return string(b)
}

// Tag returns encoded tag value (two bytes tags are supported)
func (tv TagValue) Tag() []byte {
	if (tv.T>>8)&0x1F == 0 {
		return []byte{byte(tv.T)}
	}

	return []byte{byte(tv.T >> 8), byte(tv.T & 0xff)}
}

// Len returns encoded length of the value
func (tv TagValue) Len() []byte {
	l := len(tv.V)
	if l <= 0x7f { //the first byte is a final byte?
		return []byte{byte(l)}
	}

	r := EncodeInt(l)
	numOctets := len(r)
	result := make([]byte, 1+numOctets)
	result[0] = 0x80 | byte(numOctets)

	copy(result[1:], r)

	return result
}

// WriteTo implements io.WriterTo
func (tv TagValue) WriteTo(w io.Writer) (n int64, err error) {
	tn, err := w.Write(tv.Bytes())
	return int64(tn), err
}

// ReadFrom implements io.ReaderFrom
func (tv *TagValue) ReadFrom(r io.Reader) (n int64, err error) {
	tag, tagn, err := ReadTag(r)
	if err != nil {
		return int64(tagn), err
	}

	l, ln, err := ReadLen(r)
	if err != nil {
		return int64(tagn) + int64(ln), fmt.Errorf("failed to read length: %v", err)
	}

	tv.T = tag

	if l == 0 {
		return int64(tagn) + int64(ln), nil
	}

	tv.V = make([]byte, l)
	vn, err := r.Read(tv.V)
	if vn < l {
		return int64(tagn) + int64(ln) + int64(vn), io.ErrUnexpectedEOF
	}

	if err != nil {
		return int64(tagn) + int64(ln) + int64(vn), fmt.Errorf("failed to read value: %v", err)
	}

	return int64(tagn) + int64(ln) + int64(vn), nil
}

// IsConstructed returns true if the value is constructed type(contains other TLV records)
func (tv TagValue) IsConstructed() bool {
	if tv.T <= 0xff {
		return tv.T&0x20 != 0
	}

	return (tv.T>>8)&0x20 != 0
}

// MarshalJSON implements json.Marshaller
func (tv TagValue) MarshalJSON() ([]byte, error) {
	//Contents of the tag is not a TLV structure
	if !tv.IsConstructed() {
		return []byte(fmt.Sprintf("{\"%x\":\"%x\"}", tv.T, tv.V)), nil
	}

	tvs, err := Decode(tv.V)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(tvs)
	if err != nil {
		return nil, err
	}

	result := append([]byte(fmt.Sprintf("{\"%x\":", tv.T)), b...)
	result = append(result, []byte("}")...)
	return result, nil
}

// ReadTag reads length of a tag and return it along with length of a length in bytes or an error
func ReadTag(r io.Reader) (tag int, n int, err error) {
	b := make([]byte, 1)

	//reading first byte of the tag
	n, err = r.Read(b)
	if err != nil {
		return 0, n, err
	}

	tag = int(b[0])

	if b[0]&0x1F == 0x1F { //it's a two byte tag
		tag <<= 8

		n, err = r.Read(b)
		if err != nil {
			return 0, n + 1, err
		}

		n = 2
		tag |= int(b[0])
	}

	return tag, n, nil
}

// ReadLen reads length of a tag and return it along with length of a length in bytes and/or an error
func ReadLen(r io.Reader) (length int, n int, err error) {
	b := make([]byte, 1)

	n, err = r.Read(b)
	if err != nil {
		return 0, n, err
	}

	if b[0] == 0x80 {
		return 0, 1, errIndefiniteLength
	}

	if b[0]&0x80 == 0 {
		return int(b[0]), 1, nil
	}

	nb := int(b[0] & 0x7f)
	if nb > 4 {
		return 0, 1, errInvalidLength
	}

	lenb := make([]byte, 4)
	n, err = r.Read(lenb[4-nb:])
	if err != nil {
		return 0, n + 1, err
	}

	return int(binary.BigEndian.Uint32(lenb)), n + 1, nil
}

// EncodeInt encodes an integer to BER format.
func EncodeInt(in int) []byte {
	result := make([]byte, 4)

	binary.BigEndian.PutUint32(result, uint32(in))

	var lz int
	for ; lz < 4; lz++ {
		if result[lz] != 0 {
			break
		}
	}

	return result[lz:]
}

// Decode decodes TLV encoded byte slice into slice of TagValue structs
func Decode(p []byte) ([]TagValue, error) {
	r := bytes.NewReader(p)

	var tv TagValue
	var result []TagValue
	for {
		_, err := tv.ReadFrom(r)
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		result = append(result, tv)
	}

	return result, nil
}

// ErrNotFound is returned when requested tag is not present in TLV structure
var ErrNotFound = errors.New("not found")

// Find finds first tag (DFS) in the TLV structure represented by p
func Find(tag int, p []byte) (*TagValue, error) {
	tvs, err := Decode(p)
	if err != nil {
		return nil, err
	}

	for _, tv := range tvs {
		if tv.T == tag {
			return &tv, nil
		}

		if !tv.IsConstructed() {
			continue
		}

		tvv, err := Find(tag, tv.V)
		if err != nil {
			return nil, err
		}

		if tvv != nil {
			return tvv, nil
		}
	}

	return nil, ErrNotFound
}
