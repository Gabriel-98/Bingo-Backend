package types

// ValidateCell validates if a cell has already been marked by the user and if the value
// of that cell has already been drawn. The validation for the cell (r=3,c=3) will always
// return true and for cells outside of the card (r < 1 or r > 5 or c < 1 or c > 5) will
// always return false.
func ValidateCell(drawnNumbers *DrawnNumbers, card *Card, r int, c int) bool {
	if r == 3 && c == 3 {
		return true
	}
	return card.IsMarked(r,c) && drawnNumbers.Contains(card.Value(r,c))
}

// CardValidator is the interface that wraps the method for validating bingo cards and the
// user's marks against the numbers already drawn.
type CardValidator interface {
	Validate(drawnNumbers *DrawnNumbers, card *Card) bool
}

// FullCardValidator provides the validation that all cells on the card are marked and
// their numbers have already been drawn.
type FullCardValidator struct {}

func NewFullCardValidator() *FullCardValidator {
	return &FullCardValidator{}
}

func (v *FullCardValidator) Validate(drawnNumbers *DrawnNumbers, card *Card) bool {
	for r:=1; r<=5; r++ {
		for c:=1; c<=5; c++ {
			if !ValidateCell(drawnNumbers, card, r, c) {
				return false
			}
		}
	}
	return true
}

// RowCardValidator provides the validation that all cells in any row are marked and their
// numbers have already been drawn.
type RowCardValidator struct {}

func NewRowCardValidator() *RowCardValidator {
	return &RowCardValidator{}
}

func (v *RowCardValidator) Validate(drawnNumbers *DrawnNumbers, card *Card) bool {
	for r:=1; r<=5; r++ {
		rowFilled := true
		for c:=1; c<=5 && rowFilled; c++ {
			if !ValidateCell(drawnNumbers, card, r, c) {
				rowFilled = false
			}
		}
		if rowFilled {
			return true
		}
	}
	return false
}

// ColCardValidator provides the validation that all cells in any column are marked and
// their numbers have already been drawn.
type ColCardValidator struct {}

func NewColCardValidator() *ColCardValidator {
	return &ColCardValidator{}
}

func (v *ColCardValidator) Validate(drawnNumbers *DrawnNumbers, card *Card) bool {
	for c:=1; c<=5; c++ {
		colFilled := true
		for r:=1; r<=5 && colFilled; r++ {
			if !ValidateCell(drawnNumbers, card, r, c) {
				colFilled = false
			}
		}
		if colFilled {
			return true
		}
	}
	return false
}

// DiagonalCardValidator provides the validation that all cells in either of the two
// diagonals are marked and their numbers have already been drawn.
type DiagonalCardValidator struct {}

func NewDiagonalCardValidator() *DiagonalCardValidator {
	return &DiagonalCardValidator{}
}

func (v *DiagonalCardValidator) Validate(drawnNumbers *DrawnNumbers, card *Card) bool {
	dc1, dc2 := 0, 0
	for r,c1,c2:=1,1,5; r<=5; r,c1,c2=r+1,c1+1,c2-1 {
		if ValidateCell(drawnNumbers, card, r, c1) {
			dc1++
		}
		if ValidateCell(drawnNumbers, card, r, c2) {
			dc2++
		}
	}
	return dc1 == 5 || dc2 == 5
}

// CornerCardValidator provides the validation that all corner cells are marked and their
// numbers have already been drawn.
type CornerCardValidator struct {}

func NewCornerCardValidator() *CornerCardValidator {
	return &CornerCardValidator{}
}

func (v *CornerCardValidator) Validate(drawnNumbers *DrawnNumbers, card *Card) bool {
	if ValidateCell(drawnNumbers, card, 1, 1) {
		if ValidateCell(drawnNumbers, card, 1, 5) {
			if ValidateCell(drawnNumbers, card, 5, 1) {
				return ValidateCell(drawnNumbers, card, 5, 5)
			}
		}
	}
	return false
}