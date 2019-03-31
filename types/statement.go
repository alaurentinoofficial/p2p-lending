package types

type statementType struct {
	In int
	Out int
	Dividend int
}

var Statement = statementType{ In: 0, Out: 1, Dividend: 3 }