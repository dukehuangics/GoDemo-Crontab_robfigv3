package printExample

import (
	"crontab/pkg/common"
	"crontab/pkg/common/customVar"
	"crontab/pkg/global/config"
	"log"
)


var index int

func init() {
	config.CreatSettingMap(false)
	config.SettingMap["print_example_index"] = &config.AddSetings{DefaultValue: 1, OutputPointer: &index, Custom: &customVar.SecondsInADay{}}
}

func StartJob() {
	defer common.CatchPanic("打印測試")
	index++
	log.Println("這是打印測試!", index)
	return
}
