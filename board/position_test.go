package board

import "testing"

func TestPositionString(t *testing.T) {
	tests := []struct {
		pos      *Position
		expected string
	}{
		{
			StartingPosition(),
			`
8 r n b q k b n r 
7 p p p p p p p p 
6 . . . . . . . . 
5 . . . . . . . . 
4 . . . . . . . . 
3 . . . . . . . . 
2 P P P P P P P P 
1 R N B Q K B N R 
  a b c d e f g h 
`,
		},
	}

	for _, test := range tests {
		if str := test.pos.String(); str != test.expected {
			t.Errorf("Expected position to have string representation:\n%s\nGot:\n%s", test.expected, str)
		}
	}
}
