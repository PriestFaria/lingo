package zap

// Stub-реализация go.uber.org/zap для тестов analysistest.

type Logger struct{}

func NewNop() *Logger { return &Logger{} }

func (l *Logger) Info(msg string, fields ...Field)  {}
func (l *Logger) Warn(msg string, fields ...Field)  {}
func (l *Logger) Error(msg string, fields ...Field) {}
func (l *Logger) Debug(msg string, fields ...Field) {}
func (l *Logger) Sugar() *SugaredLogger            { return &SugaredLogger{} }
func (l *Logger) With(fields ...Field) *Logger     { return l }

type SugaredLogger struct{}

func (s *SugaredLogger) Infof(template string, args ...interface{})  {}
func (s *SugaredLogger) Warnf(template string, args ...interface{})  {}
func (s *SugaredLogger) Errorf(template string, args ...interface{}) {}
func (s *SugaredLogger) Debugf(template string, args ...interface{}) {}

type Field struct{}

func String(key, val string) Field { return Field{} }

























