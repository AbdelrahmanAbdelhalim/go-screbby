package board

//Stateinfo struct stores information to restore a position object to its previous state when we
//retract a move
type StateInfo struct {

	//Copied when making a move
	Material_key    Key
	Pawn_Key        Key
	NonPawnMaterial [COLOR_NB]Value
	CastlingRights  int32
	Rule50          int32
	PliesFromNull   int32
	EqSquare        Square

	//Not Copied When making a move (Will be recomputed)
	Key             Key
	ChckersBB       BitBoard
	Previous        *StateInfo
	BlockersForKing [COLOR_NB]BitBoard
	Pinners         [COLOR_NB]BitBoard
	CheckSquares    [PIECE_TYPE_NB]BitBoard
	CapturedPiece   Piece
	Repition        int32
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
	GamePly            int32
	SideToMove         Color
}

func (board Board) pieces(pt PieceType) BitBoard {
	return FileABB
}

func (board Board) piece_on(sq Square) Piece {
	return Piece(KING)
}

func (board Board) side_to_move() Color {
	return board.SideToMove
}

// func (board Board) set(fen String) ret {
//many implementations, figure that out
// }

// func (board Board) fen() string {

// }

// func (board Board) pieces() BitBoard {
//many implementations, figure that out
// }

// func (board Board) piece_on(s Square) Piece {
// }

// func (board Board) ep_square() Square {
// }

// func (board Board) empty(s Square) Square {
// }

// func (board Board) castling_rights() Casting_Rights {
// }
// func (board Board) can_castle(cr Castling_Right) bool {
// }

// func (board Board) castling_impeded() bool {
// }

// func (board Board) castling_rook_square() Square {
// }

// func (board Board) checkers() BitBoard {
// }

// func (board Board) blockers_for_king(c Color) BitBoard {
// }

// func (board Board) pinners() BitBoard {
// }

// func (board Board) attackrs_to(sq Square) BitBoard {
//multiple impl
// }

func (board Board) attacks_by(c Color) BitBoard {
	return FileABB
}

// func (board Board) update_slider_blocker()  {
// }

// func (board Board) legal(m Move) bool {
// }
// func (board Board) pseudo_legal(m Move) bool {
// }
// func (board Board) capture(m Move) bool {
// }
// func (board Board) capture_stage(m Move) bool {
// }
// func (board Board) gives_check(m Move) bool {
// }
// func (board Board) moved_piece() Piece {
// }
// func (board Board) captured_piece() Piece {
// }

func (board Board) do_move(m Move, newSt *StateInfo) {
}

// func (board Board) do_move(m Move, newSt *StateInfo, gives_check bool) {
// }
// func (board Board) undo_move(m Move) ret {
// }
// func (board Board) do_null_move(newSt *StateInfo, tt *TranspositionTable) {
// }
// func (board Board) undo_null_move(){
// }
// func (board Board) game_ply() int32 {
// }
// func (board Board) side_to_move() Color {
// }
// func (board Board) is_draw() bool {
// }
// func (board Board) has_game_cycle() bool {
// }
// func (board Board) rule_50_count() int32 {
// }
// func (board Board) has_repeated() bool {
// }

// func (board Board) non_pawn_material(c Color) Value {
// multiple impl
// }

// func (board Board) put_piece(pc Piece, sq Square) {
// }
// func (board Board) remove_piece(sq Square) {
// }

// func (board Board) set_castling_right(c Color, rfrom Square) {
// }
// func (board Board) set_state() {
// }
// func (board Board) set_check_info() {
// }

// func (board Board) move_piece(from Square, to Square) {
// }

// func (board Board) do_castling() ret {
//templates! watchout
// }

// func (board Board) adjust_key(k Key) Key {
// }
// func (board Board) funcname() ret {
// }
// func (board Board) funcname() ret {
// }
