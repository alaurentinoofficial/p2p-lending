package types

type indexType struct {
	SELIC int
	IPCA int
	CDI int
}

func (obj *indexType) Check(_type int) bool {
	return _type >= 0 && _type <= 3
}

var Index = indexType{
	SELIC: 0,
	IPCA: 1,
	CDI: 3,
}
