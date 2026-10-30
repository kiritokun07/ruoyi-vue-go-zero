package apierror

import (
	"context"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type Business string

func (e Business) Error() string { return string(e) }

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func HTTP(ctx context.Context, err error) (int, any) {
	var business Business
	if errors.As(err, &business) {
		return http.StatusOK, response{500, business.Error()}
	}
	logx.WithContext(ctx).Error(err)
	return http.StatusInternalServerError, response{500, "系统异常，请稍后重试"}
}
