package main

import (
	"log"
)

// parses markdown
func makePost(contents string) post {
	// need to lex the input: get the title and any other attributes :O
	lex := &lexer{
		contents: []rune(contents),
		next:     make(chan rune),
		out:      make(chan token),
	}
	toks := lex.run()
	p := &parser{toks: toks}
	return parseTop(p)
}

type parser struct {
	toks chan token
	// stack: push/pop onto/off right end
	oldToks []token
}

var tokEOF = token{typ: "eof"}

func (p *parser) next() token {
	if len(p.oldToks) != 0 {
		t := p.oldToks[len(p.oldToks)-1]
		p.oldToks = p.oldToks[:len(p.oldToks)-1]
		return t
	}
	if t, ok := <-p.toks; ok {
		return t
	}
	return tokEOF
}

func (p *parser) push(t token) {
	p.oldToks = append(p.oldToks, t)
}

// inefficient! >:O
func (p *parser) peek() token {
	t := p.next()
	p.push(t)
	return t
}

// abstraction on abstraction on abstraction
func (p *parser) acceptType(typ string) bool {
	return p.peek().typ == typ
}

func parseTop(p *parser) post {
	if p.acceptType("eof") {
		log.Fatal("unexpected eof")
	}
	po := post{freqs: make(map[string]int)}
	if p.acceptType("attr") {
		attr := p.next()
		if attr.data == "title" {
			if !p.acceptType("value") {
				log.Fatal("title must have a value")
			}
			val := p.next()
			po.title = val.data
		}
	}
	// eat other attrs
	for p.acceptType("attr") {
		p.next() // eat "attr"
		if p.acceptType("value") {
			p.next() // eat optional value
		}
	}
	// parse words
	for p.acceptType("word") {
		w := p.next()
		po.freqs[w.data]++
	}
	return po
}
