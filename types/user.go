package types

type userType struct {
	Physical int
	Legal int
}

func (obj *userType) Check(_type int) bool {
	return _type == 0 || _type == 1
}

var User = userType{ Physical: 0, Legal: 1 }