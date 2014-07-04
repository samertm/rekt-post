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
