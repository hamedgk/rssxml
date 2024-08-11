package rssxml

import (
	"reflect"
	"testing"
)

var queryTable = []struct {
	input    string
	expected QueryTokens
}{
	{
		input:    "{item}",
		expected: QueryTokens{nil, []Token{{"item", nil}}},
	},
	{
		input:    "item{item}",
		expected: QueryTokens{[]Token{{"item", nil}}, []Token{{"item", nil}}},
	},
	{
		input:    "item{title,description}",
		expected: QueryTokens{[]Token{{"item", nil}}, []Token{{"title", nil}, {"description", nil}}},
	},
	{
		input:    "item{title,media:content[url],description}",
		expected: QueryTokens{[]Token{{"item", nil}}, []Token{{"title", nil}, {"media:content", []string{"url"}}, {"description", nil}}},
	},
	{
		input:    "{title,media:content[url],description}",
		expected: QueryTokens{nil, []Token{{"title", nil}, {"media:content", []string{"url"}}, {"description", nil}}},
	},
	{
		input:    "{title,description,media:content[url]}",
		expected: QueryTokens{nil, []Token{{"title", nil}, {"description", nil}, {"media:content", []string{"url"}}}},
	},
	{
		input: "rss~item{title,description,media:content[url]}",
		expected: QueryTokens{
			[]Token{
				{"rss", nil},
				{"item", nil},
			},
			[]Token{
				{"title", nil}, {"description", nil}, {"media:content", []string{"url"}},
			},
		},
	},
}

func TestQueryParser(t *testing.T) {
	for _, c := range queryTable {
		expected := ParseRSSQuery(c.input)
		if !reflect.DeepEqual(expected, c.expected) {
			t.Fail()
			t.Log("failed parsing", c.input)
		}
	}
}
