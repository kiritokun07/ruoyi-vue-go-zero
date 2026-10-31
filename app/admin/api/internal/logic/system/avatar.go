package system

import (
	"bytes"
	"context"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

// SaveAvatar decodes and re-encodes the image before exposing it as a JPEG.
func SaveAvatar(ctx context.Context, services *svc.ServiceContext, file io.Reader) (*types.AvatarResult, error) {
	const limit = 5 << 20
	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return nil, err
	}
	if len(data) > limit {
		return nil, apierror.Business("头像不能超过5 MiB")
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width <= 0 || cfg.Height <= 0 || cfg.Width > 4096 || cfg.Height > 4096 {
		return nil, apierror.Business("头像格式或尺寸无效")
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, apierror.Business("头像图片无效")
	}
	dir := services.Config.AvatarDir
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	name := uuid.NewString() + ".jpg"
	path := filepath.Join(dir, name)
	out, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	err = jpeg.Encode(out, img, &jpeg.Options{Quality: 88})
	closeErr := out.Close()
	if err != nil || closeErr != nil {
		os.Remove(path)
		if err != nil {
			return nil, err
		}
		return nil, closeErr
	}
	link := "/profile/avatar/" + name
	res, err := withSession(ctx, services, "", true, func(s *Session) (*types.AvatarResult, error) {
		return &types.AvatarResult{StatusResponse: ok(), ImgUrl: link}, s.exec("UPDATE sys_user SET avatar=?,update_by=?,update_time=NOW() WHERE user_id=?", link, s.name, s.user)
	})
	if err != nil {
		os.Remove(path)
	}
	return res, err
}
