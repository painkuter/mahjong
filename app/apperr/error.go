package apperr

import "mahjong/app/common/log"

func Check(err error) {
	if err != nil {
		//TODO: logging
		log.Error(err.Error())
		panic(err)
	}
}
