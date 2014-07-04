package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

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
		// pass file contents onto makePost
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

func main() {
	posts := makePosts("/home/samer/posts/")
}
