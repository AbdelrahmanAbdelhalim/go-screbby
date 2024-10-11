package board

//Stateinfo struct stores information to restore a position object to its previous state when we
//retract a move
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

func (board Board) do_move(m Move, newSt *StateInfo) {}

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

func (board Board) put_piece(p Piece, s Square) {
	board.Board[s] = p
	board.ByTypeBB[ALL_PIECES] |= board.ByTypeBB[p.piece_type()]
	board.ByTypeBB[ALL_PIECES] |= BitBoard(s)
	board.PieceCount[p]++
	board.PieceCount[make_piece(Piece(ALL_PIECES), p.color())]++
}
func (board Board) remove_piece(p Piece, s Square) {
	piece := board.Board[s]
	board.ByTypeBB[ALL_PIECES] ^= BitBoard(s)
	board.ByTypeBB[p.piece_type()] ^= BitBoard(s)
	board.ByColorBB[p.color()] ^= BitBoard(s)
	board.Board[s] = NO_PIECE
	board.PieceCount[piece]--
	board.PieceCount[make_piece(Piece(ALL_PIECES), p.color())]++
}
func (board Board) move_piece(from Square, to Square) {
	piece := board.Board[from]
	fromto := BitBoard(from | to)
	board.ByTypeBB[piece.piece_type()] ^= fromto
	board.ByColorBB[piece.color()] ^= fromto
	board.ByTypeBB[ALL_PIECES] ^= fromto
	board.Board[from] = NO_PIECE
	board.Board[to] = piece
}

func (board Board) init() {

}
