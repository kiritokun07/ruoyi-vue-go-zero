package system

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/generator"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func (s *Session) renderTable(id int64) (types.GenTable, []generator.Generated, error) {
	t, e := s.genTable(id)
	if e != nil {
		return t, nil, e
	}
	var child *types.GenTable
	if t.TplCategory == "sub" {
		childID, e := one[int64](s, "SELECT table_id FROM gen_table WHERE table_name=?", t.SubTableName)
		if e != nil {
			return t, nil, e
		}
		v, e := s.genTable(childID)
		if e != nil {
			return t, nil, e
		}
		child = &v
	}
	files, e := generator.Render(t, child)
	return t, files, e
}

func (s *Session) GenPreview(r *types.RecordId) (*types.Result, error) {
	_, files, e := s.renderTable(r.Id)
	if e != nil {
		return nil, e
	}
	preview := map[string]string{}
	for _, f := range files {
		preview[f.Template] = f.Text
	}
	return result(preview, nil)
}

func (s *Session) generateByName(name string) (types.GenTable, []generator.Generated, error) {
	id, e := one[int64](s, "SELECT table_id FROM gen_table WHERE table_name=?", name)
	if e != nil {
		return types.GenTable{}, nil, e
	}
	return s.renderTable(id)
}

func (s *Session) GenDownload(r *types.LookupKey) (*File, error) {
	return s.GenBatch(&types.TablesForm{Tables: r.Key})
}

func (s *Session) GenBatch(r *types.TablesForm) (*File, error) {
	names, e := tableNames(r.Tables)
	if e != nil {
		return nil, e
	}
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	used := map[string]bool{}
	for _, name := range names {
		_, files, e := s.generateByName(name)
		if e != nil {
			return nil, e
		}
		for _, f := range files {
			if used[f.Path] {
				return nil, apierror.Business("生成文件重名，请调整生成配置")
			}
			used[f.Path] = true
			file, e := w.Create(f.Path)
			if e != nil {
				return nil, e
			}
			if _, e = file.Write([]byte(f.Text)); e != nil {
				return nil, e
			}
		}
	}
	if e = w.Close(); e != nil {
		return nil, e
	}
	return &File{Name: "ruoyi.zip", Type: "application/zip", Bytes: buf.Bytes()}, nil
}

func (s *Session) GenGenerate(r *types.LookupKey) (*types.StatusResponse, error) {
	t, files, e := s.generateByName(r.Key)
	if e != nil {
		return nil, e
	}
	base, e := filepath.Abs(s.services.Config.GeneratorDir)
	if e != nil {
		return nil, e
	}
	dest := t.GenPath
	if dest == "" || dest == "/" {
		dest = "."
	}
	if filepath.IsAbs(dest) {
		dest, e = filepath.Rel(base, dest)
		if e != nil {
			return nil, e
		}
	}
	if dest != "." && !filepath.IsLocal(dest) {
		return nil, apierror.Business("生成路径必须位于配置的生成目录内")
	}
	if e = os.MkdirAll(base, 0700); e != nil {
		return nil, e
	}
	root, e := os.OpenRoot(base)
	if e != nil {
		return nil, e
	}
	defer root.Close()
	for _, f := range files {
		if !filepath.IsLocal(f.Path) {
			return nil, apierror.Business("生成文件路径无效")
		}
		name := filepath.Join(dest, filepath.FromSlash(f.Path))
		if e = root.MkdirAll(filepath.Dir(name), 0755); e != nil {
			return nil, e
		}
		tmp := name + "." + uuid.NewString()
		if e = root.WriteFile(tmp, []byte(f.Text), 0644); e != nil {
			return nil, e
		}
		if e = root.Rename(tmp, name); e != nil {
			root.Remove(tmp)
			return nil, e
		}
	}
	return &types.StatusResponse{Code: 200, Msg: fmt.Sprintf("已生成 %d 个文件至 %s", len(files), strings.TrimSuffix(filepath.Join(base, dest), "/"))}, nil
}
