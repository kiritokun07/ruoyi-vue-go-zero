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

type RemoveMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveMembersLogic {
	return &RemoveMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveMembersLogic) RemoveMembers(req *types.MemberAssignment) (resp *types.StatusResponse, err error) {
	return system.Write(l.ctx, l.svcCtx, "system:role:edit", req, func(s *system.Session) (*types.StatusResponse, error) { return s.RemoveMembers(req) })
}
