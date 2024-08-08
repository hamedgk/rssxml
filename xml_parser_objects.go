package rssxml

type AttributesData = map[string]string

type TagData struct {
	TagName, Content string
	Attributes       AttributesData
}

func ObjectAttributes(text string, item Token, idx int) (AttributesData, int) {
	success, idx := iterateOpeningTag(text, item.TagName, idx, true)
	if !success {
		return nil, idx
	}
	m := make(map[string]string)
	for _, attribute := range item.Attributes {
		success, startContentIdx, endContentIdx := iterateAttribute(text, attribute, idx)
		if !success {
			return nil, idx
		}
		m[attribute] = text[startContentIdx:endContentIdx]
		idx = endContentIdx
	}
	idx = skipUntil(text, '>', idx)
	return m, idx
}

func ObjectTag(text string, item Token, idx int) (*TagData, int) {
	success, startContentIdx := iterateOpeningTag(text, item.TagName, idx, false)
	if !success {
		return nil, startContentIdx
	}
	success, endContentIdx := iterateClosingTag(text, item.TagName, startContentIdx)
	if !success {
		return nil, endContentIdx
	}
	return &TagData{
		TagName: item.TagName,
		Content: text[startContentIdx:endContentIdx],
	}, endContentIdx
}

func ObjectTagWithAttributes(text string, item Token, idx int) (*TagData, int) {
	attributes, startContentIdx := ObjectAttributes(text, item, idx)
	if attributes == nil {
		return nil, startContentIdx
	}
	success, endContentIdx := iterateClosingTag(text, item.TagName, startContentIdx)
	if !success {
		return nil, endContentIdx
	}
	return &TagData{
		TagName:    item.TagName,
		Content:    text[startContentIdx:endContentIdx],
		Attributes: attributes,
	}, endContentIdx
}

func Extract(text string, tokens QueryTokens, idx int) ([]*TagData, int) {
	var data []*TagData
	var regression int = idx
	for _, edge := range tokens.Edges {
		var success bool
		success, idx = iterateOpeningTag(text, edge.TagName, regression, false)
		if !success {
			continue
		}
		regression = idx
	}
	for _, leaf := range tokens.Leaves {
		var tagData *TagData
		if leaf.Attributes == nil {
			tagData, idx = ObjectTag(text, leaf, regression)
		} else {
			tagData, idx = ObjectTagWithAttributes(text, leaf, regression)
		}
		if tagData == nil {
			continue
		}
		regression = idx
		data = append(data, tagData)
	}
	return data, regression
}
