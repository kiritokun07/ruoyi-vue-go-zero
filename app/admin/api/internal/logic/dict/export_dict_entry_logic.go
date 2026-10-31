// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package dict

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"
	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type ExportDictEntryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportDictEntryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportDictEntryLogic {
	return &ExportDictEntryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportDictEntryLogic) ExportDictEntry(req *types.AdminFilter) (*system.File, error) {
	return system.Read(l.ctx, l.svcCtx, "system:dict:export", func(s *system.Session) (*system.File, error) { return s.ExportDictEntry(req) })
}
