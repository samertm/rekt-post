package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var _ = fmt.Printf // debugging
var _ = os.Open    // debugging

type post struct {
	title string
	path  string
	freqs map[string]int
}

type link struct {
	posts  [2]post
	weight int
}

type graph struct {
	verticies []post
	edges     []link
}

// cover thy eyes
var bannedWords = map[string]bool{"the": true, "i": true, "to": true, "a": true, "and": true, "of": true, "my": true, "in": true, "that": true, "it": true, "is": true, "for": true, "with": true, "on": true, "was": true, "im": true, "as": true, "be": true, "how": true, "but": true, "this": true, "you": true, "about": true, "have": true, "so": true, "more": true, "an": true, "me": true, "are": true, "if": true, "all": true, "at": true, "when": true, "because": true, "ive": true, "not": true, "like": true, "one": true, "get": true, "can": true, "which": true, "would": true, "after": true, "up": true, "what": true, "from": true, "by": true, "its": true, "before": true, "first": true, "had": true, "will": true, "img": true, "also": true, "he": true, "": true}

func makePosts(path string) []post {
	if path[len(path)-1] != '/' {
		path += string(append([]byte(path), '/'))
	}
	paths, err := filepath.Glob(path + "*.md")
	if err != nil {
		// wat do
		log.Fatal(err)
	}
	posts := make([]post, 0)
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		contents, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, makePost(string(contents)))
	}
	return posts
}

type lexer struct {
	contents []rune
	store    []rune
	next     chan rune
	out      chan token
}

type token struct {
	// types used: "attr", "word", "value", "eof"
	typ  string
	data string
}

func (l *lexer) run() chan token {
	// pump runs into next
	go func() {
		for _, r := range l.contents {
			l.next <- r
		}
		close(l.next)
	}()
	// turns input into tokens
	go func() {
		for s := top; s != nil; {
			s = s(l)
		}
		close(l.out)
	}()
	return l.out
}

func (l *lexer) emit(typ string) {
	s := strings.ToLower(string(l.store))
	l.store = make([]rune, 0)
	if bannedWords[s] {
		// only returns true if s is in the set of banned words
		return
	}
	l.out <- token{typ: typ, data: s}
}

func (l *lexer) add(r rune) {
	l.store = append(l.store, r)
}

func acceptRune(r rune, acceptset string) bool {
	for _, a := range acceptset {
		if a == r {
			return true
		}
	}
	return false
}

var (
	space      = " "
	tab        = "\t"
	newline    = "\n"
	whitespace = " \n\t"
	// went through all the non-letters on my keyboard ^_^
	junk = "\"!@#$%^&*()_+1234567890-=`~,./<>?;:[]{}\\|'"
)

// for lexing markdown
type mdState func(*lexer) mdState

// parsing title:, date:, etc...
// could be in either a word or an attr
func top(l *lexer) mdState {
	r, ok := <-l.next
	if !ok {
		return nil
	}
	if acceptRune(r, whitespace) {
		l.emit("word")
		return startWord
	}
	if acceptRune(r, ":") {
		l.emit("attr")
		return startAttrValue
	}
	if acceptRune(r, junk) {
		return top
	}
	l.add(r)
	return top
}

// eats all spaces (excluding whitespace) until attr value begins
func startAttrValue(l *lexer) mdState {
	r, ok := <-l.next
	if !ok {
		return nil
	}
	if acceptRune(r, space+tab) {
		// eat space
		return startAttrValue
	}
	if acceptRune(r, newline) {
		return top
	}
	// in attr value
	l.add(r)
	return attrValue
}

// gets attr value
func attrValue(l *lexer) mdState {
	r, ok := <-l.next
	if !ok {
		return nil
	}
	if acceptRune(r, newline) {
		l.emit("value")
		return top
	}
	l.add(r)
	return attrValue
}

// eats whitespace to word
func startWord(l *lexer) mdState {
	r, ok := <-l.next
	if !ok {
		return nil
	}
	if acceptRune(r, whitespace+junk) {
		return startWord
	}
	// we are in a word
	l.add(r)
	return word
}

func word(l *lexer) mdState {
	r, ok := <-l.next
	if !ok {
		return nil
	}
	if acceptRune(r, whitespace) {
		l.emit("word")
		return startWord
	}
	if acceptRune(r, junk) {
		// emit if "]" or ")"
		if acceptRune(r, "])") {
			l.emit("word")
			return startWord
		}
		return word
	}
	// got a char
	l.add(r)
	return word
}

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

func concatFreqs(freq1, freq2 map[string]int) map[string]int {
	catted := make(map[string]int)
	for k, v := range freq1 {
		catted[k] = v
	}
	for k, v := range freq2 {
		catted[k] += v
	}
	return catted
}

type freqPair struct {
	word string
	freq int
}

func sortFreqs(freqs map[string]int) []freqPair {
	var sortedAdd func([]freqPair, freqPair) []freqPair
	sortedAdd = func(fps []freqPair, fp freqPair) []freqPair {
		for i := range fps {
			if fp.freq > fps[i].freq {
				// add where i is
				return append(fps[:i], append([]freqPair{fp}, fps[i:]...)...)
			}
		}
		// add to end
		return append(fps, fp)
	}
	storeFreqs := make([]freqPair, 0, len(freqs))
	for k, v := range freqs {
		storeFreqs = sortedAdd(storeFreqs, freqPair{word: k, freq: v})
	}
	return storeFreqs
}

func main() {
	posts := makePosts("/home/samer/posts/")
	cats := make(map[string]int)
	for _, p := range posts {
		cats = concatFreqs(cats, p.freqs)
	}
	fs := sortFreqs(cats)
	for _, fp := range fs {
		fmt.Print("\"", fp.word, "\", ", fp.freq, "\n")
	}
}
