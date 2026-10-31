package auth

import (
	"testing"

	"ruoyi-vue-go-zero/app/admin/api/internal/model"
)

func TestBuildRoutes(t *testing.T) {
	menus := []model.Menu{
		{ID: 1, Type: "M", Internal: true, Path: "system"},
		{ID: 2, Parent: 1, Type: "M", Internal: true, Path: "nested"},
		{ID: 3, Parent: 2, Type: "C", Internal: true, Path: "https://www.example.com:8080", Hidden: true},
		{ID: 4, Type: "C", Internal: true, Path: "profile", Name: "custom", Component: "user/profile", Query: "{}"},
		{ID: 5, Parent: 100, Type: "C", Path: "orphan"},
	}
	routes := buildRoutes(menus)
	if len(routes) != 2 {
		t.Fatalf("unexpected roots: %+v", routes)
	}
	root := routes[0]
	if root.Path != "/system" || root.Redirect != "noRedirect" || !root.AlwaysShow || root.Component != "Layout" {
		t.Fatalf("directory: %+v", root)
	}
	nested := root.Children[0]
	link := nested.Children[0]
	if nested.Component != "ParentView" || link.Component != "InnerLink" || link.Path != "example/com/8080" || !link.Hidden || link.Meta.Link != menus[2].Path {
		t.Fatalf("nested link: %+v", link)
	}
	page := routes[1]
	if page.Path != "/" || page.Name != "" || page.Meta != nil || page.Children[0].Name != "Custom" || page.Children[0].Query != "{}" {
		t.Fatalf("root page: %+v", page)
	}
}
