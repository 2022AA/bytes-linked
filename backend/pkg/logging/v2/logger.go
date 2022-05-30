package v2

import (
	"context"
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	KeyUidWithCtx = "uid"

	//以下变量包外可见，用于测试场景下拦截日志输出
	ZapLogger        *zap.Logger
	ZapSugaredLogger *zap.SugaredLogger
	ZapCfg           *Config

	AtomLv zap.AtomicLevel
)

type Level = zapcore.Level

const (
	typeAny = iota
	typeString
	typeBool
	typeUint32
	typeInt32
	typeError
	typeStack
)

const (
	DebugLevel = zap.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zap.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zap.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zap.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = zap.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zap.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zap.FatalLevel
)

func init() {
	Init(&Config{Level: InfoLevel})
}

type Config struct {
	Level Level
}

type logField struct {
	ty        int
	key       string
	value     interface{}
	stringVal string
	err       error
	integer   int64
	boolFlag  bool
}

func Field(key string, value interface{}) *logField {
	return &logField{key: key, value: value}
}

func String(key, value string) *logField {
	return &logField{key: key, ty: typeString, stringVal: value}
}

func Bool(key string, val bool) *logField {
	return &logField{key: key, ty: typeBool, boolFlag: val}
}

func Uint32(key string, val uint32) *logField {
	return &logField{key: key, ty: typeUint32, integer: int64(val)}
}

func Int32(key string, val int32) *logField {
	return &logField{key: key, ty: typeInt32, integer: int64(val)}
}

func ErrorField(err error) *logField {
	return &logField{ty: typeError, err: err}
}

func Stack(key string) *logField {
	return &logField{key: key, ty: typeError}
}

func Init(cfg *Config) {

	lv := zap.NewAtomicLevelAt(cfg.Level)

	AtomLv = lv
	ZapCfg = cfg

	encodeCfg := zap.NewProductionEncoderConfig()
	encodeCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encodeCfg), zapcore.AddSync(os.Stdout), zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= lv.Level()
	}))
	ZapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.NewAtomicLevelAt(zap.PanicLevel)))
	ZapSugaredLogger = ZapLogger.Sugar()
	//debug use ...
	//if lv.Level() <= DebugLevel {
	//	grpclog.SetLoggerV2(newGrpcLogger(ZapSugaredLogger))
	//}

	//SetLevelHandler(cfg.Port)

}

//func SetLevelHandler(port int) {
//	//Level handler
//	if port == 0 {
//		port = ZapCfg.Port
//	}
//
//	go levelHandler(port, AtomLv)
//}

func SetLevel(lv zapcore.Level) {
	AtomLv.SetLevel(lv)

}

//func levelHandler(port int, lv zap.AtomicLevel) {
//	if port == 0 {
//		ZapLogger.Info("levelHandler ignore Level handler ")
//		return
//	}
//	http.HandleFunc("/", lv.ServeHTTP)
//	Infof("logLevelHandler listen on %d", port)
//	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
//		ZapLogger.Error("levelHandler listen(%d) err %+v", zap.Int("Port", port), zap.Error(err))
//	}
//
//}

//---------------------------log default format---------------------------------------
func Debugf(format string, args ...interface{}) {
	ZapSugaredLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	ZapSugaredLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	ZapSugaredLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	ZapSugaredLogger.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	ZapSugaredLogger.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	ZapSugaredLogger.Fatalf(format, args...)
}

//-------------------------log ln-----------------------------------------

func Debugln(args ...interface{}) {
	ZapSugaredLogger.Debug(args...)
}

func Infoln(args ...interface{}) {
	ZapSugaredLogger.Info(args...)
}

func Warnln(args ...interface{}) {
	ZapSugaredLogger.Warn(args...)
}

func Errorln(args ...interface{}) {
	ZapSugaredLogger.Error(args...)
}

func Panicln(args ...interface{}) {
	ZapSugaredLogger.Panic(args...)
}

func Fatalln(args ...interface{}) {
	ZapSugaredLogger.Fatal(args...)
}

//------------------------------------------------------------------

func convertField(ctx context.Context, keyValues ...*logField) []zap.Field {
	fields := make([]zap.Field, 0, len(keyValues))
	for _, keyValue := range keyValues {
		switch keyValue.ty {
		case typeString:
			fields = append(fields, zap.String(keyValue.key, keyValue.stringVal))
		case typeUint32:
			fields = append(fields, zap.Uint32(keyValue.key, uint32(keyValue.integer)))
		case typeInt32:
			fields = append(fields, zap.Int32(keyValue.key, int32(keyValue.integer)))
		case typeBool:
			fields = append(fields, zap.Bool(keyValue.key, keyValue.boolFlag))
		case typeError:
			fields = append(fields, zap.Error(keyValue.err))
		case typeStack:
			fields = append(fields, zap.Stack(keyValue.key))
		default:
			fields = append(fields, zap.Any(keyValue.key, keyValue.value))
		}

	}

	if uid := uidFromContext(ctx); uid != 0 {
		fields = append(fields, zap.String(KeyUidWithCtx, strconv.Itoa(int(uid))))
	}

	return fields
}

//----------------------------log with keyValues--------------------------------------

func Debug(ctx context.Context, msg string, keyValues ...*logField) {
	fields := convertField(ctx, keyValues...)
	ZapLogger.Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, keyValues ...*logField) {
	fields := convertField(ctx, keyValues...)
	ZapLogger.Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, keyValues ...*logField) {
	fields := convertField(ctx, keyValues...)
	ZapLogger.Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, keyValues ...*logField) {
	fields := convertField(ctx, keyValues...)
	ZapLogger.Error(msg, fields...)
}

func Panic(ctx context.Context, msg string, keyValues ...*logField) {
	fields := convertField(ctx, keyValues...)
	ZapLogger.Panic(msg, fields...)
}

func Fatal(ctx context.Context, msg string, keyValues ...*logField) {
	fields := convertField(ctx, keyValues...)
	ZapLogger.Fatal(msg, fields...)
}

func uidFromContext(ctx context.Context) uint32 {
	//if ctx == nil {
	//	return 0
	//}
	//if info, ok := gwcontext.ServiceInfoFromContext(ctx); ok {
	//	return info.UserID
	//}

	return 0

}

func sugaredLogFromCtx(ctx context.Context) *zap.SugaredLogger {
	if uid := uidFromContext(ctx); uid != 0 {
		return ZapSugaredLogger.With(zap.String(KeyUidWithCtx, strconv.Itoa(int(uid))))
	}
	return ZapSugaredLogger
}

//----------------------------log with ctx format --------------------------------------

func DebugfWithCtx(ctx context.Context, format string, args ...interface{}) {
	sugaredLogFromCtx(ctx).Debugf(format, args...)
}

func InfoWithCtx(ctx context.Context, format string, args ...interface{}) {
	sugaredLogFromCtx(ctx).Infof(format, args...)
}

func WarnWithCtx(ctx context.Context, format string, args ...interface{}) {
	sugaredLogFromCtx(ctx).Warnf(format, args...)
}

func ErrorWithCtx(ctx context.Context, format string, args ...interface{}) {
	sugaredLogFromCtx(ctx).Errorf(format, args...)
}

func PanicWithCtx(ctx context.Context, format string, args ...interface{}) {
	sugaredLogFromCtx(ctx).Panicf(format, args...)
}

func FatalWithCtx(ctx context.Context, format string, args ...interface{}) {
	sugaredLogFromCtx(ctx).Fatalf(format, args...)

}

func ParseLevel(lvl string) (Level, error) {
	var l Level
	err := l.UnmarshalText([]byte(lvl))
	if err != nil {
		return InfoLevel, err
	}
	return l, nil
}
