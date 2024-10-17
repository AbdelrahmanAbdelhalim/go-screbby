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
var squareDistance [SQUARE_NB][SQUARE_NB]uint8

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
	}
}
func attacks_bb(s Square, occupied BitBoard, pt PieceType, magics Magics) BitBoard {
	switch pt {
	case BISHOP:
		return magics.BishopMagics[s].attacks[magics.BishopMagics[s].index(occupied)]

	case ROOK:
		return magics.BishopMagics[s].attacks[magics.BishopMagics[s].index(occupied)]

	case QUEEN:
		return magics.BishopMagics[s].attacks[magics.BishopMagics[s].index(occupied)]
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
	for s1 := SQ_A1; s1 <= SQ_H8; s1++ {
		for s2 := SQ_A1; s2 <= SQ_H8; s2++ {
			ass := file_distance(s1, s2)
			if ass < rank_distance(s1, s2) {
				ass = rank_distance(s1, s2)
			}
			squareDistance[s1][s2] = ass
		}
	}
}
