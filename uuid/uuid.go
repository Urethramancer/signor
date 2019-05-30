package uuid

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type UUID struct {
	strings.Builder
	id [16]uint8
}

func NewGenerator() *UUID {
	rand.Seed(time.Now().Unix())
	return &UUID{}
}

// Generate a random UUID.
func (u *UUID) Generate() string {
	var i uint8
	for i = 0; i < 16; i++ {
		u.id[i] = uint8(rand.Intn(256))
	}

	// RFC4122 adjustment
	u.id[6] = 0x40 | (u.id[6] & 0xf)
	u.id[8] = 0x80 | (u.id[8] & 0xf3)
	x := fmt.Sprintf("%x", u.id)
	u.Reset()
	u.WriteString(x[0:8])
	u.WriteString("-")
	u.WriteString(x[8:12])
	u.WriteString("-")
	u.WriteString(x[12:16])
	u.WriteString("-")
	u.WriteString(x[16:20])
	u.WriteString("-")
	u.WriteString(x[20:])
	return u.String()
}
