package conf

import (
	"gitee.com/mrmateoliu/wow_launch.git/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

// --------------------------------------------------------
// -初始化日志组件
func InitLogger() *zap.SugaredLogger {
	var writeSyncer zapcore.WriteSyncer

	logMode := zapcore.DebugLevel
	if !utils.GlobalObject.Develop {
		logMode = zapcore.InfoLevel

		// 如果 Develop 为假，只输出到文件
		writeSyncer = getWriteSyncer()
	} else {
		// 如果 Develop 为真，将日志同时输出到文件和控制台
		writeSyncer = zapcore.NewMultiWriteSyncer(getWriteSyncer(), zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(getEncoder(), writeSyncer, logMode)

	return zap.New(core).Sugar()
}

// --------------------------------------------------------
// -设置日志编码规范
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Local().Format(time.DateTime))
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

// --------------------------------------------------------
// -获取日志写入方式
func getWriteSyncer() zapcore.WriteSyncer {
	sSeparator := string(filepath.Separator) //取分隔符
	sRootDir, _ := os.Getwd()                //取当前文件目录
	// 当前文件目录 + 分隔符 + log文件夹 + 分隔符 + 当前时间 + .txt
	sLogFilePath := sRootDir + sSeparator + "log" + sSeparator + time.Now().Format(time.DateOnly) + ".txt"

	lumberjackSyncer := &lumberjack.Logger{
		Filename:   sLogFilePath,
		MaxSize:    int(utils.GlobalObject.LogMaxSize),    //日志文件最大尺寸(M),超限后自动分隔
		MaxBackups: int(utils.GlobalObject.LogMaxBackups), //保留旧文件的最大个数
		MaxAge:     int(utils.GlobalObject.LogMaxAge),     //保留旧文件的最大天数
		Compress:   false,                                 // disabled by default
	}

	return zapcore.AddSync(lumberjackSyncer)
}
