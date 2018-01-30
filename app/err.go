package app

func check(err error) {
	if err != nil {
		//TODO: logging
		panic(err)
	}
}
