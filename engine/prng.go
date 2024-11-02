package engine

type PRNG struct {
	s uint64
}

// Returns a random value of type uint64
func (prng *PRNG) Rand64() uint64 {
	prng.s ^= prng.s >> 12
	prng.s ^= prng.s << 25
	prng.s ^= prng.s >> 27
	return prng.s * 2685821657736338717
}

func (prng *PRNG) Sparse_rand() uint64 {
	return prng.Rand64() & prng.Rand64() & prng.Rand64()
}

func Init_prng(seed uint64) *PRNG {
	return &PRNG{s: seed}
}
