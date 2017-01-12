package csd

type CSD struct {
	Bin   uint64
	Signs uint64
}

func NewCSD(val int64) CSD {
	var bit uint
	var signs uint64
	var x uint64
	if val < 0 {
		x = uint64(-val)
	} else {
		x = uint64(val)
	}

	for {
		// skip zeros
		for bit < 64 && x&(1<<bit) == 0 {
			bit++
		}
		if bit == 64 {
			break
		}
		firstone := bit

		ones := 0
		for bit < 64 && x&(1<<bit) == (1<<bit) {
			bit++
			ones++
		}
		if bit == 64 {
			break
		}

		// Got next zero bit
		if ones > 1 {
			x += (1 << firstone)
			x |= (1 << firstone)
			signs |= (1 << firstone)
		}
	}

	if val < 0 {
		signs ^= x
	}

	return CSD{
		Bin:   x,
		Signs: signs,
	}
}

func (c CSD) Ones() int {
	var cnt int
	for x := c.Bin; x != 0; x >>= 1 {
		if x&1 == 1 {
			cnt++
		}
	}

	return cnt
}

func (c CSD) String() string {
	var s string
	for bit := 0; bit < 64; bit++ {

		if c.Bin&(1<<63) == 0 {
			s = s + "0"
		} else {
			if c.Signs&(1<<63) != 0 {
				s = s + "-"
			} else {
				s = s + "+"
			}
		}

		c.Bin = c.Bin << 1
		c.Signs = c.Signs << 1
	}
	return s
}

func (c CSD) Bit(n uint) bool {
	return c.Bin&(1<<n) == (1 << n)
}

func (c CSD) Sign(n uint) int {
	if c.Signs&(1<<n) == (1 << n) {
		return -1
	} else {
		return 1
	}
}

type Op struct {
	Shift int
	Sign  int
}

func (c CSD) Ops() []Op {
	ops := make([]Op, 0, 64)
	var s int
	for c.Bin != 0 {
		for c.Bin&1 == 0 {
			s++
			c.Bin >>= 1
			c.Signs >>= 1
		}

		var sign int
		if c.Signs&1 == 1 {
			sign = -1
		} else {
			sign = 1
		}

		ops = append(ops, Op{s, sign})

		s++
		c.Bin >>= 1
		c.Signs >>= 1
	}

	return ops
}
