package main

import (
	"strings"
)

// cover thy eyes
var bannedWords = map[string]bool{"the": true, "i": true, "to": true, "a": true, "and": true, "of": true, "my": true, "in": true, "that": true, "it": true, "is": true, "for": true, "with": true, "on": true, "was": true, "im": true, "as": true, "be": true, "how": true, "but": true, "this": true, "you": true, "about": true, "have": true, "so": true, "more": true, "an": true, "me": true, "are": true, "if": true, "all": true, "at": true, "when": true, "because": true, "ive": true, "not": true, "like": true, "one": true, "get": true, "can": true, "which": true, "would": true, "after": true, "up": true, "what": true, "from": true, "by": true, "its": true, "before": true, "first": true, "had": true, "will": true, "img": true, "also": true, "he": true, "": true}

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

