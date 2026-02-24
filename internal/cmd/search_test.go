package cmd

import "testing"

func TestNormalizeSearchItems(t *testing.T) {
	resp := searchAPIResponse{
		Items: []searchAPIItem{
			{Slug: "@acme/a", Summary: "first", Owner: "acme"},
		},
		Packages: []searchAPIPackage{
			{Name: "@scope/pkg", Version: "1.2.3"},
		},
	}
	items := normalizeSearchItems(resp)
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[1].Slug != "@scope/pkg" || items[1].Owner != "scope" || items[1].Summary != "version 1.2.3" {
		t.Fatalf("unexpected normalized package item: %+v", items[1])
	}
}

func TestFilterSearchItemsOwner(t *testing.T) {
	items := []searchAPIItem{
		{Slug: "@acme/a", Owner: "acme"},
		{Slug: "@foo/b", Owner: "foo"},
		{Slug: "@Acme/c", Owner: "Acme"},
	}
	filtered := filterSearchItems(items, "@acme")
	if len(filtered) != 2 {
		t.Fatalf("expected 2 items for owner acme, got %d", len(filtered))
	}
}

func TestSortSearchItems(t *testing.T) {
	items := []searchAPIItem{
		{Slug: "@zeta/b", Owner: "zeta"},
		{Slug: "@acme/c", Owner: "acme"},
		{Slug: "@acme/a", Owner: "acme"},
	}

	bySlug := sortSearchItems(items, "slug")
	if bySlug[0].Slug != "@acme/a" || bySlug[2].Slug != "@zeta/b" {
		t.Fatalf("unexpected slug sort result: %+v", bySlug)
	}

	byOwner := sortSearchItems(items, "owner")
	if byOwner[0].Owner != "acme" || byOwner[0].Slug != "@acme/a" {
		t.Fatalf("unexpected owner sort result: %+v", byOwner)
	}
}

func TestLimitSearchItems(t *testing.T) {
	items := []searchAPIItem{
		{Slug: "a"},
		{Slug: "b"},
		{Slug: "c"},
	}
	limited := limitSearchItems(items, 2)
	if len(limited) != 2 || limited[1].Slug != "b" {
		t.Fatalf("unexpected limited result: %+v", limited)
	}

	unlimited := limitSearchItems(items, 0)
	if len(unlimited) != 3 {
		t.Fatalf("expected no limit when limit=0")
	}
}

func TestOwnerFromSlug(t *testing.T) {
	if got := ownerFromSlug("@scope/name"); got != "scope" {
		t.Fatalf("expected scope owner, got %q", got)
	}
	if got := ownerFromSlug("plain"); got != "" {
		t.Fatalf("expected empty owner for plain slug, got %q", got)
	}
}
