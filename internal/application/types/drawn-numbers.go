package types

type DrawnNumbers struct {
	
}

func NewEmptyDrawnNumbers() *DrawnNumbers {
	return &DrawnNumbers{}
}

func (dn *DrawnNumbers) Contains(number int) bool {
	return true
}