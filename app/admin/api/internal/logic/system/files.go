package system

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"ruoyi-vue-go-zero/app/admin/api/internal/apierror"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

type File struct {
	Name, Type string
	Bytes      []byte
	AfterSend  func()
}

func (f *File) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", f.Type)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if f.Name != "" {
		w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(f.Name))
	}
	if _, e := io.Copy(w, bytes.NewReader(f.Bytes)); e == nil && f.AfterSend != nil {
		f.AfterSend()
	}
}

var permittedExtension = regexp.MustCompile(`^\.(jpg|jpeg|png|gif|webp|bmp|pdf|docx?|xlsx?|pptx?|txt|csv|zip|rar|7z|mp[34]|wav)$`)

func (s *Session) Upload(files []*multipart.FileHeader, multiple bool, base string) (*types.UploadResult, error) {
	if len(files) == 0 {
		return nil, apierror.Business("请选择上传文件")
	}
	dir := s.services.Config.UploadDir
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	names, links, originals, urls := []string{}, []string{}, []string{}, []string{}
	success := false
	defer func() {
		if !success {
			for _, n := range names {
				os.Remove(filepath.Join(dir, n))
			}
		}
	}()
	for _, h := range files {
		ext := strings.ToLower(filepath.Ext(h.Filename))
		if !permittedExtension.MatchString(ext) || h.Size <= 0 || h.Size > 20<<20 {
			return nil, apierror.Business("文件类型或大小不允许")
		}
		name := uuid.NewString() + ext
		in, e := h.Open()
		if e != nil {
			return nil, e
		}
		f, e := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			in.Close()
			return nil, e
		}
		names = append(names, name)
		_, e = io.Copy(f, in)
		in.Close()
		closeErr := f.Close()
		if e != nil {
			return nil, e
		}
		if closeErr != nil {
			return nil, closeErr
		}
		link := "/profile/upload/" + name
		links = append(links, link)
		urls = append(urls, base+link)
		originals = append(originals, path.Base(h.Filename))
	}
	success = true
	res := &types.UploadResult{Code: 200, Msg: "操作成功"}
	if multiple {
		res.Urls = strings.Join(urls, ",")
		res.FileNames = strings.Join(links, ",")
		res.NewFileNames = strings.Join(names, ",")
		res.OriginalFilenames = strings.Join(originals, ",")
	} else {
		res.Url = urls[0]
		res.FileName = links[0]
		res.NewFileName = names[0]
		res.OriginalFilename = originals[0]
	}
	return res, nil
}

func readResource(dir, name string) (*File, error) {
	if name == "" || name != filepath.Base(name) || !permittedExtension.MatchString(strings.ToLower(filepath.Ext(name))) {
		return nil, apierror.Business("非法文件名")
	}
	root, e := os.OpenRoot(dir)
	if e != nil {
		return nil, apierror.Business("文件不存在")
	}
	defer root.Close()
	data, e := root.ReadFile(name)
	if os.IsNotExist(e) {
		return nil, apierror.Business("文件不存在")
	}
	if e != nil {
		return nil, e
	}
	ct := mime.TypeByExtension(filepath.Ext(name))
	if ct == "" {
		ct = "application/octet-stream"
	}
	return &File{Name: name, Type: ct, Bytes: data}, nil
}

func (s *Session) DownloadResource(r *types.DownloadForm) (*File, error) {
	if name, ok := strings.CutPrefix(r.Resource, "/profile/avatar/"); ok {
		return readResource(s.services.Config.AvatarDir, name)
	}
	if name, ok := strings.CutPrefix(r.Resource, "/profile/upload/"); ok {
		return readResource(s.services.Config.UploadDir, name)
	}
	return nil, apierror.Business("非法资源路径")
}

func (s *Session) DownloadFile(r *types.DownloadForm) (*File, error) {
	f, e := readResource(s.services.Config.DownloadDir, r.FileName)
	if e != nil {
		return nil, e
	}
	if r.Delete {
		f.AfterSend = func() { os.Remove(filepath.Join(s.services.Config.DownloadDir, r.FileName)) }
	}
	return f, nil
}

func ServeUpload(dir string, w http.ResponseWriter, name string) {
	f, e := readResource(dir, name)
	if e != nil {
		http.Error(w, "file not found", 404)
		return
	}
	if strings.HasPrefix(f.Type, "image/") {
		f.Name = ""
	}
	w.Header().Set("Content-Security-Policy", "default-src 'none'; sandbox")
	f.Send(w)
}

func welcome() *File {
	return &File{Type: "text/plain; charset=utf-8", Bytes: []byte("欢迎使用若依后台管理框架（Go），请通过前端地址访问。")}
}
