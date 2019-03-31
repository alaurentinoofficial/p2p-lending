package types

type statementType struct {
	In int
	Out int
	Dividend int
}

func (obj *statementType) Check(_type int) bool {
	return _type >= 0 && _type <= 3
}

var Statement = statementType{ In: 0, Out: 1, Dividend: 3 }