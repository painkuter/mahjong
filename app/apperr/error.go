package apperr

import "fmt"

func Check(err error) {
	if err != nil {
		//TODO: logging
		fmt.Println(err.Error())
		panic(err)
	}
}
