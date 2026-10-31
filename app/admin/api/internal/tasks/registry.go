// Package tasks provides explicitly registered Go implementations of RuoYi targets.
package tasks

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Registry struct {
	mu        sync.RWMutex
	functions map[string]reflect.Value
}

func NewRegistry() *Registry { return &Registry{functions: make(map[string]reflect.Value)} }

var Default = NewRegistry()

// Register accepts func(context.Context, ...typedArgs) (string, error).
// Supported argument types match RuoYi: string, bool, int, int64 (L), float64 (D).
func (r *Registry) Register(name string, fn any) error {
	value := reflect.ValueOf(fn)
	if !strings.Contains(name, ".") || strings.ContainsAny(name, " ()\t\n") || !value.IsValid() || value.Kind() != reflect.Func || value.IsNil() {
		return fmt.Errorf("invalid task registration: %s", name)
	}
	t := value.Type()
	if t.IsVariadic() || t.NumIn() < 1 || t.In(0) != reflect.TypeFor[context.Context]() || t.NumOut() != 2 || t.Out(0) != reflect.TypeFor[string]() || t.Out(1) != reflect.TypeFor[error]() {
		return fmt.Errorf("invalid task signature: %s", name)
	}
	for i := 1; i < t.NumIn(); i++ {
		switch t.In(i) {
		case reflect.TypeFor[string](), reflect.TypeFor[bool](), reflect.TypeFor[int](), reflect.TypeFor[int64](), reflect.TypeFor[float64]():
		default:
			return fmt.Errorf("unsupported task parameter: %s", name)
		}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.functions[name]; ok {
		return fmt.Errorf("task already registered: %s", name)
	}
	r.functions[name] = value
	return nil
}

func (r *Registry) resolve(target string) (reflect.Value, []reflect.Value, error) {
	name, args, hasArgs := strings.Cut(strings.TrimSpace(target), "(")
	name = strings.TrimSpace(name)
	r.mu.RLock()
	fn, ok := r.functions[name]
	r.mu.RUnlock()
	if !ok {
		return reflect.Value{}, nil, fmt.Errorf("任务调用目标未注册")
	}
	var values []reflect.Value
	if hasArgs {
		if !strings.HasSuffix(args, ")") {
			return fn, nil, fmt.Errorf("任务参数格式无效")
		}
		args = strings.TrimSpace(strings.TrimSuffix(args, ")"))
		if args != "" {
			quoted, start := false, 0
			var tokens []string
			for i, c := range args {
				if c == '\'' {
					quoted = !quoted
				}
				if c == ',' && !quoted {
					tokens = append(tokens, args[start:i])
					start = i + 1
				}
			}
			if quoted {
				return fn, nil, fmt.Errorf("任务字符串参数未闭合")
			}
			tokens = append(tokens, args[start:])
			for _, token := range tokens {
				token = strings.TrimSpace(token)
				var v any
				var err error
				switch {
				case len(token) >= 2 && token[0] == '\'' && token[len(token)-1] == '\'':
					if strings.Contains(token[1:len(token)-1], "'") {
						return fn, nil, fmt.Errorf("任务字符串参数无效")
					}
					v = token[1 : len(token)-1]
				case token == "true" || token == "false":
					v = token == "true"
				case strings.HasSuffix(strings.ToUpper(token), "L"):
					v, err = strconv.ParseInt(token[:len(token)-1], 10, 64)
				case strings.HasSuffix(strings.ToUpper(token), "D"):
					v, err = strconv.ParseFloat(token[:len(token)-1], 64)
				default:
					var n int64
					n, err = strconv.ParseInt(token, 10, 32)
					v = int(n)
				}
				if err != nil {
					return fn, nil, fmt.Errorf("任务参数无效")
				}
				values = append(values, reflect.ValueOf(v))
			}
		}
	}
	if fn.Type().NumIn() != len(values)+1 {
		return fn, nil, fmt.Errorf("任务参数数量不匹配")
	}
	for i, v := range values {
		if fn.Type().In(i+1) != v.Type() {
			return fn, nil, fmt.Errorf("任务参数类型不匹配")
		}
	}
	return fn, values, nil
}
func (r *Registry) Validate(target string) error { _, _, err := r.resolve(target); return err }
func (r *Registry) Invoke(ctx context.Context, target string) (message string, err error) {
	fn, args, err := r.resolve(target)
	if err != nil {
		return "", err
	}
	if err = ctx.Err(); err != nil {
		return "", err
	}
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("task panic: %v", v)
		}
	}()
	values := fn.Call(append([]reflect.Value{reflect.ValueOf(ctx)}, args...))
	message = values[0].String()
	if !values[1].IsNil() {
		err = values[1].Interface().(error)
	}
	return
}
