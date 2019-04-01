package types

type indexType struct {
	SELIC int
	IPCA  int
	CDI   int
}

func (obj *indexType) Check(_type int) bool {
	return _type >= 0 && _type <= 3
}

var Index = indexType{
	SELIC: 0,
	IPCA:  1,
	CDI:   3,
}

func (obj *indexType) Porcentage(index int) float32 {
	switch index {
	case obj.SELIC:
		return 6.4
	case obj.IPCA:
		return 3.86
	case obj.CDI:
		return 6.4
	default:
		return 1
	}
}
