package apierror

import (
	"context"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

var Unauthorized = errors.New("认证失败，请重新登录")

var Forbidden = errors.New("没有权限访问此资源")

type Business string

func (e Business) Error() string { return string(e) }

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func HTTP(ctx context.Context, err error) (int, any) {
	if errors.Is(err, Unauthorized) {
		return http.StatusOK, response{401, Unauthorized.Error()}
	}
	if errors.Is(err, Forbidden) {
		return http.StatusOK, response{403, Forbidden.Error()}
	}
	var business Business
	if errors.As(err, &business) {
		return http.StatusOK, response{500, business.Error()}
	}
	logx.WithContext(ctx).Error(err)
	return http.StatusInternalServerError, response{500, "系统异常，请稍后重试"}
}
