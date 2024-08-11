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
	{
		query: "item{title,guid,description,dc:creator}",
		expected: []*TagData{
			{
				TagName: "title",
				Content: "Hamas Elevates Gaza Leader Who Planned Oct. 7 Attacks to Top Post",
			},
			{
				TagName: "guid",
				Content: "https://www.nytimes.com/2024/08/06/world/middleeast/hamas-yahya-sinwar-promoted.html",
			},
			{
				TagName: "description",
				Content: "The selection of Yahya Sinwar, a prime target of Israeli forces, to replace the assassinated Ismail Haniyeh consolidates authority in the hands of a hard-liner who is in hiding.",
			},
			{
				TagName: "dc:creator",
				Content: "Adam Rasgon, Aaron Boxerman, Euan Ward and Michael Levenson",
			},
		},
	},
	{
		query: "item{title,guid,description,category[domain],media:content[url]}",
		expected: []*TagData{
			{
				TagName: "title",
				Content: "Hamas Elevates Gaza Leader Who Planned Oct. 7 Attacks to Top Post",
			},
			{
				TagName: "guid",
				Content: "https://www.nytimes.com/2024/08/06/world/middleeast/hamas-yahya-sinwar-promoted.html",
			},
			{
				TagName: "description",
				Content: "The selection of Yahya Sinwar, a prime target of Israeli forces, to replace the assassinated Ismail Haniyeh consolidates authority in the hands of a hard-liner who is in hiding.",
			},
			{
				TagName:    "category",
				Content:    "Israel-Gaza War (2023- )",
				Attributes: map[string]string{"domain": "http://www.nytimes.com/namespaces/keywords/des"},
			},
			{
				TagName:    "media:content",
				Attributes: map[string]string{"url": "https://static01.nyt.com/images/2024/08/06/multimedia/06mideast-crisis-sinwar-photo-plwh/06mideast-crisis-sinwar-photo-plwh-mediumSquareAt3X-v2.jpg"},
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
			t.Log("failed extracting", c.query)
			t.Fail()
		}
	}
}
