package tasks

import (
	"context"
	"errors"
	"testing"
)

func TestRegisteredTasks(t *testing.T) {
	r := NewRegistry()
	called := 0
	err := r.Register("orders.sync", func(ctx context.Context, s string, b bool, l int64, d float64, n int) (string, error) {
		called++
		if s != "a,b" || !b || l != 42 || d != 1.5 || n != 3 {
			t.Fatal("wrong arguments")
		}
		return "done", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	target := "orders.sync('a,b', true, 42L, 1.5D, 3)"
	if err = r.Validate(target); err != nil {
		t.Fatal(err)
	}
	if called != 0 {
		t.Fatal("validation executed task")
	}
	msg, err := r.Invoke(context.Background(), target)
	if err != nil || msg != "done" || called != 1 {
		t.Fatal(msg, err, called)
	}
	for _, target := range []string{"os.exec('ls')", "orders.sync()", "orders.sync('unclosed)", "orders.sync('a,b', true, 42, 1.5D, 3)"} {
		if r.Validate(target) == nil {
			t.Fatal("accepted invalid target", target)
		}
	}
	if r.Register("orders.sync", func(context.Context) (string, error) { return "", nil }) == nil {
		t.Fatal("duplicate registration")
	}
	if r.Register("bad.signature", func() {}) == nil {
		t.Fatal("invalid signature")
	}
	if err = r.Register("jobs.fail", func(context.Context) (string, error) { return "", errors.New("failed") }); err != nil {
		t.Fatal(err)
	}
	if _, err = r.Invoke(context.Background(), "jobs.fail"); err == nil {
		t.Fatal("lost error")
	}
	if err = r.Register("jobs.panic", func(context.Context) (string, error) { panic("test") }); err != nil {
		t.Fatal(err)
	}
	if _, err = r.Invoke(context.Background(), "jobs.panic"); err == nil {
		t.Fatal("lost panic")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err = r.Invoke(ctx, target); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
