package engine

import (
	"math/bits"
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
var cuckoo_st *Cuckoo
var PieceToChar = " PNBRQK pnbrqk"

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
	zobrist_st, cuckoo_st = NewZobristCuckoo()
}

func (board Board) gives_check(m Move) bool {
	from := m.from_square()
	to := m.to_square()

	//Direct Check
	if bitboard_and_square(board.check_squares(board.piece_on(from).piece_type()), to) != 0 {
		return true
	}

	//Discovered Check
	if bitboard_and_square(board.blockers_for_king(^board.SideToMove), from) != 0 {
		return !aligned(from, to, board.square(KING, ^board.SideToMove)) || m.type_of() == CASTLING
	}

	switch m.type_of() {
	case NORMAL:
		return false
	case PROMOTION:
		return bitboard_and_square(attacks_bb(to, bitboard_or_square(board.pieces(), from), m.promotion_type()), board.square(KING, ^board.SideToMove)) != 0
	case EN_PASSANT:
		{
			capsq := make_square(to.file_of(), from.rank_of())
			b := bitboard_or_square((bitboard_xor_square(bitboard_xor_square(board.pieces(), from), capsq)), to)

			a := attacks_bb(board.square(KING, ^board.SideToMove), b, ROOK)
			c := board.pieces_by_color_and_piecetype(board.SideToMove, QUEEN, ROOK)
			d := a & c
			e := attacks_bb(board.square(KING, ^board.SideToMove), b, BISHOP)
			f := board.pieces_by_color_and_piecetype(board.SideToMove, QUEEN, BISHOP)
			g := e & f
			return (d | g) != 0
		}
	default:
		{
			var path Square
			if to > from {
				path = SQ_F1
			} else {
				path = SQ_D1
			}
			rto := relative_square(board.SideToMove, path)
			return bitboard_and_square(board.check_squares(ROOK), rto) != 0
		}
	}
}

func (board Board) square(pt PieceType, c Color) Square {
	return Square(bits.TrailingZeros64(uint64(board.pieces_by_color_and_piecetype(c, pt))))
}

func (board *Board) do_move(m Move, newSt *StateInfo, givesCheck bool) {
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

	if m.type_of() == CASTLING {

	}

	if captured != NO_PIECE {
		var capsq = to
		if captured.piece_type() == PAWN {
			if m.type_of() == EN_PASSANT {
				capsq -= Square(us.pawn_push())
			}
			board.St.pawn_Key ^= zobrist_st.psq[captured][capsq]
		} else {
			board.St.nonPawnMaterial[them] -= PieceValue[captured]
		}
		board.remove_piece(capsq)
		k ^= zobrist_st.psq[captured][capsq]
		board.St.material_key ^= zobrist_st.psq[captured][board.PieceCount[captured]]
		board.St.rule50 = 0
	}

	k ^= zobrist_st.psq[pc][from] ^ zobrist_st.psq[pc][to]
	if board.St.EpSquare != SQ_NONE {
		k ^= zobrist_st.enpassant[board.St.EpSquare.file_of()]
		board.St.EpSquare = SQ_NONE
	}

	//todo: castling rights update
	if m.type_of() != CASTLING {
		board.move_piece(from, to)
	}

	if pc.piece_type() == PAWN {
		if (int(to)^int(from)) == 16 && (pawn_attacks_bb(us, Square(int(to)-int(us.pawn_push())))&board.pieces_by_color_and_piecetype(them, PAWN) != 0) {
			board.St.EpSquare = Square(int(to) - int(us.pawn_push()))
			k ^= zobrist_st.enpassant[board.St.EpSquare.file_of()]
		} else if m.type_of() == PROMOTION {
			var promotion Piece = make_piece(m.promotion_type(), us)
			board.remove_piece(to)
			board.put_piece(promotion, to)
			k ^= zobrist_st.psq[pc][to] ^ zobrist_st.psq[promotion][to]
			board.St.pawn_Key ^= zobrist_st.psq[pc][to]
			board.St.material_key ^= zobrist_st.psq[promotion][board.PieceCount[promotion]-1] ^ zobrist_st.psq[pc][board.PieceCount[pc]]
			board.St.nonPawnMaterial[us] += PieceValue[promotion]
		}
		board.St.pawn_Key ^= zobrist_st.psq[pc][from] ^ zobrist_st.psq[pc][to]
		board.St.rule50 = 0
	}
	board.St.capturedPiece = captured
	board.St.key = k
	if givesCheck {
		board.St.checkersBB = attackers_to(Square(KING, them)) & pieces(us)
	} else {
		board.St.checkersBB = 0
	}
	board.SideToMove = ^board.SideToMove
	board.set_check_info()
	board.St.repition = 0
	var end int
	if board.St.rule50 < board.St.pliesFromNull {
		end = int(board.St.rule50)
	} else {
		end = int(board.St.pliesFromNull)
	}
	if end >= 4 {

	}
}

//attackers_to computes a bitboard of all pieces which attack a given square.
func (board Board) attackers_to(s Square, occupied BitBoard) BitBoard {
	return ((pawn_attacks_bb(BLACK, s) & board.pieces_by_color_and_piecetype(WHITE, PAWN)) |
		(pawn_attacks_bb(WHITE, s) & board.pieces_by_color_and_piecetype(BLACK, PAWN)) |
		(pseudo_attacks_bb(s, KNIGHT) & board.pieces(KNIGHT)) |
		(pseudo_attacks_bb(s, KING) & board.pieces(KING)) |
		(attacks_bb(s, occupied, ROOK) & board.pieces(ROOK, QUEEN)) |
		(attacks_bb(s, occupied, BISHOP) & board.pieces(BISHOP, QUEEN)))
}
func (board *Board) set_check_info() {

}

func (board *Board) undo_move(m Move) {
	board.SideToMove = ^board.SideToMove

	var us Color = board.SideToMove
	var from Square = m.from_square()
	var to Square = m.to_square()
	var pc Piece = board.piece_on(to)
	if m.type_of() == PROMOTION {
		board.remove_piece(to)
	}
	pc = make_piece(PAWN, us)
	board.put_piece(pc, to)

	if m.type_of() == CASTLING {
		//TODO:
	} else {
		board.move_piece(to, from)
	}

	if board.St.capturedPiece != NO_PIECE {
		capsq := to
		if m.type_of() == EN_PASSANT {
			capsq -= Square(us.pawn_push())
		}
		board.put_piece(board.St.capturedPiece, capsq)
	}
	board.St = board.St.Previous
	board.gamePly--
	//assedrt(board.pos_is_ok())
}

//todo: needs transposition table
// func (board *Board) do_null_move(newSt *StateInfo, tt *TranspositionTable) {
// 	// todo: std::memcpy(&newSt, st, offsetof(StateInfo, accumlatorBig))

// }

// func (board *Board) undo_null_move() bool {

// }

////////////////////////////////////////

func (board Board) key_after(m Move) Key {
	from := m.from_square()
	to := m.to_square()
	pc := board.piece_on(from)
	captured := board.piece_on(to)
	k := board.St.key ^ zobrist_st.side
	if captured != NO_PIECE {
		k ^= zobrist_st.psq[captured][to]
	}
	k ^= zobrist_st.psq[pc][to] ^ zobrist_st.psq[pc][from]
	return k
}

func (board Board) has_game_cycle(ply int) bool {
	var j int
	var end int
	if board.St.rule50 < board.St.pliesFromNull {
		end = int(board.St.rule50)
	} else {
		end = int(board.St.pliesFromNull)
	}
	if end < 3 {
		return false
	}
	var originalKey Key = board.St.key
	var stp *StateInfo = board.St.Previous

	for i := 3; i <= end; i += 2 {
		stp = stp.Previous.Previous
		moveKey := originalKey ^ stp.key
		if cuckoo_st.cuckoo[H1(moveKey)] == moveKey || cuckoo_st.cuckoo[H1(moveKey)] == moveKey {
			if cuckoo_st.cuckoo[H1(moveKey)] == moveKey {
				j = H1(moveKey)
			} else if cuckoo_st.cuckoo[H2(moveKey)] == moveKey {
				j = H2(moveKey)
			}
			move := cuckoo_st.cuckooMove[j]
			s1 := move.from_square()
			s2 := move.to_square()
			if bitboard_or_square(between_bb(s1, s2), s2)&board.pieces_by_type(ALL_PIECES) == 0 {
				if ply > i {
					return true
				}
				var s3 Square
				if board.empty(s1) {
					s3 = s2
				} else {
					s3 = s1
				}
				if board.piece_on(s3).color() != board.side_to_move() {
					continue
				}
				if stp.repition == 0 {
					return true
				}
			}
		}

	}
	return false
}

func isDigit(b rune) bool {
	return b >= '0' && b <= '9'
}
func isLower(c rune) bool {
	return c >= 'a' && c <= 'z'
}

func (board *Board) set_position_from_fen(fen string, st *StateInfo) {
	var col, row, token byte
	var idx uint
	var sq Square = SQ_A8
	partidx := 0
	parts := strings.Split(fen, " ")

	//Piece Placement
	//todo: Possibly optimize PieceToChar. Stockfish uses StringView which needs a look on how to replicate
	part := parts[partidx]
	for _, c := range part {
		if isDigit(c) {
			sq += Square((Direction(c - '0')) * EAST)
		} else if c == '/' {
			sq += Square(2 * SOUTH)
		} else {
			idx := strings.IndexRune(PieceToChar, c)
			if idx != -1 {
				board.put_piece(Piece(idx), sq)
				sq++
			}
		}
	}

	partidx++
	if parts[partidx] == "w" {
		board.SideToMove = WHITE
	} else {
		board.SideToMove = BLACK
	}

	partidx++

	part = parts[partidx]
	for _, c := range part {
		var rsq Square
		var color Color
		if isLower(c) {
			color = WHITE
		} else {
			color = BLACK
		}
		rook := make_piece(ROOK, color)
	}
}

func (board Board) pieces_by_type(pt PieceType) BitBoard {
	return board.ByTypeBB[pt]
}

//double check if it causes problem: different implementation
func (board Board) pieces(pt ...PieceType) BitBoard {
	bb := board.ByTypeBB[pt[0]]
	for _, p := range pt {
		bb ^= board.ByTypeBB[p]
	}
	return bb
}

func (board Board) pieces_by_color(c Color) BitBoard {
	return board.ByColorBB[c]
}

func (board Board) pieces_by_color_and_piecetype(c Color, pt ...PieceType) BitBoard {
	return board.pieces_by_color(c) & board.pieces(pt...)
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
func (board *Board) remove_piece(s Square) {
	p := board.Board[s]
	board.ByTypeBB[ALL_PIECES] ^= square_bb(s)
	board.ByTypeBB[p.piece_type()] ^= square_bb(s)
	board.ByColorBB[p.color()] ^= square_bb(s)
	board.Board[s] = NO_PIECE
	board.PieceCount[p]--
	board.PieceCount[make_piece((ALL_PIECES), p.color())]++
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
