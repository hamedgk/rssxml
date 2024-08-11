package rssxml

import (
	"io"
	"os"
	"reflect"
	"testing"
)

var extractTable = []struct {
	query    string
	expected []*TagData
}{
	{
		query: "item{title}",
		expected: []*TagData{
			{
				TagName: "title",
				Content: "Hamas Elevates Gaza Leader Who Planned Oct. 7 Attacks to Top Post",
			},
		},
	},
}

func TestXMLParser(t *testing.T) {
	file, err := os.Open("test_data/nytimes.xml")
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	runTest(t, string(data))
}

func runTest(t *testing.T, text string) {
	for _, c := range extractTable {
		tokens := ParseRSSQuery(c.query)
		output, _ := Extract(text, tokens, 0)
		if !reflect.DeepEqual(c.expected, output) {
			t.Log(c.query)
			t.Fail()
		}
	}
}
