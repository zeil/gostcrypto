package gost341112

import (
	//"crypto"
	"hash"
)

func init() {
	// TODO uncomment when integrate to golang crypto package
	//crypto.RegisterHash(crypto.GOSTR3411_2012_256, New256)
	//crypto.RegisterHash(crypto.GOSTR3411_2012_512, New)
}

// The size of a GOST R 34.11-2012 512 bit checksum in bytes.
const Size = 64

// The size of a GOST R 34.11-2012 256 bit checksum in bytes.
const Size256 = 32

// The blocksize of GOST R 34.11-2012 512 and 256 bit in bytes.
const BlockSize = 64

// digest represents the partial evaluation of a checksum.
type digest struct {
	h     [BlockSize]uint8
	N     [BlockSize]uint8
	Σ     [BlockSize]uint8
	x     [BlockSize]byte
	nx    int
	len   uint64
	is256 bool // mark if this digest is GOST R 34.11-2012 256 bit
}

func (d *digest) Reset() {
	for i := 0; i < BlockSize; i++ {
		if d.is256 {
			d.h[i] = 1
		} else {
			d.h[i] = 0
		}
		d.N[i] = 0
		d.Σ[i] = 0
		d.x[i] = 0
	}
	d.nx = 0
	d.len = 0
}

// New returns a new hash.Hash computing the GOST R 34.11-2012 512 bit checksum.
func New() hash.Hash {
	d := new(digest)
	d.Reset()
	return d
}

// New256 returns a new hash.Hash computing the GOST R 34.11-2012 256 bit checksum.
func New256() hash.Hash {
	d := new(digest)
	d.is256 = true
	d.Reset()
	return d
}

func (d *digest) Size() int {
	if !d.is256 {
		return Size
	}
	return Size256
}

func (d *digest) BlockSize() int { return BlockSize }

func (d *digest) Write(p []byte) (nn int, err error) {
	nn = len(p)
	d.len += uint64(nn)
	if d.nx > 0 {
		n := copy(d.x[d.nx:], p)
		d.nx += n
		if d.nx == BlockSize {
			block(d, d.x[:])
			d.nx = 0
		}
		p = p[n:]
	}
	if len(p) >= BlockSize {
		n := len(p) &^ (BlockSize - 1)
		block(d, p[:n])
		p = p[n:]
	}
	if len(p) > 0 {
		d.nx = copy(d.x[:], p)
	}
	return
}

func (d0 *digest) Sum(in []byte) []byte {
	// Make a copy of d0 so that caller can keep writing and summing.
	d := *d0
	hash := d.checkSum()
	if d.is256 {
		return append(in, hash[Size256:]...)
	}
	return append(in, hash[:]...)
}

func (d *digest) checkSum() [Size]byte {
	d.x[d.nx] = 1
	for i := d.nx + 1; i < BlockSize; i++ {
		d.x[i] = 0
	}
	m := &d.x
	N := &d.N
	h := &d.h
	Σ := &d.Σ
	g(N, h, m)
	addLen(N, uint64(d.nx))
	add(Σ, m[:])
	g(nil, h, N)
	g(nil, h, Σ)
	var digest [Size]byte
	copy(digest[:], d.h[:])
	return digest
}

// Sum512 returns the GOST R 34.11-2012 512 bit checksum of the data.
func Sum512(data []byte) [Size]byte {
	var d digest
	d.Reset()
	d.Write(data)
	return d.checkSum()
}

// Sum256 returns the GOST R 34.11-2012 256 bit checksum of the data.
func Sum256(data []byte) (sum256 [Size256]byte) {
	var d digest
	d.is256 = true
	d.Reset()
	d.Write(data)
	sum := d.checkSum()
	copy(sum256[:], sum[Size256:])
	return
}
