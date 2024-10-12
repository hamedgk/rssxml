package rssxml

type TagData struct {
	TagName, Content string
	Attributes       map[string]string
}

func ObjectAttributes(text string, item Token, idx int) (map[string]string, int) {
	idx = iterateOpeningTag(text, item.TagName, idx, true)
	if idx == 0 {
		return nil, idx
	}
	m := map[string]string{}
	for _, attribute := range item.Attributes {
		startContentIdx, endContentIdx := iterateAttribute(text, attribute, idx)
		if startContentIdx == 0 && endContentIdx == 0 {
			return nil, idx
		}
		m[attribute] = text[startContentIdx:endContentIdx]
		idx = endContentIdx
	}
	idx = skipUntil(text, '>', idx)
	return m, idx
}

func ObjectTag(text string, item Token, idx int) (TagData, int) {
	startContentIdx := iterateOpeningTag(text, item.TagName, idx, false)
	if startContentIdx == 0 {
		return TagData{}, 0
	}
	endContentIdx := iterateClosingTag(text, item.TagName, startContentIdx)
	if endContentIdx == 0 {
		return TagData{}, 0
	}
	return TagData{
		TagName: item.TagName,
		Content: text[startContentIdx:endContentIdx],
	}, endContentIdx
}

func ObjectTagWithAttributes(text string, item Token, idx int) (TagData, int) {
	attributes, startContentIdx := ObjectAttributes(text, item, idx)
	if attributes == nil {
		return TagData{}, 0
	}
	endContentIdx := iterateClosingTag(text, item.TagName, startContentIdx)
	if endContentIdx == 0 {
		return TagData{}, 0
	}
	return TagData{
		TagName:    item.TagName,
		Content:    text[startContentIdx:endContentIdx],
		Attributes: attributes,
	}, endContentIdx
}

func Extract(text string, tokens QueryTokens, idx int) (map[string]TagData, int) {
	data := map[string]TagData{}
	regression := idx
	for _, edge := range tokens.Edges {
		idx = iterateOpeningTag(text, edge.TagName, regression, false)
		if idx == 0 {
			continue
		}
		regression = idx
	}
	for _, leaf := range tokens.Leaves {
		var tagData TagData
		if leaf.Attributes == nil {
			tagData, idx = ObjectTag(text, leaf, regression)
		} else {
			tagData, idx = ObjectTagWithAttributes(text, leaf, regression)
		}
		//returning zero indicates failure
		if idx == 0 {
			continue
		}
		regression = idx
		data[tagData.TagName] = tagData
	}
	return data, regression
}

func iterateOpeningTag(text, tag string, idx int, onlyTagName bool) int {
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
			return idx
		}
		idx = skipUntil(text, '>', idx)
		return idx
	}
	return 0
}

func iterateClosingTag(text, tag string, idx int) int {
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
		return beginTagIdx
	}
	return 0
}

func iterateAttribute(text, attribute string, idx int) (int, int) {
	idx = skipWhile(text, ' ', idx)
	i := 0
	for i < len(attribute) {
		if text[idx] != attribute[i] {
			for text[idx] != ' ' {
				idx++
				if text[idx] == '>' {
					return 0, 0
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
	return beginAttributeIdx, idx - 1
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
