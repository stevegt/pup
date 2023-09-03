package main

import (
	"bufio"
	"bytes"
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

	ckrefGm(inbuf)
}

// ckrefGm checks references and headings using goldmark
func ckrefGm(inbuf []byte) {

	// create a text.Reader from the input file
	// - this is required by the parser
	txt := text.NewReader(inbuf)

	// parse the input file into an AST node tree
	// p := parser.NewParser()
	// tree := p.Parse(txt)
	// tree.Dump(txt.Source(), 0)

	// md := goldmark.New()
	// var obuf bytes.Buffer
	// err = md.Convert(inbuf, &obuf)
	// Ck(err)

	p := goldmark.DefaultParser()
	doc := p.Parse(txt)

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

	// check for headings that are not in the links map
	for key, _ := range headings {
		_, ok := links[key]
		if !ok {
			Pf("unreferenced heading: %s\n", key)
		}
	}

	// check for heading links that are not in the headings map
	for key, _ := range links {
		_, ok := headings[key]
		if !ok {
			Pf("link to missing heading: %s\n", key)
		}
	}

}

// ckrefRe checks references and headings using regular
// expressions because goldmark and its docs are byzantine
func ckrefRe(inbuf []byte) {

	// heading regular expression
	hre := regexp.MustCompile(`(?m)^#+\s+(.*)\s*$`)

	// link regular expression
	lre := regexp.MustCompile(`(?m)\[.*\]\(#(.*)\)`)

	// scan the input file for headings and links to headings, one
	// line at a time
	// - if a heading is found, add it to a map
	// - if a link to a heading is found, check the map
	// - if the link is not in the map, print a message
	scanner := bufio.NewScanner(bytes.NewReader(inbuf))
	headings := make(map[string]*ast.Heading)
	for scanner.Scan() {
		line := string(scanner.Text())
		Pl("line:", line)

		// check for a heading
		hm := hre.FindStringSubmatch(line)
		if len(hm) > 0 {
			// found a heading
			// key is the heading text
			key := string(hm[1])
			// lowercase the heading text
			key = strings.ToLower(key)
			// replace spaces with hyphens
			key = strings.Replace(key, " ", "-", -1)
			_, ok := headings[key]
			if ok {
				Pf("Duplicate heading: %s\n", key)
			} else {
				Pf("Adding heading: %s\n", key)
				headings[key] = nil
			}
		}

		// check for a link to a heading
		lms := lre.FindAllStringSubmatch(line, -1)
		for _, lm := range lms {
			// found a link to a heading
			// key is the link destination
			key := string(lm[1])
			Pf("found link: %s\n", key)
			_, ok := headings[key]
			if ok {
				Pf("link after heading: %s\n", key)
			} else {
				Pf("link before heading: %s\n", key)
			}
		}
	}

}
