// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package role

import (
	"context"
	"ruoyi-vue-go-zero/app/admin/api/internal/logic/system"

	"ruoyi-vue-go-zero/app/admin/api/internal/svc"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMembersLogic {
	return &AddMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMembersLogic) AddMembers(req *types.MemberAssignment) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:role:edit", req, func(s *system.Session) (*types.StatusResponse, error) { return s.AddMembers(req) })
}
