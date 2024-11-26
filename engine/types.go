package engine

type BitBoard uint64
type Key uint64
type Value int
type Color int8

var VALUE_ZERO Value = 0
var PawnValue Value = 200
var KnightValue Value = 781
var BishopValue Value = 825
var RookValue Value = 1276
var QueenValue Value = 2538

var PieceValue = []Value{VALUE_ZERO, PawnValue, KnightValue, BishopValue, RookValue, QueenValue, VALUE_ZERO, VALUE_ZERO, VALUE_ZERO, PawnValue, KnightValue, BishopValue, QueenValue, VALUE_ZERO, VALUE_ZERO}

type Player struct {
	Player Color
}

const (
	WHITE Color = iota
	BLACK
	COLOR_NB
)

type Castling_Rights int32
type CastlingRights struct {
	Castling_rights Castling_Rights
}

const (
	NO_CASTLING Castling_Rights = iota
	WHITE_SHORT Castling_Rights = iota
	WHITE_LONG  Castling_Rights = WHITE_SHORT << 1
	BLACK_SHORT Castling_Rights = WHITE_SHORT << 2
	BLACK_LONG  Castling_Rights = WHITE_SHORT << 3

	KING_SIDE         Castling_Rights = WHITE_SHORT | BLACK_SHORT
	QUEEN_SIDE        Castling_Rights = WHITE_LONG | BLACK_LONG
	WHITE_CASTLING    Castling_Rights = WHITE_SHORT | WHITE_LONG
	BLACK_CASTLING    Castling_Rights = BLACK_SHORT | BLACK_LONG
	ANY_CASTLING      Castling_Rights = WHITE_CASTLING | BLACK_CASTLING
	CASTLING_RIGHT_NB Castling_Rights = 16
)

type Bound uint8

const (
	BOUND_NONE Bound = iota
	BOUND_UPPER
	BOUND_LOWER
	BOUND_EXACT
)

type MoveType uint16

const (
	NORMAL     MoveType = iota
	PROMOTION  MoveType = 1 << 14
	EN_PASSANT MoveType = 2 << 14
	CASTLING   MoveType = 3 << 14
)

// In a square, the first 3 bits indicate the file and the 3 bits after indicate the rank
type Square int32

const (
	SQ_A1 Square = iota
	SQ_B1
	SQ_C1
	SQ_D1
	SQ_E1
	SQ_F1
	SQ_G1
	SQ_H1

	SQ_A2
	SQ_B2
	SQ_C2
	SQ_D2
	SQ_E2
	SQ_F2
	SQ_G2
	SQ_H2

	SQ_A3
	SQ_B3
	SQ_C3
	SQ_D3
	SQ_E3
	SQ_F3
	SQ_G3
	SQ_H3

	SQ_A4
	SQ_B4
	SQ_C4
	SQ_D4
	SQ_E4
	SQ_F4
	SQ_G4
	SQ_H4

	SQ_A5
	SQ_B5
	SQ_C5
	SQ_D5
	SQ_E5
	SQ_F5
	SQ_G5
	SQ_H5

	SQ_A6
	SQ_B6
	SQ_C6
	SQ_D6
	SQ_E6
	SQ_F6
	SQ_G6
	SQ_H6

	SQ_A7
	SQ_B7
	SQ_C7
	SQ_D7
	SQ_E7
	SQ_F7
	SQ_G7
	SQ_H7

	SQ_A8
	SQ_B8
	SQ_C8
	SQ_D8
	SQ_E8
	SQ_F8
	SQ_G8
	SQ_H8

	SQ_ZERO   = 0
	SQUARE_NB = 64
	SQ_NONE   = 64
)

// A move needs 16 bits

// bits 0 - 5 destination square (from 0 - 63)
// bits 6 - 11 origin square (from 0 - 63)
// bits 12-13 promotion piece type (from KNIGHT-2 to QUEEN-2)
// bits 14-15 special move flag: promotion (1), en passant (2) castling (3)
// NOTE: en passant bit is set only when a pawn can be captured
// Special cases are Move::none() and Move::null(). Which have the same origin and destination squares
type Move struct {
	data uint16
}

type Direction int16

const (
	EAST       Direction = iota + 1
	NORTH      Direction = 8
	SOUTH      Direction = -NORTH
	WEST       Direction = -EAST
	NORTH_EAST Direction = NORTH + EAST
	SOUTH_EAST Direction = SOUTH + EAST
	SOUTH_WEST Direction = SOUTH + WEST
	NORTH_WEST Direction = NORTH + WEST
)

const FILE_NB = 8
const CASTLING_RIGHTS_NB = 16

type File int32
type Rank int32

const (
	FILE_A File = iota
	FILE_B
	FILE_C
	FILE_D
	FILE_E
	FILE_F
	FILE_G
	FILE_H
)

const (
	RANK_1 Rank = iota
	RANK_2
	RANK_3
	RANK_4
	RANK_5
	RANK_6
	RANK_7
	RANK_8
)

type PieceType uint16

const (
	NO_PIECE_TYPE PieceType = iota
	PAWN
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
	ALL_PIECES
	PIECE_TYPE_NB
)

type Piece uint16

const (
	NO_PIECE Piece = iota
	W_PAWN
	W_KNIGHT
	W_BISHOP
	W_ROOK
	W_QUEEN
	W_KING

	B_PAWN = iota + 2 // B_PAWN - W_PAWN + 8 (iota = 7 here)
	B_KNIGHT
	B_BISHOP
	B_ROOK
	B_QUEEN
	B_KING
	PIECE_NB = 16
)

func make_square(f File, r Rank) Square {
	return Square(Square(r<<3) + Square(f))
}

func (mv *Move) SetMove(from, to Square) {
	mv.data = uint16((from<<6 + to))
}
func (sq Square) rank_of() Rank {
	return Rank(sq >> 3)
}

func (sq Square) file_of() File {
	return File(sq & 7)
}
func (mv Move) from_square() Square {
	return Square(mv.data >> 6 & 0x3F)
}

func (mv Move) to_square() Square {
	return Square(mv.data & 0x3F)
}

func (mv Move) from_to_squares() (Square, Square) {
	return mv.from_square(), mv.to_square()
}

func (mv Move) type_of() MoveType {
	return MoveType(mv.data & (3 << 14))
}
func (mv Move) promotion_type() PieceType {
	return PieceType(PieceType(((mv.data >> 12) & 3)) + KNIGHT)
}

func relative_square(c Color, s Square) Square {
	return Square(s ^ Square((c * 56)))
}
func make_move_from_to(from Square, to Square) {
}

func nullMove() Move {
	return Move{65}
}

func noneMove() Move {
	return Move{0}
}

func (mv Move) is_ok() bool {
	return nullMove().data != mv.data && noneMove().data != mv.data
}

func (p Piece) piece_type() PieceType {
	return PieceType(p & 7)
}

func (p Piece) color() Color {
	return Color(p >> 3)
}

func make_piece(p PieceType, c Color) Piece {
	return Piece(Piece((c << 3)) + Piece(p))
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
func bitboard_oreq_square(b *BitBoard, s Square) {
	*b |= square_bb(s)
}

func bitboard_andeq_square(b *BitBoard, s Square) {
	*b &= square_bb(s)
}

func (c Color) pawn_push() Direction {
	if c == WHITE {
		return NORTH
	} else {
		return SOUTH
	}
}

// BitBoards for a file or rank
func rank_bb(r Rank) BitBoard {
	return Rank1BB << (8 * r)
}

func file_bb(f File) BitBoard {
	return FileABB << f
}

func file_distance(s1 Square, s2 Square) uint8 {
	ret := s1.file_of() - s2.file_of()
	if ret < 0 {
		return uint8(-ret)
	} else {
		return uint8(ret)
	}
}

func rank_distance(s1 Square, s2 Square) uint8 {
	ret := s1.rank_of() - s2.rank_of()
	if ret < 0 {
		return uint8(-ret)
	} else {
		return uint8(ret)
	}
}

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
func (s Square) is_square_ok() bool {
	return s >= SQ_A1 && s <= SQ_H8
}
func (bb BitBoard) shift_double_north() BitBoard {
	return bb << 16
}

func (bb BitBoard) shift_double_south() BitBoard {
	return bb >> 16
}
func piece_to_char(p Piece) string {
	switch p {
	case W_PAWN:
		return "♙"
	case W_KNIGHT:
		return "♘"
	case W_BISHOP:
		return "♗"
	case W_ROOK:
		return "♖"
	case W_QUEEN:
		return "♕"
	case W_KING:
		return "♔"
	case B_PAWN:
		return "♟"
	case B_KNIGHT:
		return "♞"
	case B_BISHOP:
		return "♝"
	case B_ROOK:
		return "♜"
	case B_QUEEN:
		return "♛"
	case B_KING:
		return "♚"
	default:
		return ""
	}
}
