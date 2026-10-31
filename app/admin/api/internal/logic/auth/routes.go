package auth

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"ruoyi-vue-go-zero/app/admin/api/internal/model"
	"ruoyi-vue-go-zero/app/admin/api/internal/types"
)

func buildRoutes(menus []model.Menu) []types.Route {
	byParent := make(map[int64][]model.Menu)
	for _, menu := range menus {
		if menu.ID > 0 {
			byParent[menu.Parent] = append(byParent[menu.Parent], menu)
		}
	}
	return childRoutes(byParent, 0)
}

func routeName(name, path string) string {
	if name == "" {
		name = path
	}
	if name == "" {
		return ""
	}
	first, n := utf8.DecodeRuneInString(name)
	return string(unicode.ToUpper(first)) + name[n:]
}

var linkPath = strings.NewReplacer("https://", "", "http://", "", "www.", "", ".", "/", ":", "/")

func childRoutes(byParent map[int64][]model.Menu, parent int64) []types.Route {
	routes := make([]types.Route, 0, len(byParent[parent]))
	for _, m := range byParent[parent] {
		meta := &types.RouteMeta{Title: m.Title, Icon: m.Icon, NoCache: m.NoCache}
		if strings.HasPrefix(m.Path, "http://") || strings.HasPrefix(m.Path, "https://") {
			meta.Link = m.Path
		}
		inner := m.Internal && meta.Link != ""
		rootPage := parent == 0 && m.Type == "C" && m.Internal
		r := types.Route{Name: routeName(m.Name, m.Path), Hidden: m.Hidden, Path: m.Path, Query: m.Query, Meta: meta, Component: m.Component}
		if r.Component == "" || rootPage {
			switch {
			case parent != 0 && inner:
				r.Component = "InnerLink"
			case parent != 0 && m.Type == "M":
				r.Component = "ParentView"
			default:
				r.Component = "Layout"
			}
		}
		if parent == 0 && m.Type == "M" && m.Internal {
			r.Path = "/" + m.Path
		} else if parent != 0 && inner {
			r.Path = linkPath.Replace(m.Path)
		}
		switch {
		case m.Type == "M" && len(byParent[m.ID]) > 0:
			r.AlwaysShow = true
			r.Redirect = "noRedirect"
			r.Children = childRoutes(byParent, m.ID)
		case rootPage:
			r.Children = []types.Route{{Name: r.Name, Path: m.Path, Component: m.Component, Query: m.Query, Meta: meta}}
			r.Name, r.Path, r.Meta = "", "/", nil
		case parent == 0 && inner:
			path := linkPath.Replace(m.Path)
			r.Path = "/"
			r.Meta = &types.RouteMeta{Title: m.Title, Icon: m.Icon}
			r.Children = []types.Route{{Name: routeName(m.Name, path), Path: path, Component: "InnerLink", Meta: &types.RouteMeta{Title: m.Title, Icon: m.Icon, Link: m.Path}}}
		}
		routes = append(routes, r)
	}
	return routes
}
