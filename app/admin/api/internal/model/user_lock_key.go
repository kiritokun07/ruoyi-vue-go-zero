package model

import (
	"context"
	"crypto/sha256"
	"fmt"
	"regexp"
)

var sqlIdentifier = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// UserLockKey uses the existing column collation, so names considered equal by
// MySQL (including case/accent variants) also share a Redis lock. No schema change.
func (m *AdminModel) UserLockKey(ctx context.Context, name string) (string, error) {
	var column struct {
		Schema    string `db:"schema_name"`
		Charset   string `db:"charset_name"`
		Collation string `db:"collation_name"`
	}
	err := m.db.QueryRowCtx(ctx, &column, "SELECT TABLE_SCHEMA schema_name, CHARACTER_SET_NAME charset_name, COLLATION_NAME collation_name FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='sys_user' AND COLUMN_NAME='user_name'")
	if err != nil {
		return "", err
	}
	if !sqlIdentifier.MatchString(column.Charset) || !sqlIdentifier.MatchString(column.Collation) {
		return "", fmt.Errorf("invalid username collation")
	}
	var weight string
	err = m.db.QueryRowCtx(ctx, &weight, "SELECT HEX(WEIGHT_STRING(RTRIM(CONVERT(? USING "+column.Charset+")) COLLATE "+column.Collation+"))", name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("ruoyi-go:user-write:%x", sha256.Sum256([]byte(column.Schema+":"+weight))), nil
}
