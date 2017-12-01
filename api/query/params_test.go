// Copyright 2017 Intuit, Inc.  All rights reserved.
// Use of this source code is governed the MIT license
// that can be found in the LICENSE file.

package query

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/intuit/ami-query/amicache"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  Params
	}{
		{
			amicache.StateTag,
			fmt.Sprintf("%s=available&%[1]s=deprecated&%[1]s=available", amicache.StateTag),
			Params{
				regions: []string{},
				images:  []string{},
				tags: map[string][]string{
					amicache.StateTag: []string{"available", "deprecated"},
				},
			},
		},
		{
			"tags",
			"tag=foo1:bar&tag=foo2:bar&tag=foo1:baz&tag=foo1:baz&tag=foo2:baz",
			Params{
				regions: []string{},
				images:  []string{},
				tags: map[string][]string{
					"foo1": []string{"bar", "baz"},
					"foo2": []string{"bar", "baz"},
				},
			},
		},
		{
			"ami",
			"ami=ami-1a2b3c4d&ami=ami-2a2b3c4d&ami=ami-3a2b3c4d&ami=ami-2a2b3c4d",
			Params{
				regions: []string{},
				images:  []string{"ami-1a2b3c4d", "ami-2a2b3c4d", "ami-3a2b3c4d"},
				tags:    map[string][]string{},
			},
		},
		{
			"account_id",
			"account_id=foo&account_id=bar&account_id=foo",
			Params{
				acctID:  "foo",
				regions: []string{},
				images:  []string{},
				tags:    map[string][]string{},
			},
		},
		{
			"callback",
			"callback=foo&callback=bar&callback=foo",
			Params{
				callback: "foo",
				regions:  []string{},
				images:   []string{},
				tags:     map[string][]string{},
			},
		},
		{
			"pretty",
			"pretty",
			Params{
				pretty:  true,
				regions: []string{},
				images:  []string{},
				tags:    map[string][]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{}
			if err := p.Decode(&url.URL{RawQuery: tt.query}); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.want, p) {
				t.Errorf("\n\twant: %#v\n\t got: %#v", tt.want, p)
			}
		})
	}
}

func TestDecodeBadKey(t *testing.T) {
	p := &Params{}
	err := p.Decode(&url.URL{RawQuery: "foo=bar"})
	if want, got := "unknown query key: foo", err.Error(); want != got {
		t.Errorf("\n\twant err: %q\n\t got err: %q", want, got)
	}
}

func TestDecodeBadTagValue(t *testing.T) {
	p := &Params{}
	err := p.Decode(&url.URL{RawQuery: "tag=foo:bar:baz"})
	if want, got := "invalid query tag value: foo:bar:baz", err.Error(); want != got {
		t.Errorf("\n\twant err: %q\n\t got err: %q", want, got)
	}
}

func TestDecodeParseError(t *testing.T) {
	p := &Params{}
	err := p.Decode(&url.URL{RawQuery: `foo=%%bar`})
	if want, got := `invalid URL escape "%%b"`, err.Error(); want != got {
		t.Errorf("\n\twant err: %q\n\t got err: %q", want, got)
	}
}

func TestDedup(t *testing.T) {
	got := dedup([]string{"foo", "bar", "baz", "foo"})
	want := []string{"foo", "bar", "baz"}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
