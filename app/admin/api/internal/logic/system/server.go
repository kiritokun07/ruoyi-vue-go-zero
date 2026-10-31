package system

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"

	"ruoyi-vue-go-zero/app/admin/api/internal/types"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

var startedAt = time.Now()

func (s *Session) ServerInfo() (*types.Result, error) {
	vm, e := mem.VirtualMemoryWithContext(s.ctx)
	if e != nil {
		return nil, e
	}
	times, e := cpu.TimesWithContext(s.ctx, false)
	if e != nil {
		return nil, e
	}
	c := times[0]
	total := c.Total()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	host, _ := os.Hostname()
	cwd, _ := os.Getwd()
	ip := "127.0.0.1"
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if n, ok := addr.(*net.IPNet); ok && !n.IP.IsLoopback() && n.IP.To4() != nil {
			ip = n.IP.String()
			break
		}
	}
	partitions, e := disk.PartitionsWithContext(s.ctx, false)
	if e != nil {
		return nil, e
	}
	files := []map[string]any{}
	for _, p := range partitions {
		d, e := disk.UsageWithContext(s.ctx, p.Mountpoint)
		if e == nil {
			files = append(files, map[string]any{"dirName": p.Mountpoint, "sysTypeName": p.Fstype, "typeName": p.Device, "total": bytesText(d.Total), "free": bytesText(d.Free), "used": bytesText(d.Used), "usage": d.UsedPercent})
		}
	}
	const gb = 1 << 30
	const mb = 1 << 20
	return result(map[string]any{
		"cpu": map[string]any{"cpuNum": runtime.NumCPU(), "used": 100 * c.User / total, "sys": 100 * c.System / total, "free": 100 * c.Idle / total},
		"mem": map[string]any{"total": float64(vm.Total) / gb, "used": float64(vm.Used) / gb, "free": float64(vm.Available) / gb, "usage": vm.UsedPercent},
		"jvm": map[string]any{"name": "Go runtime", "version": runtime.Version(), "total": m.Sys / mb, "used": m.Alloc / mb, "free": (m.Sys - m.Alloc) / mb, "usage": float64(m.Alloc) * 100 / float64(m.Sys), "startTime": startedAt.Format(time.DateTime), "runTime": time.Since(startedAt).Round(time.Second).String(), "home": runtime.GOROOT()},
		"sys": map[string]any{"computerName": host, "computerIp": ip, "osName": runtime.GOOS, "osArch": runtime.GOARCH, "userDir": cwd}, "sysFiles": files}, nil)
}

func bytesText(n uint64) string { return fmt.Sprintf("%.2f GB", float64(n)/(1<<30)) }
