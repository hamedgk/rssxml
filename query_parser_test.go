package rssxml

import (
	"reflect"
	"testing"
)

var table = []struct {
	got  string
	want QueryTokens
}{
	{
		got:  "{item}",
		want: QueryTokens{nil, []Token{{"item", nil}}},
	},
	{
		got:  "item{item}",
		want: QueryTokens{[]Token{{"item", nil}}, []Token{{"item", nil}}},
	},
	{
		got:  "item{title,description}",
		want: QueryTokens{[]Token{{"item", nil}}, []Token{{"title", nil}, {"description", nil}}},
	},
	{
		got:  "item{title,media:content[url],description}",
		want: QueryTokens{[]Token{{"item", nil}}, []Token{{"title", nil}, {"media:content", []string{"url"}}, {"description", nil}}},
	},
	{
		got:  "{title,media:content[url],description}",
		want: QueryTokens{nil, []Token{{"title", nil}, {"media:content", []string{"url"}}, {"description", nil}}},
	},
}

func TestQuery(t *testing.T) {
	for _, c := range table {
		got := ParseRSSQuery(c.got)
		if !reflect.DeepEqual(got, c.want) {
			t.Fail()
			t.Log("failed parsing", c.got)
		}
	}
}
