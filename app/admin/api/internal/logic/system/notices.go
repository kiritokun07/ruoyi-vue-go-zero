package system

import (
	"strings"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

const noticeFrom = ` FROM sys_notice n LEFT JOIN sys_notice_read nr ON nr.notice_id=n.notice_id AND nr.user_id=? WHERE n.status='0'`
const noticeColumns = `n.notice_id,n.notice_title,n.notice_type,n.create_by,IFNULL(DATE_FORMAT(n.create_time,'%Y-%m-%d %H:%i:%s'),'') create_time,nr.user_id IS NOT NULL is_read`

func (s *Session) HeaderNotices() (*types.NoticeList, error) {
	rows, err := many[types.Notice](s, "SELECT "+noticeColumns+noticeFrom+" ORDER BY n.create_time DESC,n.notice_id DESC LIMIT 10", s.user)
	unread := 0
	for _, r := range rows {
		if !r.IsRead {
			unread++
		}
	}
	return &types.NoticeList{StatusResponse: ok(), Data: rows, UnreadCount: unread}, err
}

func (s *Session) NoticeDetail(r *types.RecordId) (*types.NoticeDetail, error) {
	from := noticeFrom
	if s.user == 1 || s.grants["system:notice:query"] || s.grants["*:*:*"] {
		from = strings.Replace(from, "n.status='0'", "1=1", 1)
	}
	n, err := one[types.Notice](s, "SELECT "+noticeColumns+",COALESCE(n.notice_content,'') notice_content"+from+" AND n.notice_id=?", s.user, r.Id)
	return &types.NoticeDetail{StatusResponse: ok(), Data: n}, err
}

func (s *Session) markRead(ids []int64) (*types.StatusResponse, error) {
	return saved(s.exec("INSERT INTO sys_notice_read(user_id,notice_id,read_time) SELECT ?,notice_id,NOW() FROM sys_notice WHERE status='0' AND notice_id IN ("+numbers(ids)+") ON DUPLICATE KEY UPDATE read_time=NOW()", s.user))
}

func (s *Session) MarkNoticeRead(r *types.ReadNoticeRequest) (*types.StatusResponse, error) {
	return s.markRead([]int64{r.NoticeId})
}

func (s *Session) MarkNoticesRead(r *types.ReadNoticesRequest) (*types.StatusResponse, error) {
	ids, e := parseIDs(r.Ids)
	if e != nil {
		return nil, e
	}
	return s.markRead(ids)
}
