package client

// Logger интерфейс для логирования
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

// Field представляет поле для структурированного логирования
type Field struct {
	Key   string
	Value interface{}
}

// NoOpLogger реализует Logger без вывода (по умолчанию)
type NoOpLogger struct{}

func (n *NoOpLogger) Debug(msg string, fields ...Field) {}
func (n *NoOpLogger) Info(msg string, fields ...Field)  {}
func (n *NoOpLogger) Warn(msg string, fields ...Field)  {}
func (n *NoOpLogger) Error(msg string, fields ...Field) {}

// SetLogger устанавливает логгер для клиента
func (c *Client) SetLogger(logger Logger) {
	if logger == nil {
		c.logger = &NoOpLogger{}
	} else {
		c.logger = logger
	}
}

