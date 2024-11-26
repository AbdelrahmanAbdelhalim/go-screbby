package engine

type Zobrist struct {
	psq           [PIECE_NB][SQUARE_NB]Key
	enpassant     [FILE_NB]Key
	castling      [CASTLING_RIGHTS_NB]Key
	side, noPawns Key
}
type Cuckoo struct {
	cuckoo     [8192]Key
	cuckooMove [8192]Move
}

var pieces = []Piece{
	W_PAWN, W_KNIGHT, W_ROOK, W_BISHOP, W_KING, W_QUEEN,
	B_PAWN, B_KNIGHT, B_ROOK, B_BISHOP, B_KING, B_QUEEN,
}

// initialize zobrist tables for tt
func initialize(z *Zobrist, c *Cuckoo) {
	rand := Init_prng(1070372)
	for _, pc := range pieces {
		for s := SQ_A1; s <= SQ_H8; s++ {
			z.psq[pc][s] = Key(rand.Rand64())
		}
	}

	for f := FILE_A; f <= FILE_H; f++ {
		z.enpassant[f] = Key(rand.Rand64())
	}

	for cr := NO_CASTLING; cr <= ANY_CASTLING; cr++ {
		z.castling[cr] = Key(rand.Rand64())
	}

	z.side = Key(rand.Rand64())
	z.noPawns = Key(rand.Rand64())

	for i, _ := range c.cuckoo {
		c.cuckoo[i] = 0
	}

	for i, _ := range c.cuckooMove {
		c.cuckooMove[i] = 0
	}

	//get back to this
	for _, pc := range pieces {
		for s1 := SQ_A1; s1 <= SQ_H8; s1++ {
			for s2 := Square(s1 + 1); s2 <= SQ_H8; s2++ {
				if PieceType(pc) != PAWN {
					move := Move{0}
					move.SetMove(s1, s2)
				}
			}
		}
	}
}

func NewZobristCuckoo() (*Zobrist, *Cuckoo) {
	z := &Zobrist{}
	c := &Cuckoo{}
	initialize(z, c)
	return z, c
}
