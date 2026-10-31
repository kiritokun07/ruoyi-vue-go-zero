package online

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"testing"
	"time"
)

func TestSessionSlidingExpiryAndRevocation(t *testing.T) {
	server := miniredis.RunT(t)
	cache := redis.New(server.Addr())
	ctx := context.Background()
	session := Session{UserID: 7, TokenId: "session", Identity: Identity{UserName: "member"}}
	if err := Create(ctx, cache, session, 1800); err != nil {
		t.Fatal(err)
	}
	server.FastForward(5 * time.Minute)
	if ok, err := Touch(ctx, cache, "session", 7, 1800); err != nil || !ok {
		t.Fatalf("touch: %v %v", ok, err)
	}
	if ttl := server.TTL(Prefix + "session"); ttl != 25*time.Minute {
		t.Fatalf("renewed too early: %v", ttl)
	}
	server.FastForward(6 * time.Minute)
	if ok, err := Touch(ctx, cache, "session", 7, 1800); err != nil || !ok {
		t.Fatal(ok, err)
	}
	if ttl := server.TTL(Prefix + "session"); ttl != 30*time.Minute {
		t.Fatalf("not renewed: %v", ttl)
	}
	if ok, err := Touch(ctx, cache, "session", 8, 1800); err != nil || ok {
		t.Fatal("wrong account accepted", err)
	}
	if _, err := Revoke(ctx, cache, "session", 0); err != nil {
		t.Fatal(err)
	}
	if ok, err := Touch(ctx, cache, "session", 7, 1800); err != nil || ok {
		t.Fatal("revoked token accepted", err)
	}
	if server.Exists(Prefix + "session") {
		t.Fatal("session recreated")
	}
	session.TokenId = "missing"
	if err := Create(ctx, cache, session, 1800); err != nil {
		t.Fatal(err)
	}
	server.Del(Prefix + "missing")
	if ok, err := Touch(ctx, cache, "missing", 7, 1800); err != nil || ok {
		t.Fatal("missing session accepted", err)
	}
	session.TokenId = "idle"
	if err := Create(ctx, cache, session, 1800); err != nil {
		t.Fatal(err)
	}
	server.FastForward(31 * time.Minute)
	if ok, err := Touch(ctx, cache, "idle", 7, 1800); err != nil || ok {
		t.Fatal("expired session accepted", err)
	}
}
