package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	. "github.com/stevegt/goadapt"
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
	citations := doc.getSection("References")
	concepts.ckForwards(citations)
	concepts.ckUnrefs()
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
	retxt := Spf(`\[%s\]\[%s\]`, display, display)
	re := regexp.MustCompile(retxt)
	return re.MatchString(raw)
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

// mapCitations returns a map of citations in the given section
// func (section *Section) mapCitations() (citations map[string]*Node) {

// ckForwards checks for forward references and references that have
// no heading or citation
func (section *Section) ckForwards(citations *Section) {
	heads := section.mapHeadings()
	refs := section.mapReferences()
	// Pprint(refs)
	for key, ref := range refs {
		head, ok := heads[key]
		if !ok {
			if ref.isRef() {
				Pl("XXX check citation:", key)
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
	for key, _ := range heads {
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

/*
// ckrefs checks references and headings
func (doc *Node) ckrefs() {
	// get headings and references
	headings, refs := mapHeadRefs(doc)

	// ckref checks references and headings

	// check for unreferenced headings
	for key, heading := range headings {
		// skip if the heading is not in the concepts section

		if !heading.inConcepts {
			_, ok := refs[key]
			if !ok {
				Pf("unreferenced heading: %s\n", key)
			}
		}

		// check for forward references (references before headings)
		for key, link := range links {
			_, ok := headings[key]
			if !ok {
				Pf("link to missing heading: %s\n", key)
			}
		}

		// parseConcepts parses the Concepts section of the document
		// - returns a map of headings and a map of references

		// walk the tree looking for headings and references
		headings := make(map[string]*ast.Heading)
		links := make(map[string]*ast.Link)
		inConcepts := false
		walker := func(node ast.Node, entering bool) (ast.WalkStatus, error) {
			// Pf("node: %T, entering: %v\n", node, entering)
			if entering {
				switch node.(type) {
				case *ast.Heading:
					// get position of the heading
					lines := node.Lines()
					start := lines.At(0).Start
					Pf("start: %d\n", start)

					// add the heading to the map
					// - the key is the heading text
					// - the value is the heading node
					// - if the heading is already in the map, print a message
					key := string(node.Text(txt.Source()))
					// look for the concepts section
					if key == "Concepts" {
						inConcepts = true
						return ast.WalkContinue, nil
					}
					// skip headings that are not in the concepts section
					if !inConcepts {
						return ast.WalkContinue, nil
					}
					// lowercase
					key = strings.ToLower(key)
					// replace spaces with hyphens
					key = strings.Replace(key, " ", "-", -1)
					_, ok := headings[key]
					if ok {
						Pf("Duplicate heading: %s\n", key)
					} else {
						// Pf("Adding heading: %s\n", key)
						headings[key] = node.(*ast.Heading)
					}
				case *ast.Link:
					// check the headings map for the reference
					// - if the reference is not in the map, print a message
					key := string(node.(*ast.Link).Destination)
					// skip if the link is not to a heading
					if key[0] != '#' {
						return ast.WalkContinue, nil
					}
					// remove the leading '#'
					key = key[1:]
					// add it to the links map
					links[key] = node.(*ast.Link)
					// Pf("Checking reference: %s\n", key)
					_, ok := headings[key]
					if ok {
						// Pf("reference after heading: %s\n", key)
					} else {
						Pf("reference before heading: %s\n", key)
					}
				}
			} else {
				switch node.(type) {
				case *ast.Heading:
					// watch for when we leave the concepts section
					key := string(node.Text(txt.Source()))
					if key == "Concepts" {
						inConcepts = false
					}
				}
			}
			return ast.WalkContinue, nil
		}

		err := ast.Walk(doc, walker)
		Ck(err)

	}
}
*/
