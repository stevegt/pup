package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	. "github.com/stevegt/goadapt"
	"github.com/stevegt/grokker"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func main() {
	fn := os.Args[1]

	inbuf, err := ioutil.ReadFile(fn)
	Ck(err)

	doc := parse(inbuf)
	concepts := doc.getSection("Concepts")
	concepts.ckForwards()
	concepts.ckUnrefs()
	doc.ckCitations()
	doc.showTodos()
	// doc.showTerms(fn)
}

// parse simply parses the input buffer and returns a Node tree
func parse(inbuf []byte) *Node {
	// create a text.Reader from the input file
	// - this is required by the parser
	txt := text.NewReader(inbuf)
	// create a parser
	p := goldmark.DefaultParser()
	// parse the input file
	doc := p.Parse(txt)
	// walk the tree and build a Node tree
	top := &Node{}
	root := top.fromAst(doc, inbuf)
	return root
}

// Node is a parsed node
//   - we use this as a wrapper around ast.Node because the goldmark
//     API and docs are byzantine and this lets us keep all the
//     crazy in one place
type Node struct {
	astNode ast.Node
	// parent is the parent node
	parent *Node
	// children is a slice of child nodes
	children []*Node
	// source is the source text for the entire document because
	// goldmark doesn't keep track of the source text for
	// a node itself so we have to pass it in every time
	// a goldmark ast.Node method needs it
	source []byte
}

// fromAst walks an ast.Node tree and builds a Node tree
func (parent *Node) fromAst(astNode ast.Node, source []byte) *Node {
	// parent is empty for the root node
	if parent.astNode == nil {
		parent = nil
	}
	// create a Node for the current ast.Node
	node := &Node{astNode: astNode, parent: parent, source: source}
	// get the children of the current ast.Node
	astChild := astNode.FirstChild()
	for astChild != nil {
		// walk the child and add it to the current node
		node.children = append(node.children, node.fromAst(astChild, source))
		// get the next child
		astChild = astChild.NextSibling()
	}
	return node
}

// dump prints the Node
func (node *Node) dump() {
	node.astNode.Dump(node.source, 0)
}

// walk is the walk function we wish goldmark had
// - it lets us carry source and other context in the Node struct
func (node *Node) walk(walker func(node *Node) (stop bool)) (stop bool) {
	// Pf("node: kind=%s, text=%s\n", node.astNode.Kind(), node.String())
	stop = walker(node)
	if stop {
		return
	}
	for _, child := range node.children {
		stop = child.walk(walker)
		if stop {
			return
		}
	}
	return
}

// key is the map key for the node
func (node *Node) key() string {
	// heading key is the heading text lowercased and with spaces
	// replaced by dashes
	if node.isHeading() {
		return strings.ReplaceAll(strings.ToLower(node.String()), " ", "-")
	}
	// link key is the link destination without the leading hash
	if node.isLink() {
		return strings.TrimPrefix(node.linkDst(), "#")
	}
	Assert(false, "node is not a heading or link")
	return ""
}

// isHeading returns true if the node is a heading
func (node *Node) isHeading() bool {
	_, ok := node.astNode.(*ast.Heading)
	return ok
}

// isLink returns true if the node is a Link
func (node *Node) isLink() bool {
	_, ok := node.astNode.(*ast.Link)
	return ok
}

// isRef returns true if the node is a reference to a citation
// - a citation is a link with the format [id][id]
func (node *Node) isRef() bool {
	if !node.isLink() {
		return false
	}
	// get the link display text -- this is the
	// text between the first set of square brackets and
	// is found in the first child node
	display := node.children[0].String()
	// it's a reference if the format is [text][id]
	// get the raw text of the node -- this may include
	// text from the parent node
	raw := node.raw()
	// look for the raw text of the link including both sets
	// of square brackets
	retxt := Spf(`\[%s\]\[(\S+)\]`, display)
	re := regexp.MustCompile(retxt)
	// if we find a match, it's a citation reference
	matches := re.FindStringSubmatch(raw)
	if len(matches) > 0 {
		return true
	}
	return false
}

// text returns the text string for the node
func (node *Node) String() string {
	return string(node.astNode.Text(node.source))
}

// linkDst returns the destination of a link node
func (node *Node) linkDst() string {
	Assert(node.isLink(), "node is not a link")
	dest := node.astNode.(*ast.Link).Destination
	return string(dest)
}

// isBlock returns true if the node is a block
func (node *Node) isBlock() bool {
	return node.astNode.Type() == ast.TypeBlock
}

// pos returns the position of the node in the source text
func (node *Node) pos() int {
	// if the node is not a block, we can't get its position, so
	// we return the position of its parent
	if !node.isBlock() {
		if node.parent == nil {
			return 0
		}
		return node.parent.pos()
	}
	segments := node.astNode.Lines()
	if segments.Len() == 0 {
		// Pf("child of node at pos %d has no lines\n", node.parent.pos())
		// node.dump()
		// Assert(false, "node has no lines")
		if node.parent == nil {
			return 0
		}
		return node.parent.pos()
	}
	segment := segments.At(0)
	start := segment.Start
	return start
}

// end returns the end position of the node in the source text
func (node *Node) end() int {
	// if the node is not a block, we can't get its position, so
	// we return the position of its parent
	if !node.isBlock() {
		if node.parent == nil {
			return 0
		}
		return node.parent.end()
	}
	segments := node.astNode.Lines()
	if segments.Len() == 0 {
		if node.parent == nil {
			return 0
		}
		return node.parent.end()
	}
	segment := segments.At(segments.Len() - 1)
	end := segment.Stop
	return end
}

// raw returns the raw text of the node
func (node *Node) raw() (txt string) {
	start := node.pos()
	end := node.end()
	if start == 0 || end <= start {
		return ""
	}
	txt = string(node.source[start:end])
	return
}

/*
// inConcepts returns true if the node is in the Concepts section.
// if node is not a heading, we check its parent.  if node is itself
// the Concepts heading, we return false.
func (node *Node) inConcepts() bool {
	parent := node.parent
	if parent == nil {
		return false
	}
	if !parent.isHeading() {
		return parent.inConcepts()
	}
	parentTxt := parent.String()
	if parentTxt == "Concepts" {
		return true
	}
	return false
}
*/

// headingLevel returns the heading level of the node, or -1 if it's not a heading
func (node *Node) headingLevel() int {
	if !node.isHeading() {
		return -1
	}
	return node.astNode.(*ast.Heading).Level
}

// getSection returns the named section of the document as a list of nodes
func (doc *Node) getSection(name string) (section *Section) {
	section = &Section{}

	// find the section start
	var start *Node
	doc.walk(func(node *Node) (stop bool) {
		if node.isHeading() && node.String() == name {
			start = node
			stop = true
		}
		return
	})
	Assert(start != nil, "section not found")
	startPos := start.pos()
	topLevel := start.headingLevel()
	Assert(topLevel >= 0, "section start is not a heading")

	// find the start of the next section
	var next *Node
	doc.walk(func(node *Node) (stop bool) {
		if node.pos() > startPos {
			level := node.headingLevel()
			if level >= 0 && level <= topLevel {
				next = node
				stop = true
				return
			}
		}
		return
	})

	// collect nodes between start and next
	doc.walk(func(node *Node) (stop bool) {
		if node.pos() >= startPos {
			if next != nil && node.pos() >= next.pos() {
				stop = true
				return
			}
			section.Add(node)
		}
		return
	})

	return

	/*

		// walk the tree looking for the start and end of the section,
		// collecting nodes in between
		topLevel := -1
		doc.walk(func(node *Node) (stop bool) {
			level := node.headingLevel()
			// look for the start of the section
			txt := node.String()
			if txt == name {
				topLevel = level
				section.Add(node)
				return
			}
			// keep going if we're not yet in the section
			if topLevel < 0 {
				return
			}
			// stop at the end of the section
			if level >= 0 && level <= topLevel {
				stop = true
				return
			}
			// add the node to the section
			section.Add(node)
			return
		})
		return
	*/
}

// Section is a section of the document, a list of Nodes
type Section struct {
	// nodes is the list of nodes in the section
	nodes []*Node
}

// key is the name of the section, lowercased and with spaces replaced by dashes
func (section *Section) key() string {
	return strings.ToLower(strings.Replace(section.name(), " ", "-", -1))
}

// name is the text of the first heading in the section
func (section *Section) name() string {
	for _, node := range section.nodes {
		if node.isHeading() {
			return node.String()
		}
	}
	return ""
}

// Add adds a node to the section
func (section *Section) Add(node *Node) {
	section.nodes = append(section.nodes, node)
}

// mapHeadings returns a map of headings in the given section
func (section *Section) mapHeadings() (headings map[string]*Node) {
	headings = make(map[string]*Node)
	// look for headings
	for _, node := range section.nodes {
		if node.isHeading() {
			headings[node.key()] = node
		}
	}
	return
}

// mapReferences returns a map of references in the given section
func (section *Section) mapReferences() (refs map[string]*Node) {
	refs = make(map[string]*Node)
	// walk each node looking for references
	for _, node := range section.nodes {
		node.walk(func(node *Node) (stop bool) {
			if node.isLink() {
				// Pf("node: %s kind %v\n", node.String(), node.astNode.Kind())
				refs[node.key()] = node
			}
			return
		})
	}
	return
}

// ckForwards checks the section for forward references and references
// that have no heading or citation
func (section *Section) ckForwards() {
	heads := section.mapHeadings()
	refs := section.mapReferences()
	// Pprint(refs)
	for key, ref := range refs {
		head, ok := heads[key]
		if !ok {
			if ref.isRef() {
				// Pl("XXX check citation:", key)
			} else {
				Pl("reference has no heading:", key)
			}
			/*
				Pl("string:", ref.String())
				Pl("pos:", ref.pos())
				Pl("children:")
				for i, child := range ref.children {
					Pf("  %d: %s\n", i, child.String())
				}
				Pl("len(children):", len(ref.children))
				Pl("raw:", ref.raw())
				Pl("linkDst:", ref.linkDst())
				Pl("isLink:", ref.isLink())
				Pl("isRef:", ref.isRef())
				Pf("attrs: %#v\n", ref.astNode.Attributes())
			*/
		} else {
			// reference has a heading
			// check for forward reference
			refpos := ref.pos()
			headpos := head.pos()
			if refpos < headpos {
				Pl("forward reference:", key)
			}
		}
	}
}

// ckUnrefs checks for unreferenced headings in the given section
func (section *Section) ckUnrefs() {
	heads := section.mapHeadings()
	refs := section.mapReferences()
	// sort heads by position
	sortedHeads := make([]*Node, 0, len(heads))
	for _, head := range heads {
		sortedHeads = append(sortedHeads, head)
	}
	sort.Slice(sortedHeads, func(i, j int) bool {
		return sortedHeads[i].pos() < sortedHeads[j].pos()
	})
	for _, node := range sortedHeads {
		key := node.key()
		// Pl(section.String(), "has heading:", key)
		if key == section.key() {
			// skip the section heading
			continue
		}
		_, ok := refs[key]
		if !ok {
			// heading has no reference
			Pl("heading has no reference:", key)
		}
	}
}

// ckCitations checks for citations that have no entry in the
// References section.
func (doc *Node) ckCitations() {
	// get the start of the References section
	refSection := doc.getSection("References")

	// okay, we aren't using goldmark for this because goldmark is too
	// smart by half when parsing links; in the case of GFM reference
	// links, it conflates the link text with the link destination.
	// That all fine, but if the link text is a citation, and the
	// citation is missing, then goldmark ignores it.  So we grunt
	// through this using the raw text of the document and some
	// regular expressions instead.

	// get the raw text of the References section
	// - we make the possibly incorrect assumption that the
	//   References section is the last section in the document
	refStart := refSection.nodes[0].pos()
	rawRefs := doc.source[refStart:]
	// Pl("rawRefs:", string(rawRefs))

	// find all citations in the document
	// - a citation is a link with the format [id][id]
	re := regexp.MustCompile(`\[([\s\w]+)\]\s*\[([\s\w]+)\]`)
	citations := re.FindAllStringSubmatch(string(doc.source), -1)
	// Pl("citations:", citations)

	// check each citation
	for _, citation := range citations {
		// Pl("citation:", citation)
		// make sure both ids are the same
		/*
			if citation[1] != citation[2] {
				Pl("citation ids don't match:", citation[1], citation[2])
			}
		*/
		id := citation[2]
		// check if the citation is in the References section
		resrc := `\[` + id + `\]:\s+`
		re := regexp.MustCompile(resrc)
		if !re.Match(rawRefs) {
			Pl("reference has no citation:", citation[0])
		}
	}

}

// showTodos shows all TODO and XXX marks in the document, with line numbers
func (doc *Node) showTodos() {
	// iterate over the lines in the document
	// - we do this using a line scanner instead of the AST because
	//   the AST doesn't preserve line numbers
	lines := strings.Split(string(doc.source), "\n")
	for i, line := range lines {
		hit := false
		// check for TODO
		if strings.Contains(line, "TODO") {
			hit = true
		}
		// check for XXX
		if strings.Contains(line, "XXX") {
			hit = true
		}
		if hit {
			Pf("%5d: %s\n", i+1, line)
		}
	}
	return
}

// showTodos shows all TODO and XXX marks in the document
func (doc *Node) XXXshowTodos() {
	// walk the doc
	doc.walk(func(node *Node) (stop bool) {
		if !node.isBlock() {
			return
		}
		// get the lines for the node
		segments := node.astNode.Lines()
		// iterate over the segments
		for i := 0; i < segments.Len(); i++ {
			segment := segments.At(i)
			// get the raw text of the segment
			txt := string(segment.Value(doc.source))
			hit := false
			// check for TODO
			if strings.Contains(txt, "TODO") {
				hit = true
			}
			// check for XXX
			if strings.Contains(txt, "XXX") {
				hit = true
			}
			if hit {
				Pl("TODO:", txt)
			}
		}
		return
	})
}

// showTerms shows all terms in the document that are not defined in
// the same document.
func (doc *Node) showTerms(fn string) {

	// walk the doc, collecting text blocks
	txts := make([]string, 0)
	doc.walk(func(node *Node) (stop bool) {
		if !node.isBlock() {
			return
		}
		// get the lines for the node
		segments := node.astNode.Lines()
		// iterate over the segments
		for i := 0; i < segments.Len(); i++ {
			segment := segments.At(i)
			// get the raw text of the segment
			txt := string(segment.Value(doc.source))
			txts = append(txts, txt)
		}
		return
	})

	// concatenate the text blocks into larger blocks of < 7000 tokens
	// - this is to avoid hitting the limit on the number of tokens
	//   that can be sent to an 8k grokker; we allow around 1k tokens
	//   for the response
	// - XXX we don't have a token counter yet, so we just use the
	//   length of the text / 3 as an estimate
	maxLen := 7000 / 3
	bigTxts := make([]string, 0)
	bigTxt := ""
	for _, txt := range txts {
		if len(bigTxt)+len(txt) < maxLen {
			bigTxt += txt
		} else {
			bigTxts = append(bigTxts, bigTxt)
			bigTxt = txt
		}
	}

	// send the prompt with each big text to grokker and show the results
	basename := filepath.Base(fn)
	prompt := Spf("List all terms used below that are not lay terms and that are not defined in %s.  Make the list in markdown format, one hyphen bullet point per term.", basename)
	Pl("undefined terms:")
	grok, _, _, _, err := grokker.Load()
	Ck(err)
	for _, bigTxt := range bigTxts {
		query := Spf("%s\n\n%s", prompt, bigTxt)
		res, err := grok.Answer(query, false)
		Ck(err)
		Pl(res)
	}
}
