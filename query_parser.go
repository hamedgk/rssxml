package rssxml

type QueryTokens struct {
	Edges  []Token
	Leaves []Token
}

type Token struct {
	TagName    string
	Attributes []string
}

func ParseRSSQuery(q string) QueryTokens {
	var (
		tokens     = QueryTokens{}
		i          = 0
		edge, leaf = Token{}, Token{}
		tagName    = make([]byte, 0, 4)
	)
Loop:
	for {
		switch {
		case i == len(q):
			edge.TagName = string(tagName)
			tokens.addEdge(edge)
			return tokens
		case q[i] == '~':
			edge.TagName = string(tagName)
			tagName = make([]byte, 0, 4)
			tokens.addEdge(edge)
			edge = Token{}
			i++
			continue
		case q[i] == '{':
			if i > 0 {
				edge.TagName = string(tagName)
				tagName = make([]byte, 0, 4)
				tokens.addEdge(edge)
			}
			i++
			break Loop
		case q[i] == '[':
			i++
			i = parseTagAttributes(q, &edge, i)
		default:
			tagName = append(tagName, q[i])
			i++
		}
	}

	for {
		switch {
		case i == len(q):
			return tokens
		case q[i] == ',':
			leaf.TagName = string(tagName)
			tagName = make([]byte, 0, 4)
			tokens.addLeaf(leaf)
			leaf = Token{}
			i++
			continue
		case q[i] == '}':
			leaf.TagName = string(tagName)
			tokens.addLeaf(leaf)
			return tokens
		case q[i] == '[':
			i++
			i = parseTagAttributes(q, &leaf, i)
		default:
			tagName = append(tagName, q[i])
			i++
		}
	}
}

func (t *Token) addAttribute(attribute []byte) {
	t.Attributes = append(t.Attributes, string(attribute))
}

func (qt *QueryTokens) addEdge(t Token) {
	qt.Edges = append(qt.Edges, t)
}

func (qt *QueryTokens) addLeaf(t Token) {
	qt.Leaves = append(qt.Leaves, t)
}

func parseTagAttributes(q string, t *Token, i int) int {
	attribute := make([]byte, 0, 4)
	for {
		switch {
		case q[i] == ']':
			t.addAttribute(attribute)
			i++
			return i
		case q[i] == ',':
			t.addAttribute(attribute)
			attribute = make([]byte, 0, 4)
			i++
			continue
		default:
			attribute = append(attribute, q[i])
			i++
		}
	}
}
