package tasks

import (
	"context"
	"fmt"
)

func init() {
	registrations := map[string]any{
		"ryTask.ryNoParams": func(context.Context) (string, error) { return "执行无参方法", nil },
		"ryTask.ryParams":   func(_ context.Context, value string) (string, error) { return "执行有参方法：" + value, nil },
		"ryTask.ryMultipleParams": func(_ context.Context, s string, b bool, l int64, d float64, i int) (string, error) {
			return fmt.Sprintf("执行多参方法：字符串 %s，布尔 %t，长整型 %d，浮点型 %g，整形 %d", s, b, l, d, i), nil
		},
	}
	for name, fn := range registrations {
		if err := Default.Register(name, fn); err != nil {
			panic(err)
		}
	}
}
