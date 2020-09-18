package main

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

func initLogger() (err error) {
	//`{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`
	config := make(map[string]interface{})
	config["filename"] = AppConfig.Log.LogPath
	config["level"]    = convertLogLevel(AppConfig.Log.LogLevel)
	config["maxlines"] = 10000000
	//config["maxsize"]  = "256MB"

	bytes, err := json.Marshal(config)
	if err != nil {
		return
	}
	// 日志记录调用的文件名和文件行号
	logs.EnableFuncCallDepth(true)
	// 自定义log日志的记录方式
	return logs.SetLogger(logs.AdapterFile, string(bytes))
}

func convertLogLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}
