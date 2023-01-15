package gee

import "testing"

func TestNestGroup(t *testing.T) {
	g := New()
	a := g.Group("/a")
	b := a.Group("/b")
	c := b.Group("/c")
	if b.prefix != "/a/b" {
		t.Fatal("b.prefix should be /a/b")
	}
	if c.prefix != "/a/b/c" {
		t.Fatal("c.prefix should be /a/b/c")
	}
}
