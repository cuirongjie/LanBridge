/*
简易日志
*/
package logger

import (
	"LanBridge_Client/cache"
	"log"
)

func Debug(a ...any) {
	if cache.Conf.Debug {
		log.Println(a)
	}
}

func Info(a ...any) {
	log.Println(a)
}
