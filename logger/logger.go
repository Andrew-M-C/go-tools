package logger

type L interface{
    Debug(v ...interface{})
    Debugf(format string, a ...interface{})
    Info(v ...interface{})
    Infof(format string, a ...interface{})
    Warn(v ...interface{})
    Warnf(format string, a ...interface{})
    Error(v ...interface{})
    Errorf(format string, a ...interface{})
}
