package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var _ = fmt.Println // debugging

type post struct {
	title string
	path  string
	// content pre-processing
	content string
	freqs   map[string]int
	edges   []*edge
}

func (p *post) String() string {
	return p.title
}

type edge struct {
	posts  [2]*post
	weight int
}

func (e *edge) String() string {
	return "[" + e.posts[0].title + " " +
		e.posts[1].title + "] " + fmt.Sprint(e.weight)
}

// generates a link for p
func (e *edge) Link(p *post, folderPath string) string {
	var otherPost *post
	if e.posts[0].title != p.title {
		otherPost = e.posts[0]
	} else {
		otherPost = e.posts[1]
	}
	s := "[" + otherPost.title + "](|filename|" +
		newPath(otherPost, folderPath) + ")"
	return s
}

func newEdge(p0, p1 *post) *edge {
	return &edge{posts: [2]*post{p0, p1}}
}

func (l edge) contains(p post) bool {
	return l.posts[0].title == p.title || l.posts[1].title == p.title
}

type graph struct {
	vertices []*post
	edges     []*edge
}

func (g graph) String() string {
	var s string
	for i := range g.vertices {
		s += fmt.Sprint(g.vertices[i])
	}
	for i := range g.edges {
		s += fmt.Sprint(g.edges[i])
	}
	return s
}

func (g graph) containsEdge(p0, p1 post) bool {
	for _, l := range g.edges {
		if l.contains(p0) && l.contains(p1) {
			return true
		}
	}
	return false
}

func makePosts(path string) []*post {
	if path[len(path)-1] != '/' {
		path += string(append([]byte(path), '/'))
	}
	paths, err := filepath.Glob(path + "*.md")
	if err != nil {
		// wat do
		log.Fatal(err)
	}
	posts := make([]*post, 0)
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

func makePost(contents string) *post {
	// need to lex the input: get the title and any other attributes :O
	lex := &lexer{
		contents: []rune(contents),
		next:     make(chan rune),
		out:      make(chan token),
	}
	toks := lex.run()
	p := &parser{toks: toks}
	post := parseTop(p)
	post.content = contents
	return post
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

func sortedAppendEdge(es []*edge, e *edge) []*edge {
	for i := range es {
		if e.weight > es[i].weight {
			// add where i is
			return append(es[:i], append([]*edge{e}, es[i:]...)...)
		}
	}
	// add to end
	return append(es, e)
}

func generateEdge(p0, p1 *post) *edge {
	e := newEdge(p0, p1)
	keys := union(p0.freqs, p1.freqs)
	for _, k := range keys {
		e.weight += min(p0.freqs[k], p1.freqs[k])
	}
	p0.edges = sortedAppendEdge(p0.edges, e)
	p1.edges = sortedAppendEdge(p1.edges, e)
	return e
}

func createEdges(g *graph) []*edge {
	p := g.vertices
	edges := make([]*edge, 0, 2*len(p)) // not sure how much of a cap i should allocate...
	for i := 0; i < len(p); i++ {
		for j := i + 1; j < len(p); j++ {
			edges = append(edges, generateEdge(p[i], p[j]))
		}
	}
	return edges
}

func newPath(p *post, folderPath string) string {
	return folderPath + strings.Replace(p.title, " ", "-", -1) + ".md"
}

func generatePosts(g *graph, folderPath string) {
	for _, p := range g.vertices {
		f, err := os.Create(newPath(p, folderPath))
		if err != nil {
			log.Fatal(err)
		}
		var links string
		for i := 0; i < len(p.edges) && i < 3; i++ {
			links += "* " + p.edges[i].Link(p, "/posts/") + "\n"
		}
		tagline := "\n*Similar Links*: (powered by [rekt-post](https://github.com/samertm/rekt-post))\n"
		_, err = f.WriteString(p.content + tagline + links)
		if err != nil {
			log.Fatal(err)
		}
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	g := &graph{}
	g.vertices = makePosts("/home/samer/posts/")
	g.edges = createEdges(g)
	generatePosts(g, "/home/samer/pelican/content/posts/")
}
