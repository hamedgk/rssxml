package rssxml

func iterateOpeningTag(text, tag string, idx int, onlyTagName bool) (bool, int) {
OuterLoop:
	for idx < len(text) {
		if text[idx] != '<' {
			idx++
			continue
		}
		idx++
		for i := range tag {
			if tag[i] != text[idx] {
				idx = skipUntil(text, '>', idx)
				continue OuterLoop
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
OuterLoop:
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
				continue OuterLoop
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
