package system

import (
	"net/http/httptest"
	"ruoyi-vue-go-zero/app/admin/api/internal/tasks"
	"strings"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func TestOnlyLocalTableDDL(t *testing.T) {
	for _, stmt := range []string{"DROP TABLE sys_user", "CREATE TABLE otherdb.example (id bigint)", "CREATE TABLE copy AS SELECT * FROM sys_user", "CREATE VIEW v AS SELECT 1"} {
		if _, _, e := checkedDDL(stmt); e == nil {
			t.Fatalf("accepted: %s", stmt)
		}
	}
	n, q, e := checkedDDL("CREATE TABLE biz_sample (id bigint PRIMARY KEY) ENGINE=MyISAM COMMENT='业务表'")
	if e != nil || n != "biz_sample" || !strings.Contains(q, "ENGINE=InnoDB") || !strings.Contains(q, "业务表") {
		t.Fatalf("name %q, sql %q, err %v", n, q, e)
	}
}

func TestQuartzDayAndSeconds(t *testing.T) {
	start := time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		expr string
		want time.Time
	}{
		{"0/10 * * * * ?", start.Add(10 * time.Second)},
		{"0 0 8 ? * 2", time.Date(2026, 11, 2, 8, 0, 0, 0, time.UTC)},
		{"0 0 8 ? * 2#2", time.Date(2026, 11, 9, 8, 0, 0, 0, time.UTC)},
	} {
		c, e := schedule(tc.expr)
		if e != nil {
			t.Fatal(e)
		}
		if got := c.Next(start); !got.Equal(tc.want) {
			t.Errorf("%s: got %v want %v", tc.expr, got, tc.want)
		}
	}
	if e := tasks.Default.Validate("os.exec('echo unexpected')"); e == nil {
		t.Fatal("unregistered target allowed")
	}
}

func TestCreateTaskRequestDoesNotRequireComputedNextFire(t *testing.T) {
	req := httptest.NewRequest("POST", "/monitor/job", strings.NewReader(`{"jobName":"probe","jobGroup":"DEFAULT","invokeTarget":"ryTask.ryNoParams","cronExpression":"0/10 * * * * ?","concurrent":"1"}`))
	req.Header.Set("Content-Type", "application/json")
	var task types.ScheduleEntry
	if err := httpx.Parse(req, &task); err != nil {
		t.Fatal(err)
	}
}
