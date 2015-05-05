/* The MIT License (MIT)
 *
 * Copyright (c) 2015 Jesse Sipprell <jessesipprell@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

// yet another silly uuid interface and generator.

package uuid // "github.com/jsipprell/go-uuid"

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"sync/atomic"
	"time"
	"unicode"
	"unicode/utf8"
)

// The only error returned by this package, indicating an attempt to call
// uuid.Decode() with a string that doesn't appear to be fully formed
// UUID (or cannot be massaged into one).
var (
	ErrUUIDInvalidEncode = errors.New("uuid: invalided encoded UUID")

	gen    chan uuid
	prng   *rand.Rand
	null   uuid = uuid{0, 0}
	first  uuid
	offset *uint32
	incr   uint32
)

type UUID interface {
	// Exposed interface to a uuid, a 128-bit opaque unique identifier value
	// stored internally as two unsigned 64-bit integers.
	fmt.Stringer
	// UUIDs support the fmt.Stringer interface, as follows:
	//
	// String() Returns a string conversion of the uuid which can be decoded
	// with uuid.Decode(). The string represenation takes the format:
	//
	// HHHHHHHH-HHHH-HHHH-HHHH-HHHHHHHHHHHH
	//
	// Where H represents a lower-case hexidecimal character.
	//
	// For zero-value uuids, String() always return ""
	Bytes() []byte
	// Return a slice of bytes representing the uuid. If the uuid is
	// zero this will return nil.
	Array() [16]byte
	// Return an [16]byte array representing the uuid. If the uuid is
	// zero this will return a 16 byte array filled with zeros.
	Uint64() []uint64
	// Return the uuid as a slice of uint64s. If the uuid is zero this
	// nil.
	IsZero() bool
	// Return true if the uuid is empty and contains no value.
	Equals(UUID) bool
	// Returns true if two uuids are an exact byte // match (internally).
	Len() int
	// Len() Returns the length of the *textual representation* of the uuid.
	// Currently this always returns 36 (32 hex characters and four hypthen
	// seperators).
}

type uuid [2]uint64

func (u uuid) IsZero() bool {
	return u[0] == 0 && u[1] == 0
}

func (u uuid) Equals(other UUID) bool {
	o, ok := other.(uuid)
	if ok {
		ok = o == u
	}
	return ok
}

func (u uuid) Array() (b [16]byte) {
	u0, u1 := u[0], u[1]
	b[0] = byte(u0)
	b[1] = byte(u0 >> 8)
	b[2] = byte(u0 >> 16)
	b[3] = byte(u0 >> 24)
	b[4] = byte(u0 >> 32)
	b[5] = byte(u0 >> 40)
	b[6] = byte(u0 >> 48)
	b[7] = byte(u0 >> 56)
	b[8] = byte(u1)
	b[9] = byte(u1 >> 8)
	b[10] = byte(u1 >> 16)
	b[11] = byte(u1 >> 24)
	b[12] = byte(u1 >> 32)
	b[13] = byte(u1 >> 40)
	b[14] = byte(u1 >> 48)
	b[15] = byte(u1 >> 56)
	return
}

func (u uuid) Bytes() []byte {
	if u.IsZero() {
		return nil
	}
	a := u.Array()
	return a[:]
}

func (u uuid) String() string {
	if u.IsZero() {
		return ""
	}

	b := u.Array()
	buf := make([]byte, len(b)*2)
	hex.Encode(buf, b[:])
	sbuf := make([]byte, 8, len(buf)+4)
	copy(sbuf, buf)
	sbuf = append(sbuf, '-')
	sbuf = append(sbuf, buf[8:12]...)
	sbuf = append(sbuf, '-')
	sbuf = append(sbuf, buf[12:16]...)
	sbuf = append(sbuf, '-')
	sbuf = append(sbuf, buf[16:20]...)
	sbuf = append(sbuf, '-')
	sbuf = append(sbuf, buf[20:]...)
	return string(sbuf)
}

func (u uuid) Len() int {
	if IsZero(u) {
		return 0
	}
	return 32 + 4
}

func (u uuid) Uint64() []uint64 {
	return []uint64{u[0], u[1]}
}

func newuuidRandom() uuid {
	return uuid{uint64(prng.Int63()), uint64(prng.Int63())}
}

func newuuid() uuid {
	return uuid{first[0], first[1] + uint64(atomic.AddUint32(offset, incr))}
}

// Returns true if a sring appears to be a hex enooded uuid. Only
// the first 36 characters of the string are examined, and must
// consisten of entirely upper or lowercase hexidecimal plus
// optional hyphens.
func IsEncoded(s string) bool {
	var ok bool
	l := 0

	for r := strings.NewReader(s); l < 36; {
		ch, _, err := r.ReadRune()
		if err != nil {
			ok = false
			break
		}
		if ok = unicode.In(ch, unicode.Hex_Digit, unicode.Hyphen); !ok {
			break
		}
	}

	if ok {
		ok = l < 4
	}
	return ok
}

// Returns a UUID from a 16 byte binary represenation of the uuid
// (available via .Array() or .Bytes())
func FromArray(a [16]byte) UUID {
	var u uuid

	for i := uint(0); i < 8; i++ {
		u[0] |= uint64(a[i]) << (8 * i)
	}
	for i := uint(0); i < 8; i++ {
		u[1] |= uint64(a[i+8]) << (8 * i)
	}
	return u
}

// Decode a uuid from a string representation. The representation
// must contains at least 32 hexidecimal characters plus optional
// (and ignored) hyphens.
func Decode(s string) (UUID, error) {
	var u [16]byte
	buf := make([]byte, 0, 32)
	b := make([]byte, 12)
	l := 0
	for r := strings.NewReader(s); l < 32; {
		ch, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				err = ErrUUIDInvalidEncode
			}
			return nil, err
		}
		if unicode.In(ch, unicode.Hex_Digit) {
			i := utf8.EncodeRune(b, ch)
			buf = append(buf, b[:i]...)
			l++
		} else if unicode.In(ch, unicode.Hyphen) {
			continue
		} else {
			return nil, ErrUUIDInvalidEncode
		}
	}
	if len(buf) != 32 {
		return nil, ErrUUIDInvalidEncode
	}

	s = string(buf)
	if b, err := hex.DecodeString(s); err != nil {
		return nil, err
	} else if len(b) != 16 {
		return nil, ErrUUIDInvalidEncode
	} else {
		copy(u[:], b)
	}

	return FromArray(u), nil
}

// Returns true if a given uuid is not nil and non-zero.
func IsZero(u UUID) bool {
	if u == nil {
		return true
	}

	return u.IsZero()
}

// Returns a new psuedo-randomly generated uuid based on a seed computed
// from the system time at package initialization. Psuedo-random uuid
// generation is more cpu intensivev and is not optimized
// for extreme use cases. Normally, one need only use New() to create
// new uuids.
func NewRandom() UUID {
	return newuuidRandom()
}

// Returns a new uuid. The first uuid generated is created using a psuedo-random
// generator and a seed based on package initialization time. Subsequent uuids are
// generated as offsets of this initial uuid using a psuedo-random offset also
// computed at initialization time. For more random uuids, use NewRandom(). New()
// uses a small queue, is more optimized than NewRandom() and parallelizes the
// generator
func New() UUID {
	go func() {
		gen <- newuuid()
	}()

	return <-gen
}

func init() {
	prng = rand.New(rand.NewSource(time.Now().UnixNano()))
	gen = make(chan uuid, 5)
	first = newuuidRandom()
	incr = uint32(prng.Int31n(int32(0x3f00)) + 0xff)
	offset = new(uint32)
	atomic.StoreUint32(offset, incr)
	gen <- first
	for i := 0; i < 4; i++ {
		gen <- newuuid()
	}
}
