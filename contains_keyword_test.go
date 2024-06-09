package main_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/sue445/feed_proxy"
	"testing"
)

func TestContainsKeyword(t *testing.T) {
	type args struct {
		source  string
		keyword string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "not contains",
			args: args{
				source:  "XXX",
				keyword: "AAA",
			},
			want: false,
		},
		{
			name: "single keyword matches",
			args: args{
				source:  "XXXAAAABBB",
				keyword: "AAA",
			},
			want: true,
		},
		{
			name: "single keyword matches (full-width character)",
			args: args{
				source:  "ＸＸＸＡＡＡＸＸＸ",
				keyword: "ＡＡＡ",
			},
			want: true,
		},
		{
			name: "AAA BBB",
			args: args{
				source:  "AAABBB",
				keyword: "AAA BBB",
			},
			want: true,
		},
		{
			name: "(AAA BBB) | CCC",
			args: args{
				source:  "AAA BBB",
				keyword: "(AAA BBB) | CCC",
			},
			want: true,
		},
		{
			name: "(AAA | BBB) CCC",
			args: args{
				source:  "AAACCC",
				keyword: "(AAA | BBB) CCC",
			},
			want: true,
		},
		{
			name: "(A OR B) AND B",
			args: args{
				source:  "sue445",
				keyword: "(sue aaa) | 445",
			},
			want: true,
		},
		{
			name: "empty keyword",
			args: args{
				source:  "XXX",
				keyword: "",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := main.ContainsKeyword(tt.args.source, tt.args.keyword)

			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
