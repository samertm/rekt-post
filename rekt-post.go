package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var _ = fmt.Println // debugging

type post struct {
	title string
	path  string
	freqs map[string]int
	edges []edge
}

type edge struct {
	posts  [2]post
	weight int
}

func newEdge(p0, p1 post) edge {
	return edge{posts: [2]post{p0, p1}}
}

func (l edge) contains(p post) bool {
	return l.posts[0].title == p.title || l.posts[1].title == p.title
}

type graph struct {
	verticies []post
	edges     []edge
}

func (g graph) containsEdge(p0, p1 post) bool {
	for _, l := range g.edges {
		if l.contains(p0) && l.contains(p1) {
			return true
		}
	}
	return false
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

func union(freq0, freq1 map[string]int) []string {
	u := make([]string, 0)
	for k, _ := range freq0 {
		if _, ok := freq1[k]; ok {
			u = append(u, k)
		}
	}
	return u
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateEdge(p0, p1 post) edge {
	e := newEdge(p0, p1)		
	keys := union(p0.freqs, p1.freqs)
	for _, k := range keys {
		e.weight += min(p0.freqs[k], p1.freqs[k])
	}
	return e
}

func createEdges(g *graph) []edge {
	p := g.verticies
	edges := make([]edge, 0, 2*len(p)) // not sure how much of a cap i should allocate...
	for i := 0; i < len(p); i++ {
		for j := i + 1; j < len(p); j++ {
			edges = append(edges, generateEdge(p[i], p[j]))
		}
	}
	return edges
}

func main() {
	g := &graph{}
	g.verticies = makePosts("/home/samer/posts/")
	g.edges = createEdges(g)
}

