package model

import (
	"context"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const BulletinColumns = `IFNULL(notice_id,0) notice_id,IFNULL(notice_title,'') notice_title,IFNULL(notice_type,'') notice_type,IFNULL(notice_content,'') notice_content,IFNULL(status,'') status,IFNULL(create_by,'') create_by,IFNULL(DATE_FORMAT(create_time,'%Y-%m-%d %H:%i:%s'),'') create_time,IFNULL(update_by,'') update_by,IFNULL(DATE_FORMAT(update_time,'%Y-%m-%d %H:%i:%s'),'') update_time,IFNULL(remark,'') remark`

func (m *AdminModel) BulletinRows(ctx context.Context, r *types.AdminFilter, paged bool) ([]types.BulletinEntry, int64, error) {
	f := filter{}
	f.text("notice_title", r.NoticeTitle)
	f.equal("notice_type", r.NoticeType)
	f.text("create_by", r.CreateBy)
	f.dates("create_time", &types.ListFilter{BeginTime: r.BeginTime, EndTime: r.EndTime})
	from := " FROM sys_notice" + f.where()
	n, err := queryOne[int64](ctx, m.db, "SELECT COUNT(*)"+from, f.args...)
	if err != nil {
		return nil, 0, err
	}
	order := " ORDER BY notice_id DESC"
	if paged {
		order += page(&types.ListFilter{PageNum: r.PageNum, PageSize: r.PageSize})
	}
	rows, err := queryRows[types.BulletinEntry](ctx, m.db, "SELECT "+BulletinColumns+from+order, f.args...)
	return rows, n, err
}

func (m *AdminModel) DeleteBulletin(ctx context.Context, ids []int64) error {
	return m.exec(ctx, "DELETE FROM sys_notice WHERE notice_id IN ("+numbers(ids)+")")
}

func (m *AdminModel) EditBulletin(ctx context.Context, r *types.BulletinEntry, actor string) error {
	return m.exec(ctx, "UPDATE sys_notice SET notice_title=?,notice_type=?,notice_content=?,status=?,remark=?,update_by=?,update_time=NOW() WHERE notice_id=?", r.NoticeTitle, r.NoticeType, r.NoticeContent, r.Status, r.Remark, actor, r.NoticeId)
}

func (m *AdminModel) EditBulletinExists(ctx context.Context, id int64) (int, error) {
	return queryOne[int](ctx, m.db, "SELECT COUNT(*) FROM sys_notice WHERE notice_id=?", id)
}

func (m *AdminModel) CreateBulletin(ctx context.Context, r *types.BulletinEntry, actor string) error {
	return m.exec(ctx, "INSERT INTO sys_notice(notice_title,notice_type,notice_content,status,remark,create_by,create_time) VALUES(?,?,?,?,?,?,NOW())", r.NoticeTitle, r.NoticeType, r.NoticeContent, r.Status, r.Remark, actor)
}
