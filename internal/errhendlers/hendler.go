package errhendlers

import (
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/logger"
)

func PanicErr(err error, msg string) {
	if err != nil {
		logger.Error(msg, nil, err)
		panic(err)
	}
}
