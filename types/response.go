package types

type responseType struct {
	Ok               int
	Unauthorized     int
	InvalidArguments int
	NotFound         int
}

var Response = responseType{0, 1, 2, 3}

