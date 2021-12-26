package main

import (
	"fmt"
	"strings"
	"sync"
)

/**
--- Day 12: Passage Pathing ---

With your submarine's subterranean subsystems subsisting suboptimally, the only way you're getting out of this cave
anytime soon is by finding a path yourself. Not just a path - the only way to know if you've found the best path is to
find all of them. Fortunately, the sensors are still mostly working, and so you build a rough map of the remaining
caves (your puzzle input). For example:

start-A
start-b
A-c
A-b
b-d
A-end
b-end

This is a list of how all of the caves are connected. You start in the cave named start, and your destination is the
cave named end. An entry like b-d means that cave b is connected to cave d - that is, you can move between them. So,
the above cave system looks roughly like this:

    start
    /   \
c--A-----b--d
    \   /
     end

Your goal is to find the number of distinct paths that start at start, end at end, and don't visit small caves more than
once. There are two types of caves: big caves (written in uppercase, like A) and small caves (written in lowercase,
like b). It would be a waste of time to visit any small cave more than once, but big caves are large enough that it
might be worth visiting them multiple times. So, all paths you find should visit small caves at most once, and can
visit big caves any number of times. Given these rules, there are 10 paths through this example cave system:

start,A,b,A,c,A,end
start,A,b,A,end
start,A,b,end
start,A,c,A,b,A,end
start,A,c,A,b,end
start,A,c,A,end
start,A,end
start,b,A,c,A,end
start,b,A,end
start,b,end

(Each line in the above list corresponds to a single path; the caves visited by that path are listed in the order they
are visited and separated by commas.) Note that in this cave system, cave d is never visited by any path: to do so,
cave b would need to be visited twice (once on the way to cave d and a second time when returning from cave d), and
since cave b is small, this is not allowed. Here is a slightly larger example:

dc-end
HN-start
start-kj
dc-start
dc-HN
LN-dc
HN-end
kj-sa
kj-HN
kj-dc

The 19 paths through it are as follows:

start,HN,dc,HN,end
start,HN,dc,HN,kj,HN,end
start,HN,dc,end
start,HN,dc,kj,HN,end
start,HN,end
start,HN,kj,HN,dc,HN,end
start,HN,kj,HN,dc,end
start,HN,kj,HN,end
start,HN,kj,dc,HN,end
start,HN,kj,dc,end
start,dc,HN,end
start,dc,HN,kj,HN,end
start,dc,end
start,dc,kj,HN,end
start,kj,HN,dc,HN,end
start,kj,HN,dc,end
start,kj,HN,end
start,kj,dc,HN,end
start,kj,dc,end

Finally, this even larger example has 226 paths through it:

fs-end
he-DX
fs-he
start-DX
pj-DX
end-zg
zg-sl
zg-pj
pj-he
RW-he
fs-DX
pj-RW
zg-RW
start-pj
he-WI
zg-he
pj-fs
start-RW

How many paths through this cave system are there that visit small caves at most once?
*/

func Twelve() interface{} {
	gc := GraphConstructor{
		vertexPairs: sync.Map{},
	}
	ch := make(chan string)
	go func(ch chan string) {
		defer close(ch)
		s := GetScanner(12)
		for s.Scan() {
			ch <- s.Text()
		}
	}(ch)

	// Process inputs
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(gc *GraphConstructor, inputs chan string, wg *sync.WaitGroup) {
		defer wg.Done()
		for i := range inputs {
			r := strings.SplitN(i, "-", 2)
			gc.AddPair(r[0], r[1])
		}
	}(&gc, ch, &wg)
	wg.Wait()

	// Create graph from GraphConstructor
	g := CaveGraph{
		adjacencyList: map[string][]string{},
	}
	for k := range getVertices(&gc) {
		a, ok := gc.vertexPairs.Load(k)
		if !ok {
			continue
		}
		g.adjacencyList[k] = a.(*Appender).GetSet()
	}

	// Use Start and recurse, send solutions to channel
	results := make(chan string, 100)
	wg = sync.WaitGroup{}
	wg.Add(len(g.GetStartEdges()))

	for _, e := range g.GetStartEdges() {
		go func(g *CaveGraph, wg *sync.WaitGroup, outputs chan string, i InterimPath) {
			defer wg.Done()
			recursePathsFinding(g, &i, outputs)
			fmt.Println(e, "defer wg.Done()")
		}(&g, &wg, results, InterimPath{
			path: []string{"start", e},
			included: map[string]bool{
				"start": true,
				e:       true,
			},
		})
	}

	count := 0
	rWg := sync.WaitGroup{}
	rWg.Add(1)
	go func(count *int, wg *sync.WaitGroup) {
		defer wg.Done()
		for r := range results {
			fmt.Println("out", r)
			*count++
		}
	}(&count, &rWg)

	wg.Wait()
	close(results)
	rWg.Wait()

	// 4885
	return count
}

func recursePathsFinding(g *CaveGraph, p *InterimPath, outputs chan string) {
	if p.isComplete() {
		outputs <- strings.Join(p.path, "-")
		return
	}

	for _, e := range g.adjacencyList[p.getLast()] {
		if p.canEnter(e) && e != p.getLast() {
			n := p.CopyAndAppend(e)
			recursePathsFinding(g, &n, outputs)
		}
	}
}

func getVertices(gc *GraphConstructor) chan string {
	vertexCh := make(chan string)
	go func(ch chan string) {
		defer close(vertexCh)
		gc.vertexPairs.Range(func(k, v interface{}) bool {
			vertexCh <- k.(string)
			return true
		})
	}(vertexCh)
	return vertexCh
}

type GraphConstructor struct {
	vertexPairs sync.Map // [string]Appender
}

func (gc *GraphConstructor) addDirected(x string, y string) {
	gc.vertexPairs.LoadOrStore(x, CreateAppender())
	i, _ := gc.vertexPairs.Load(x)
	a := i.(*Appender)
	a.Append(y)
}

func (gc *GraphConstructor) AddPair(x string, y string) {
	gc.addDirected(x, y)
	gc.addDirected(y, x)
}

// CaveGraph once constructed, is read-only
type CaveGraph struct {
	adjacencyList map[string][]string
}

func (cg CaveGraph) GetStartEdges() []string {
	return cg.adjacencyList["start"]
}

func isSmall(c string) bool {
	// ASCII numbers, is small iff lowercase, which starts at 97.
	return c[0] > 96
}

type InterimPath struct {
	path     []string
	included map[string]bool
}

func (p InterimPath) getLast() string {
	return p.path[len(p.path)-1]
}

func (i InterimPath) CopyAndAppend(a string) InterimPath {
	n := InterimPath{
		path:     append(i.path, a),
		included: map[string]bool{},
	}
	for _, e := range n.path {
		n.included[e] = true
	}
	return n
}

func (i InterimPath) isComplete() bool {
	return i.getLast() == "end"
}

func (i InterimPath) canEnter(c string) bool {
	if "start" == c {
		return false
	}
	visited, ok := i.included[c]

	// Can enter iff !ok, ok && !visited, ok && visited && !isSmall(c)
	return !ok || (ok && !visited) || (ok && visited && !isSmall(c))
}
