package example

func IfNotNilPanic(err error) {
	if err != nil {
		panic(err)
	}
}
