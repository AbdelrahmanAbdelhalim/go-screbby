package board

import (
	"testing"
)

func TestMove(t *testing.T) {
	tests := []struct {
		name             string
		wantFrom, wantTo Square
		expected_data    uint16
	}{
		{"Move from 0 to 3", 0, 3, 3},   //From square 0 to some square
		{"Move from 0 to 3", 3, 0, 192}, //from some square to square 0
		{"Move from 6 to 3", 6, 3, 387}, // from some square to some square
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			move := Move{}
			move.SetMove(tt.wantFrom, tt.wantTo)
			if move.data != tt.expected_data {
				t.Errorf("Expected move data %v, but got %v", move.data, tt.expected_data)
			}
		})
	}
}

func TestNullMove(t *testing.T) {
	expected_null_move := Move{65}
	null_move := nullMove()
	if expected_null_move.data != null_move.data {
		t.Errorf("Expected null move %v, but got %v", expected_null_move.data, null_move.data)
	}
}

func TestNoneMove(t *testing.T) {
	expected_null_move := Move{0}
	null_move := noneMove()
	if expected_null_move.data != null_move.data {
		t.Errorf("Expected null move %v, but got %v", expected_null_move.data, null_move.data)
	}
}

func TestFromSquare(t *testing.T) {
	data := 4 << 6
	mv := Move{uint16(data)}
	from_square := mv.from_square()
	if from_square != Square(data>>6) {
		t.Errorf("Expected from_square %v, but got %v", data, from_square)
	}
}

func TestToSquare(t *testing.T) {
	data := 4
	mv := Move{uint16(data)}
	to_square := mv.to_square()
	if to_square != Square(data) {
		t.Errorf("Expected from_square %v, but got %v", data, to_square)
	}
}

func TestIs_ok(t *testing.T) {
	test_cases := []struct {
		test_name       string
		mv              Move
		expected_result bool
	}{
		{"Test no null_move", Move{45}, true},
		{"Test null_move", Move{65}, false},
		{"Test none_move", Move{0}, false},
	}

	for _, tt := range test_cases {
		t.Run(tt.test_name, func(t *testing.T) {
			if tt.mv.is_ok() != tt.expected_result {
				t.Errorf("Expected is_ok() result %v, found %v", tt.mv.is_ok(), tt.expected_result)
			}
		})
	}
}

func Test_ToFromSquare(t *testing.T) {
	test_cases := []struct {
		test_name string
		mv        Move
		exp_from  Square
		exp_to    Square
	}{
		{"test_1", Move{330}, 5, 10},
		{"test_2", Move{4094}, 63, 62},
		{"test_3", Move{3393}, 53, 1},
		{"test_4", Move{125}, 1, 61},
	}

	for _, tt := range test_cases {
		t.Run(tt.test_name, func(t *testing.T) {
			rec_square_from, rec_square_to := tt.mv.from_to_squares()
			if rec_square_from != tt.exp_from {
				t.Errorf("Expected square_from %v, but found %v", tt.exp_from, rec_square_from)
			}
			if rec_square_to != tt.exp_to {
				t.Errorf("Expected Square_to %v, but found %v", tt.exp_to, rec_square_to)
			}
		})
	}
}

func TestRankOfSquare(t *testing.T) {
	sq := SQ_E6
	var exp_rank Rank = 5
	rk := sq.rank_of()
	if exp_rank != rk {
		t.Errorf("Rank_of for square failed expected %v, got %v", exp_rank, rk)
	}
}

func TestFileOfSquare(t *testing.T) {
	sq := SQ_A2
	var exp_file File = 0
	f := sq.file_of()
	if exp_file != f {
		t.Errorf("Rank_of for square failed expected %v, got %v", exp_file, f)
	}
}
