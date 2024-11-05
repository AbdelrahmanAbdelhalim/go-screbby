package engine

import (
	"fmt"
	"testing"
)

func TestSquarebb(t *testing.T) {
	sq := SQ_E4
	sq_bb := square_bb(sq)
	exp_res := BitBoard(1 << 28)
	if sq_bb != exp_res {
		t.Errorf("Square bb conversion failed, expected %v, found %v", exp_res, sq_bb)
	}
}

func TestMakeSquare(t *testing.T) {
	file := FILE_E
	rank := RANK_4
	exp_sq := SQ_E4
	sq := make_square(file, rank)
	if exp_sq != sq {
		t.Errorf("make_square failed, expected %v found %v", exp_sq, sq)
	}
}

func TestBbXorSquare(t *testing.T) {
	sq := SQ_E4
	bb := BitBoard(1<<29) | BitBoard(1<<28)
	exp_bb := BitBoard(1 << 29)
	bb = bitboard_xor_square(bb, sq)
	if bb != exp_bb {
		t.Errorf("BB and Square failed, expected %v, found %v", exp_bb, bb)
	}
}

func TestBbOrSquare(t *testing.T) {
	sq := SQ_E4
	bb := BitBoard(1 << 29)
	exp_bb := BitBoard(1<<28) | BitBoard(1<<29)
	bb = bitboard_or_square(bb, sq)
	if bb != exp_bb {
		t.Errorf("BB and Square failed, expected %v, found %v", exp_bb, bb)
	}
}

func TestBbAndSquare(t *testing.T) {
	sq := SQ_E4
	bb := BitBoard(1<<28) | BitBoard(1<<29)
	exp_bb := BitBoard(1 << 28)
	bb = bitboard_and_square(bb, sq)
	if bb != exp_bb {
		t.Errorf("BB and Square failed, expected %v, found %v", exp_bb, bb)
	}
}

func TestPrintBB(t *testing.T) {
	sq_1 := SQ_E4
	sq_2 := SQ_E5
	bb := square_bb(sq_1) | square_bb(sq_2)
	s := printbb(bb)
	fmt.Println(s)
	fmt.Println("----------------------------------")
}

func TestShiftBBEast(t *testing.T) {
	sq := []Square{SQ_A1, SQ_H1, SQ_E4, SQ_E5}
	dest_sq := []Square{SQ_B1, SQ_F4, SQ_F5}
	var bb BitBoard = 0
	for _, square := range sq {
		bb = bitboard_or_square(bb, square)
	}
	var exp_bb BitBoard = 0
	for _, square := range dest_sq {
		exp_bb = bitboard_or_square(exp_bb, square)
	}
	fmt.Println("Test Shift East Before")
	fmt.Println(printbb(bb))
	new_bb := bb.shift(EAST)
	fmt.Println("Test Shift East After")
	fmt.Println(printbb(new_bb))

	fmt.Println("----------------------------------")
	if exp_bb != new_bb {
		t.Error("Shift East Failed, Review Output")
	}
}

func TestShiftBBWest(t *testing.T) {
	sq := []Square{SQ_A1, SQ_H1, SQ_E4, SQ_E5}
	dest_sq := []Square{SQ_G1, SQ_D4, SQ_D5}
	var bb BitBoard = 0
	for _, square := range sq {
		bb = bitboard_or_square(bb, square)
	}
	var exp_bb BitBoard = 0
	for _, square := range dest_sq {
		exp_bb = bitboard_or_square(exp_bb, square)
	}

	fmt.Println("Test Shift West Before")
	fmt.Println(printbb(bb))
	new_bb := bb.shift(WEST)

	fmt.Println("Test Shift West After")
	fmt.Println(printbb(new_bb))

	fmt.Println("----------------------------------")
	if exp_bb != new_bb {
		t.Error("Shift West Failed, Review Output")
	}
}

func TestShiftBBNorthEast(t *testing.T) {
	sq := []Square{SQ_A1, SQ_H1, SQ_E4, SQ_E5, SQ_G8}
	dest_sq := []Square{SQ_B2, SQ_F5, SQ_F6}
	var bb BitBoard = 0
	for _, square := range sq {
		bb = bitboard_or_square(bb, square)
	}
	var exp_bb BitBoard = 0
	for _, square := range dest_sq {
		exp_bb = bitboard_or_square(exp_bb, square)
	}

	fmt.Println("Test Shift North East Before")
	fmt.Println(printbb(bb))
	new_bb := bb.shift(NORTH_EAST)

	fmt.Println("Test Shift North East After")
	fmt.Println(printbb(new_bb))

	fmt.Println("----------------------------------")
	if exp_bb != new_bb {
		t.Error("Shift North East Failed, Review Output")
	}
}

func TestShiftBBNorthWest(t *testing.T) {
	sq := []Square{SQ_A1, SQ_H1, SQ_E4, SQ_E5, SQ_H8}
	dest_sq := []Square{SQ_G2, SQ_D5, SQ_D6}
	var bb BitBoard = 0
	for _, square := range sq {
		bb = bitboard_or_square(bb, square)
	}
	var exp_bb BitBoard = 0
	for _, square := range dest_sq {
		exp_bb = bitboard_or_square(exp_bb, square)
	}

	fmt.Println("Test Shift North West Before")
	fmt.Println(printbb(bb))
	new_bb := bb.shift(NORTH_WEST)

	fmt.Println("Test Shift North West After")
	fmt.Println(printbb(new_bb))

	fmt.Println("----------------------------------")
	if exp_bb != new_bb {
		t.Error("Shift North West Failed, Review Output")
	}
}

func TestShiftBBSouthEast(t *testing.T) {
	sq := []Square{SQ_A1, SQ_H1, SQ_E4, SQ_E5, SQ_H8}
	dest_sq := []Square{SQ_F3, SQ_F4}
	var bb BitBoard = 0
	for _, square := range sq {
		bb = bitboard_or_square(bb, square)
	}
	var exp_bb BitBoard = 0
	for _, square := range dest_sq {
		exp_bb = bitboard_or_square(exp_bb, square)
	}

	fmt.Println("Test Shift South East Before")
	fmt.Println(printbb(bb))
	new_bb := bb.shift(SOUTH_EAST)

	fmt.Println("Test Shift South East After")
	fmt.Println(printbb(new_bb))

	fmt.Println("----------------------------------")
	if exp_bb != new_bb {
		t.Error("Shift North West Failed, Review Output")
	}
}

func TestShiftBBSouthWest(t *testing.T) {
	sq := []Square{SQ_A1, SQ_H1, SQ_E4, SQ_E5, SQ_H8, SQ_A8}
	dest_sq := []Square{SQ_D3, SQ_D4, SQ_G7}
	var bb BitBoard = 0
	for _, square := range sq {
		bb = bitboard_or_square(bb, square)
	}
	var exp_bb BitBoard = 0
	for _, square := range dest_sq {
		exp_bb = bitboard_or_square(exp_bb, square)
	}

	fmt.Println("Test Shift South West Before")
	fmt.Println(printbb(bb))
	new_bb := bb.shift(SOUTH_WEST)

	fmt.Println("Test Shift South West After")
	fmt.Println(printbb(new_bb))

	fmt.Println("----------------------------------")
	if exp_bb != new_bb {
		t.Error("Shift North West Failed, Review Output")
	}
}

func Testmagics(t *testing.T) {

}
