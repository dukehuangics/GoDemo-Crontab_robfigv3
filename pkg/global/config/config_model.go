package config

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
量測並建立SettingMap
@param needClear: 是否強制清除SettingMap
*/
func CreatSettingMap(needClear bool) {
	if needClear || len(SettingMap)<=0{
		SettingMap = make(map[string]*AddSetings) //儲存環境變數,可依需求自行擴充
	}
	return
}

/*
讀取環境變數設定檔
@param configPath: 設定檔路徑
@param isInit: 是否第首次執行(遇錯則關閉整個程序)
*/
func (cfgs *CfgType) ReadFile(configPath string, isInit bool) {

	if !isInit || configPath == "" {
		configPath = Cfgs.ConfigPath
	}

	configBtye, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("讀取環境設定值失敗: %v (%v)\n", configPath, err)
		if isInit {
			os.Exit(Code_error_exit)
		}
	}

	Cfgs.ConfigPath = configPath

	// →#註解← 一般注解方式
	// →  #註解  ←設定值前後可留半形空白
	// →settingKey1=value1← 一般設定方式
	// →settingKey2 = value1← 錯誤設定方式

	for index, value := range strings.Split(string(configBtye), "\n") {

		value = strings.Trim(value, " ") //清除前後多餘空白
		indexComment := strings.Index(value, "#")
		if indexComment >= 0 && indexComment < 3 {
			continue //此#為設定的註解，故忽略
		}

		if indexSetChar := strings.Index(value, "="); indexSetChar < 1 {
			continue //此狀況設定值"="位置錯誤或無"=""
		} else {

			settingKey := value[:indexSetChar]
			settingValue := value[indexSetChar+1:]

			if _, isExsit := SettingMap[settingKey]; !isExsit {
				log.Printf("環境設定值錯誤警告: %v (Line:%v) > %v，檢測到無效設定!", configPath, index, value)
				continue
			}

			if err := SettingMap[settingKey].Run(settingValue); err != nil {
				log.Printf("環境設定值錯誤警告: %v (Line:%v) > %v (%v)，將使用預設值'%v'代替.", configPath, index, value, err, SettingMap[settingKey].DefaultValue)
			}

		}
	}
}
