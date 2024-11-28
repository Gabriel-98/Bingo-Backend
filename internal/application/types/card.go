package types

// A Card represents a card bingo, and has the following restrictions: 
// - Consist of 5 rows x 5 columns (25 cells).
// - All numbers in the card are different.
// - The central cell is empty.
// - For each column i (1 <= i <= 5), all cells have values into the
// range [15*(i-1)+1, 15*i].
type Card struct {
	cells [5][5]int
	marked [5][5]bool
	counter int
}

// NewRandomCard creates a new Card.
func NewRandomCard() *Card {
	var cells [5][5]int
	var marked [5][5]bool
	for i:=0; i<5; i++ {
		for j:=0; j<5; j++ {
			cells[i][j] = 0
			marked[i][j] = false
		}
	}
	return &Card{
		cells: cells,
		marked: marked,
		counter: 0,
	}
}

func (card *Card) isValidCell(r int, c int) bool {
	return (1 <= r && r <= 5) && (1 <= c && c <= 5) && !(r == 3 && c == 3)
}

func (card *Card) Value(r int, c int) int {
	return card.cells[r-1][c-1]
}

// Mark marks a cell.
func (card *Card) Mark(r int, c int) {
	if card.isValidCell(r, c) {
		if !card.marked[r-1][c-1] {
			card.marked[r-1][c-1] = true
			card.counter++
		}
	}
}

// Unmark unmarks a cell.
func (card *Card) Unmark(r int, c int) {
	if card.isValidCell(r, c) {
		if card.marked[r-1][c-1] {
			card.marked[r-1][c-1] = false
			card.counter--
		}
	}
}

func (card *Card) IsMarked(r int, c int) bool {
	return card.marked[r-1][c-1]
}