package common

func Check(err error) {
	if err != nil {
		//TODO: logging
		panic(err)
	}
}
