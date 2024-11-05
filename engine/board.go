package engine

import (
	"strconv"
	"strings"
)

// Stateinfo struct stores information to restore a position object to its previous state when we
// retract a move

type StateInfo struct {

	//Copied when making a move
	material_key    Key
	pawn_Key        Key
	nonPawnMaterial [COLOR_NB]Value
	CastlingRights  int32
	rule50          int32
	pliesFromNull   int32
	EpSquare        Square

	//Not Copied When making a move (Will be recomputed)
	key             Key
	checkersBB      BitBoard
	Previous        *StateInfo
	blockersForKing [COLOR_NB]BitBoard
	pinners         [COLOR_NB]BitBoard
	checkSquares    [PIECE_TYPE_NB]BitBoard
	capturedPiece   Piece
	repition        int32
}

type Board struct {
	//Data Members
	Board              [SQUARE_NB]Piece
	ByTypeBB           [PIECE_TYPE_NB]BitBoard
	ByColorBB          [COLOR_NB]BitBoard
	PieceCount         [PIECE_NB]int32
	CastlingRightsMask [SQUARE_NB]int32
	CastlingRookSquare [CASTLING_RIGHT_NB]Square
	CastlingPath       [CASTLING_RIGHT_NB]BitBoard
	St                 *StateInfo
	gamePly            int32
	SideToMove         Color
}

var zobrist_st *Zobrist
var cuckkoo_st *Cuckoo

func (board Board) get_sides_piecesbb(c Color) BitBoard {
	return board.ByColorBB[c]
}
func (board Board) get_piecebb_by_type(pt PieceType) BitBoard {
	return board.ByTypeBB[pt]
}

func (board Board) get_pieces_bb(pieces ...PieceType) BitBoard {
	bb := BitBoard(0)
	for _, pt := range pieces {
		bb |= board.get_piecebb_by_type(pt)
	}
	return bb
}

func (board Board) get_typed_sides_pieces_bb(c Color, pieces ...PieceType) BitBoard {
	return board.get_sides_piecesbb(c) & board.get_pieces_bb(pieces...)
}

func (board Board) get_piece_bb(pt PieceType, c Color) BitBoard {
	return board.ByTypeBB[pt] & board.ByColorBB[c]
}

func (board Board) print_position() {
	var builder strings.Builder
	builder.WriteString("\n +---+---+---+---+---+---+---+---+\n")

	for rank := RANK_8; rank >= RANK_1; rank-- {
		for file := FILE_A; file <= FILE_H; file++ {
			builder.WriteString(" | ")
			builder.WriteString(piece_to_char(board.piece_on(make_square(file, rank))))
		}
		builder.WriteString(" | ")
		builder.WriteString(strconv.Itoa(int(1 + rank)))
		builder.WriteString(" | ")
	}
	builder.WriteString("   a   b   c   d   e   f   g   h\n")

	//Probe TableBase and output results
	//Print FEN
}

//functions for cuckoo tables for repitition detection

func H1(h Key) int {
	return int(h & 0x1fff)
}

func H2(h Key) int {
	return int((h >> 16) & 0x1fff)
}

func init_zobrist_cuckoo() {
	zobrist_st, cuckoo_st := NewZobristCuckoo()
}

func (board *Board) do_move(m Move, newSt *StateInfo) {
	k := board.St.key ^ zobrist_st.side
	newSt.copy_from_old_st(*board.St)
	newSt.Previous = board.St
	board.St = newSt

	board.gamePly++
	board.St.rule50++
	board.St.pliesFromNull++

	var us Color = board.SideToMove
	var them Color = ^us
	var from Square = m.from_square()
	var to Square = m.to_square()
	var pc Piece = board.piece_on(from)
	var captured Piece
	if m.type_of() == EN_PASSANT {
		captured = make_piece(PAWN, them)
	} else {
		board.piece_on(to)
	}
}

func (board *Board) set_position_from_fen(fen string, st *StateInfo) {

}

func (st *StateInfo) copy_from_old_st(oldSt StateInfo) {

}
func (board Board) piece_on(sq Square) Piece {
	return board.Board[sq]
}

func (board Board) side_to_move() Color {
	return board.SideToMove
}

func (board Board) ep_square() Square {
	return board.St.EpSquare
}

func (board Board) empty(s Square) bool {
	return board.piece_on(s) == NO_PIECE
}

func (board Board) castling_rights(c Color) Castling_Rights {
	return Castling_Rights(c) & Castling_Rights(board.St.CastlingRights)
}

func (board Board) can_castle(cr Castling_Rights) bool {
	return Castling_Rights(board.St.CastlingRights)&cr > 0
}

func (board Board) castling_rook_square(cr Color) Square {
	return board.CastlingRookSquare[cr]
}

func (board Board) checkers() BitBoard {
	return board.St.checkersBB
}

func (board Board) blockers_for_king(c Color) BitBoard {
	return board.St.blockersForKing[c]
}

func (board Board) pinners(c Color) BitBoard {
	return board.St.pinners[c]
}

func (board Board) check_squares(pt PieceType) BitBoard {
	return board.St.checkSquares[pt]
}

func (board Board) attacks_by(c Color) BitBoard {
	return FileABB
}

func (board Board) game_ply() int32 {
	return board.gamePly
}

func (board Board) rule_50_count() int32 {
	return board.St.rule50
}

func (board Board) pawn_key() Key {
	return board.St.pawn_Key
}

func (board Board) material_key() Key {
	return board.St.material_key
}

func (board Board) captured_piece() Piece {
	return board.St.capturedPiece
}

func (board *Board) put_piece(p Piece, s Square) {
	board.Board[s] = p
	board.ByTypeBB[ALL_PIECES] |= board.ByTypeBB[p.piece_type()]
	board.ByTypeBB[ALL_PIECES] |= square_bb(s)
	board.PieceCount[p]++
	board.PieceCount[make_piece(Piece(ALL_PIECES), p.color())]++
}
func (board *Board) remove_piece(p Piece, s Square) {
	piece := board.Board[s]
	board.ByTypeBB[ALL_PIECES] ^= square_bb(s)
	board.ByTypeBB[p.piece_type()] ^= square_bb(s)
	board.ByColorBB[p.color()] ^= square_bb(s)
	board.Board[s] = NO_PIECE
	board.PieceCount[piece]--
	board.PieceCount[make_piece(Piece(ALL_PIECES), p.color())]++
}
func (board *Board) move_piece(from Square, to Square) {
	piece := board.Board[from]
	fromto := BitBoard(from | to)
	board.ByTypeBB[piece.piece_type()] ^= fromto
	board.ByColorBB[piece.color()] ^= fromto
	board.ByTypeBB[ALL_PIECES] ^= fromto
	board.Board[from] = NO_PIECE
	board.Board[to] = piece
}
