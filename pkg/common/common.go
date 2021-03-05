package common

import (
	"fmt"
	"log"
	"runtime"
)

//攔截Panic，防止非預期錯誤導致當機
func CatchPanic(errorTitle string) {
	if err := recover(); err != nil {
		var logText string
		for index := 0; index < 8; index++ { //最多捕捉10層log
			ptr, filename, line, ok := runtime.Caller(index)
			if !ok {
				break
			}
			logText = logText + fmt.Sprintf(" %v:%d,%v > %v\n", filename, line, runtime.FuncForPC(ptr).Name(), err)
		}
		log.Printf("%v : 發生嚴重錯誤 %v\n", errorTitle, logText)
		return
	}
}
