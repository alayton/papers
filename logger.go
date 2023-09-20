package papers

import (
	"fmt"
)

type Logger interface {
	Print(v ...interface{})
}

type SystemLogger struct{}

func (s SystemLogger) Print(v ...interface{}) {
	fmt.Println(v...)
}
