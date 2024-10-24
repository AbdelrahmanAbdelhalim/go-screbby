package engine

import (
	"math/bits"
	"strconv"
	"strings"
)

// BitBoards For Each File

const FileABB BitBoard = 0x0101010101010101
const FileBBB BitBoard = FileABB << 1
const FileCBB BitBoard = FileABB << 2
const FileDBB BitBoard = FileABB << 3
const FileEBB BitBoard = FileABB << 4
const FileFBB BitBoard = FileABB << 5
const FileGBB BitBoard = FileABB << 6
const FileHBB BitBoard = FileABB << 7

// BitBoards for Each Rank
const Rank1BB BitBoard = 0xFF
const Rank2BB BitBoard = 0xFF << (8 * 1)
const Rank3BB BitBoard = 0xFF << (8 * 2)
const Rank4BB BitBoard = 0xFF << (8 * 3)
const Rank5BB BitBoard = 0xFF << (8 * 4)
const Rank6BB BitBoard = 0xFF << (8 * 5)
const Rank7BB BitBoard = 0xFF << (8 * 6)
const Rank8BB BitBoard = 0xFF << (8 * 7)

var popCnt16 = make([]uint8, 65535)
var pawnAttacks [COLOR_NB][SQUARE_NB]BitBoard
var squareDistance [SQUARE_NB][SQUARE_NB]uint8
var pseudoAttacks [PIECE_NB][SQUARE_NB]BitBoard
var PawnAttacks [COLOR_NB][SQUARE_NB]BitBoard
var lineBB [SQUARE_NB][SQUARE_NB]BitBoard
var betweenBB [SQUARE_NB][SQUARE_NB]BitBoard

func printbb(b BitBoard) string {
	var builder strings.Builder
	builder.WriteString("+---+---+---+---+---+---+---+---+\n")
	for r := RANK_8; r >= RANK_1; r -= 1 {
		for f := FILE_A; f <= FILE_H; f += 1 {
			if bitboard_and_square(b, make_square(f, r)) != 0 {
				builder.WriteString("| X ")
			} else {
				builder.WriteString("|   ")
			}
		}
		builder.WriteString("| ")
		builder.WriteString(" " + strconv.Itoa(int((1 + r))))
		builder.WriteString("\n+---+---+---+---+---+---+---+---+\n")
	}
	builder.WriteString("  a   b   c   d   e   f   g   h\n")
	return builder.String()
}

func (bb BitBoard) calc_pawn_attacks_bb(c Color) BitBoard {
	if c == WHITE {
		return bb.shift(NORTH_EAST) | bb.shift(NORTH_WEST)
	} else {
		return bb.shift(SOUTH_EAST) | bb.shift(SOUTH_WEST)
	}
}

func (bb BitBoard) popcount() int {
	return bits.OnesCount64(uint64(bb))
}
func (bb BitBoard) msb() int {
	if bb == 0 {
		return -1
	}
	return 63 - bits.LeadingZeros64(uint64(bb))
}

var RookTable [0x19000]BitBoard
var BishopTable [0x1490]BitBoard
var RookMagics [SQUARE_NB]Magic
var BishopMagics [SQUARE_NB]Magic

type Magic struct {
	mask    BitBoard
	magic   BitBoard
	attacks []BitBoard
	shift   uint64
}
type Magics struct {
	RookMagics   [SQUARE_NB]Magic
	BishopMagics [SQUARE_NB]Magic
}

// todo: Optimize this
func (m Magic) index(occupied BitBoard) uint64 {
	//todo: implement haspext directive
	lo := uint64(((occupied & m.mask) * m.magic) >> m.shift)
	hi := uint64(occupied>>32) & uint64(m.magic>>32) >> m.shift
	return (lo*uint64(m.magic) ^ hi*uint64(m.magic>>32)) >> m.shift
}

func distance(s1 Square, s2 Square) uint8 {
	return squareDistance[s1][s2]
}

func init_magics(pt PieceType, table []BitBoard, magics []Magic) {
	var seeds = [][]int32{{8977, 44560, 54343, 38998, 5731, 95205, 104912, 17020},
		{728, 10316, 55013, 32803, 12281, 15100, 16645, 255}}
	var occupancy [4096]BitBoard
	var reference [4096]BitBoard
	var b BitBoard
	var epoch = [4096]int32{}
	var cnt, size int32

	for s := SQ_A1; s <= SQ_H8; s++ {
		edges := ((Rank1BB | Rank8BB) & ^rank_bb(s.rank_of())) | ((FileABB | FileHBB) & ^file_bb(s.file_of()))
		m := &magics[s]
		m.mask = sliding_attacks(pt, s, 0) & ^edges
		m.shift = uint64(64 - m.mask.popcount()) //todo: Take into account 32bit vs 64 bit

		b = 0
		size = 0
		runonce := false
		//todo: sort out has pext
		for !runonce || (b != 0) {
			occupancy[size] = b
			reference[size] = sliding_attacks(pt, s, b)
			size++
			b = (b - m.mask) & m.mask
			runonce = true
		}
		rng := Init_prng(uint64(seeds[1][s.rank_of()])) //todo: adjust for 32 bit vs 64 bit

		for i := 0; i < int(size); {
			for m.magic = 0; ((m.magic * m.mask) >> 56).popcount() == 0; {
				m.magic = BitBoard(rng.Sparse_rand())
			}
			for i := 0; i < int(size); i++ {
				cnt++
				idx := m.index(occupancy[i])
				if epoch[idx] < cnt {
					epoch[idx] = cnt
					m.attacks[idx] = reference[i]
				} else if m.attacks[idx] != reference[i] {
					break
				}
			}
		}
	}

}
func attacks_bb_empty_board(s Square, pt PieceType) BitBoard {
	return pseudoAttacks[pt][s]
}
func attacks_bb(s Square, occupied BitBoard, pt PieceType) BitBoard {
	switch pt {
	case BISHOP:
		return BishopMagics[s].attacks[BishopMagics[s].index(occupied)]

	case ROOK:
		return RookMagics[s].attacks[RookMagics[s].index(occupied)]

	case QUEEN:
		return attacks_bb(s, occupied, BISHOP) | attacks_bb(s, occupied, ROOK)
	default:
		return BitBoard(0)
	}
}
func safe_destination(s Square, d int) BitBoard {
	var to Square = Square(int(s) + d)
	if to.is_square_ok() && distance(s, to) <= 2 {
		return square_bb(to)
	} else {
		return 0
	}
}

func sliding_attacks(pt PieceType, sq Square, occupied BitBoard) BitBoard {
	var attacks BitBoard = 0
	var RookDirections = [4]Direction{NORTH, SOUTH, EAST, WEST}
	var BishopDirection = [4]Direction{NORTH_EAST, SOUTH_EAST, SOUTH_WEST, NORTH_WEST}
	var dir *[4]Direction
	if pt == ROOK {
		dir = &RookDirections
	} else {
		dir = &BishopDirection
	}

	for _, d := range dir {
		s := sq
		for safe_destination(s, int(d)) != 0 {
			s += Square(d)
			attacks = bitboard_or_square(attacks, s)
			if bitboard_and_square(occupied, s) > 0 {
				break
			}
		}
	}
	return attacks
}

func init_bitboards() {
	var i uint32 = 0
	for ; i < (1 << 16); i++ {
		popCnt16[i] = uint8(bits.OnesCount16(uint16(i))) //Might Blow up later "watch out"
	}
	for s1 := SQ_A1; s1 <= SQ_H8; s1++ {
		for s2 := SQ_A1; s2 <= SQ_H8; s2++ {
			ass := file_distance(s1, s2)
			if ass < rank_distance(s1, s2) {
				ass = rank_distance(s1, s2)
			}
			squareDistance[s1][s2] = ass
		}
	}
	init_magics(ROOK, RookTable[:], RookMagics[:])
	init_magics(BISHOP, BishopTable[:], BishopMagics[:])
	for s1 := SQ_A1; s1 <= SQ_H8; s1++ {
		PawnAttacks[WHITE][s1] = square_bb(s1).calc_pawn_attacks_bb(WHITE)
		PawnAttacks[BLACK][s1] = square_bb(s1).calc_pawn_attacks_bb(BLACK)
		for _, step := range []int{-9, -8, -7, -1, 1, 7, 8, 9} {
			pseudoAttacks[KING][s1] |= safe_destination(s1, step)
		}
		for _, step := range []int{-17, -15, -10, -6, 6, 10, 15, 17} {
			pseudoAttacks[KNIGHT][s1] |= safe_destination(s1, step)
		}
		pseudoAttacks[BISHOP][s1] = attacks_bb(s1, 0, BISHOP)
		pseudoAttacks[ROOK][s1] = attacks_bb(s1, 0, ROOK)
		pseudoAttacks[QUEEN][s1] = pseudoAttacks[BISHOP][s1] | pseudoAttacks[ROOK][s1]
		for _, pt := range []PieceType{BISHOP, ROOK} {
			for s2 := SQ_A1; s2 <= SQ_H8; s2++ {
				if bitboard_and_square(pseudoAttacks[pt][s1], s2) != 0 {
					lineBB[s1][s2] = bitboard_or_square(bitboard_or_square((attacks_bb(s1, 0, pt)&attacks_bb(s2, 0, pt)), s1), s2)
					betweenBB[s1][s2] = (attacks_bb(s1, square_bb(s2), pt) & attacks_bb(s2, square_bb(s1), pt))
				}
				bitboard_oreq_square(&betweenBB[s1][s2], s2)
			}
		}
	}
}
