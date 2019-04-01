package types

type responseType struct {
	Ok                  int
	Unauthorized        int
	InvalidArguments    int
	NotFound            int
	AlreadyExists       int
	InsufficientFunds   int
	PayPreviousPortions int
	PaymentCeiling      int
}

var Response = responseType{
	0,
	1,
	2,
	3,
	4,
	5,
	6,
	7,
}
