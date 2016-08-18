package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	element "github.com/omniscale/imposm3/element"
	imposm "github.com/omniscale/imposm3/parser/pbf"
	terminal "golang.org/x/crypto/ssh/terminal"
)

func main() {

	// ------------------------------------------------
	// record statistical data about PBF
	// ------------------------------------------------

	// store metrics
	var metric = metrics{
		ErrorCount:     0,
		WarningCount:   0,
		TotalNodes:     0,
		TotalWays:      0,
		TotalRelations: 0,
	}

	// bitmasks to store element ids
	masks := bitmasks{
		Nodes:        newMask(),
		Ways:         newMask(),
		Relations:    newMask(),
		NodeRefs:     newMask(),
		WayRefs:      newMask(),
		RelationRefs: newMask(),
	}

	// ------------------------------------------------
	// first pass: populate bitmasks
	// ------------------------------------------------

	// parser channels
	ch := newChannels()

	// open pbf file
	file := openFile()

	// synchronize goroutines
	wg := sync.WaitGroup{}
	wg.Add(1)

	// parser routines
	go populateBitmasks(ch, &masks, &metric, &wg)

	// imposm3 parser
	parser := imposm.NewParser(file, ch.Coords, nil, ch.Ways, ch.Relations)
	parser.Parse()
	wg.Wait()

	// ------------------------------------------------
	// second pass: perform linting
	// ------------------------------------------------

	// parser channels
	ch = newChannels()

	// re-open pbf file
	file.Close()
	file = openFile()
	defer file.Close()

	// synchronize goroutines
	wg = sync.WaitGroup{}
	wg.Add(3)

	// linter routines
	go lintNodes(ch, &masks, &metric, &wg)
	go lintWays(ch, &masks, &metric, &wg)
	go lintRelations(ch, &masks, &metric, &wg)

	// imposm3 parser
	parser = imposm.NewParser(file, ch.Coords, nil, ch.Ways, ch.Relations)
	parser.Parse()
	wg.Wait()

	// print metrics
	metric.Print()

	// exit(1) on error
	if metric.ErrorCount > 0 {
		os.Exit(1)
	}
}

func populateBitmasks(ch channels, masks *bitmasks, metric *metrics, wg *sync.WaitGroup) {
	var routines = &sync.WaitGroup{}
	routines.Add(2)

	// nodes
	go func() {
		for set := range ch.Coords {
			for _, node := range set {
				masks.Nodes.Insert(node.Id)
				metric.TotalNodes++
			}
		}
		routines.Done()
	}()

	// ways
	go func() {
		for set := range ch.Ways {
			for _, way := range set {
				masks.Ways.Insert(way.Id)
				metric.TotalWays++
				for _, ref := range way.Refs {
					masks.NodeRefs.Insert(ref)
				}
			}
		}
		routines.Done()
	}()

	routines.Wait()
	routines.Add(1)

	// relations
	go func() {
		for set := range ch.Relations {
			for _, relation := range set {
				masks.Relations.Insert(relation.Id)
				metric.TotalRelations++
				for _, member := range relation.Members {
					switch member.Type {
					case 0:
						masks.NodeRefs.Insert(member.Id)
					case 1:
						masks.WayRefs.Insert(member.Id)
					case 2:
						masks.RelationRefs.Insert(member.Id)
					}
				}
			}
		}
		routines.Done()
	}()

	routines.Wait()
	wg.Done()
}

// node linter
func lintNodes(ch channels, masks *bitmasks, metric *metrics, wg *sync.WaitGroup) {
	for set := range ch.Coords {
		for _, node := range set {
			if !masks.NodeRefs.Has(node.Id) {
				switch len(node.Tags) {
				case 0:
					metric.Warning("node %d not referenced by way/relation and has no tags\n", node.Id)
				case 1:
					if _, ok := node.Tags["created_by"]; ok {
						metric.Warning("node %d not referenced by way/relation and has no valid tags\n", node.Id)
					}
				}
			}
		}
	}
	wg.Done()
}

// way linter
func lintWays(ch channels, masks *bitmasks, metric *metrics, wg *sync.WaitGroup) {
	for set := range ch.Ways {
		for _, way := range set {

			// way has too few refs
			if len(way.Refs) < 2 {
				metric.Error("way %d invalid refcount %d\n", way.Id, len(way.Refs))
			}

			// way is missing a node
			for _, ref := range way.Refs {
				if !masks.Nodes.Has(ref) {
					metric.Error("way %d missing member node %d\n", way.Id, ref)
				}
			}
		}
	}
	wg.Done()
}

// relation linter
func lintRelations(ch channels, masks *bitmasks, metric *metrics, wg *sync.WaitGroup) {
	for set := range ch.Relations {
		for _, relation := range set {

			// missing relation members
			for _, member := range relation.Members {
				switch member.Type {
				case 0:
					if !masks.Nodes.Has(member.Id) {
						metric.Error("relation %d missing member node %d\n", relation.Id, member.Id)
					}
				case 1:
					if !masks.Ways.Has(member.Id) {
						metric.Error("relation %d missing member way %d\n", relation.Id, member.Id)
					}
				case 2:
					// super relations
					if !masks.Relations.Has(member.Id) {
						metric.Error("relation %d missing member relation %d\n", relation.Id, member.Id)
					}
				}
			}
		}
	}
	wg.Done()
}

// metrics
type metrics struct {
	ErrorCount     int
	WarningCount   int
	TotalNodes     int
	TotalWays      int
	TotalRelations int
}

var isTerminal = terminal.IsTerminal(int(os.Stdout.Fd()))

func (m *metrics) Error(format string, a ...interface{}) {
	if isTerminal {
		fmt.Fprintf(os.Stderr, "\033[0;31m")
	}
	fmt.Fprintf(os.Stdout, "error: ")
	fmt.Fprintf(os.Stdout, format, a...)
	m.ErrorCount++
	if isTerminal {
		fmt.Fprintf(os.Stderr, "\033[0m")
	}
}

func (m *metrics) Warning(format string, a ...interface{}) {
	if isTerminal {
		fmt.Fprintf(os.Stderr, "\033[0;33m")
	}
	fmt.Fprintf(os.Stdout, "warning: ")
	fmt.Fprintf(os.Stdout, format, a...)
	m.WarningCount++
	if isTerminal {
		fmt.Fprintf(os.Stderr, "\033[0m")
	}
}

func (m *metrics) Print() {
	if isTerminal {
		fmt.Fprintf(os.Stderr, "\033[1;37m")
	}
	fmt.Fprintf(os.Stderr, "ErrorCount: %d\n", m.ErrorCount)
	fmt.Fprintf(os.Stderr, "WarningCount: %d\n", m.WarningCount)
	fmt.Fprintf(os.Stderr, "TotalNodes: %d\n", m.TotalNodes)
	fmt.Fprintf(os.Stderr, "TotalWays: %d\n", m.TotalWays)
	fmt.Fprintf(os.Stderr, "TotalRelations: %d\n", m.TotalRelations)
	if isTerminal {
		fmt.Fprintf(os.Stderr, "\033[0m")
	}
}

// parser channels
type channels struct {
	Coords    chan []element.Node
	Nodes     chan []element.Node
	Ways      chan []element.Way
	Relations chan []element.Relation
}

// newChannels - constructor
func newChannels() channels {
	return channels{
		make(chan []element.Node, 100),
		make(chan []element.Node, 100),
		make(chan []element.Way, 100),
		make(chan []element.Relation, 100),
	}
}

// bitmasks to store ids
type bitmasks struct {
	Nodes        *bitmask
	Ways         *bitmask
	Relations    *bitmask
	NodeRefs     *bitmask
	WayRefs      *bitmask
	RelationRefs *bitmask
}

// simple bitmask
type bitmask struct {
	I map[uint64]uint64
}

// basic get/set methods
func (b *bitmask) Has(val int64) bool {
	var v = uint64(val)
	return (b.I[v/64] & (1 << (v % 64))) != 0
}
func (b *bitmask) Insert(val int64) {
	var v = uint64(val)
	b.I[v/64] |= (1 << (v % 64))
}
func newMask() *bitmask {
	return &bitmask{
		make(map[uint64]uint64),
	}
}

// convenience file to open pbf file from argv
func openFile() *imposm.Pbf {
	if len(os.Args) < 1 {
		log.Fatal("Invalid PBF File")
		os.Exit(2)
	}
	path := os.Args[1]
	file, err := imposm.Open(path)
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}
	return file
}
