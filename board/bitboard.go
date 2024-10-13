package board

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

// Return BitBoard Representation of a Square
func square_bb(s Square) BitBoard {
	var bb BitBoard = 1
	return bb << s
}

func bitboard_and_square(b BitBoard, s Square) BitBoard {
	return b & square_bb(s)
}

func bitboard_or_square(b BitBoard, s Square) BitBoard {
	return b | square_bb(s)
}
func bitboard_xor_square(b BitBoard, s Square) BitBoard {
	return b ^ square_bb(s)
}

// BitBoards for a file or rank
func rank_bb(r Rank) BitBoard {
	return Rank1BB << (8 * r)
}

func file_bb(f File) BitBoard {
	return FileABB << f
}

// Not sure if the compiler would optimize (a >> 1) << 8 to a << 7
// Possibly needs optimization, more verbose atm to make it clear
func (bb BitBoard) shift(d Direction) BitBoard {
	switch d {
	case NORTH:
		return bb << 8
	case SOUTH:
		return bb >> 8
	case EAST:
		return (bb &^ FileHBB) << 1
	case WEST:
		return (bb &^ FileABB) >> 1
	case NORTH_EAST:
		return ((bb &^ FileHBB) << 1) << 8
	case NORTH_WEST:
		return ((bb &^ FileABB) >> 1) << 8
	case SOUTH_EAST:
		return (bb &^ FileHBB << 1) >> 8
	case SOUTH_WEST:
		return (bb &^ FileABB >> 1) >> 8
	default:
		return 0
	}
}

func (bb BitBoard) shift_double_north() BitBoard {
	return bb << 16
}

func (bb BitBoard) shift_double_south() BitBoard {
	return bb >> 16
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
