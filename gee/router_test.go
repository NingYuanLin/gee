package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	pattern := "/p/:name"
	parts := parsePattern(pattern)
	realParts := []string{"p", ":name"}
	ok := reflect.DeepEqual(parts, realParts)
	if ok == false {
		t.Fatal(fmt.Sprintf("%q should be %q, rather than%q\n", pattern, realParts, parts))
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, params := r.getRoute("GET", "/hello/ning")

	if n == nil {
		t.Fatal("can't get route")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if n.part != ":name" {
		t.Fatal("part should be :name")
	}

	if params["name"] != "ning" {
		t.Fatal("name in params should be \"ning\"")
	}

}
