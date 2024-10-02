package rssxml

type TagData struct {
	TagName, Content string
	Attributes       map[string]string
}

func ObjectAttributes(text string, item Token, idx int) (map[string]string, int) {
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

func iterateOpeningTag(text, tag string, idx int, onlyTagName bool) (bool, int) {
Outer:
	for idx < len(text) {
		if text[idx] != '<' {
			idx++
			continue
		}
		idx++
		for i := range tag {
			if tag[i] != text[idx] {
				idx = skipUntil(text, '>', idx)
				continue Outer
			}
			idx++
		}
		if onlyTagName {
			return true, idx
		}
		idx = skipUntil(text, '>', idx)
		return true, idx
	}
	return false, idx
}

func iterateClosingTag(text, tag string, idx int) (bool, int) {
	var beginTagIdx int
Outer:
	for idx < len(text) {
		if text[idx] != '<' || text[idx+1] != '/' {
			idx++
			continue
		}
		beginTagIdx = idx
		idx += 2
		for i := range tag {
			if tag[i] != text[idx] {
				idx = skipUntil(text, '>', idx)
				continue Outer
			}
			idx++
		}
		return true, beginTagIdx
	}
	return false, idx
}

func iterateAttribute(text, attribute string, idx int) (bool, int, int) {
	idx = skipWhile(text, ' ', idx)
	i := 0
	for i < len(attribute) {
		if text[idx] != attribute[i] {
			for text[idx] != ' ' {
				idx++
				if text[idx] == '>' {
					return false, idx + 1, 0
				}
			}
			idx = skipWhile(text, ' ', idx)
			i = 0
			continue
		}
		idx++
		i++
	}
	idx = skipUntil(text, '"', idx)
	beginAttributeIdx := idx
	idx = skipUntil(text, '"', idx)
	return true, beginAttributeIdx, idx - 1
}

func skipUntil(text string, c byte, idx int) int {
	for text[idx] != c {
		idx++
	}
	idx++
	return idx
}

func skipWhile(text string, c byte, idx int) int {
	for text[idx] == c {
		idx++
	}
	return idx
}
