package osmaps

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/pbberlin/tools/util"
)

const (
	zeroVal   = ""
	cFanout   = 12   // key-values per node
	cPrealloc = 1500 // kv-slices pre-allocated on map-creation

	noSuccessor = -1 // constant indicating no successor
)

var (
	cntr   = 0
	ndStat = make([]map[*node]int, 20) // node statistics
)

func init() {
	for i := 0; i < len(ndStat); i++ {
		ndStat[i] = make(map[*node]int, 0)
	}
}

// key-value type
type kvt struct {
	key, val string
	succ     int
}

type OSMap struct {
	root   *node
	less   func(string, string) bool
	length int

	reservoir [][]kvt
	allocCntr int
}

type node struct {
	min, max string // key
	minIdx   int    // index to smallest key
	kv       []kvt

	red         bool
	left, right *node
}

// New returns an empty Map
func New() *OSMap {
	f := func(a, b string) bool {
		return a < b
	}
	m := &OSMap{less: f}

	x := make([][]kvt, cPrealloc)
	for i := 0; i < len(x); i++ {
		x[i] = make([]kvt, 0, cFanout)
	}
	m.reservoir = x
	return m
}

// Insert inserts a new key-value into the Map returning true;
// or replaces an existing key-value pair's value returning false.
//      inserted := myMap.Insert(key, value).
func (m *OSMap) Insert(key, value string) (inserted bool) {
	m.root, inserted = m.insert(m.root, key, value, 0)
	m.root.red = false
	if inserted {
		m.length++
	}
	return inserted
}

// For compatibility
func (m *OSMap) Set(key, value string) (inserted bool) { return m.Insert(key, value) }

// Find returns the value and true if the key is in the Map
// or nil and false otherwise.
//      value, found := myMap.Find(key).
func (m *OSMap) Find(key string) (value string, found bool) {
	nd := m.root
	for nd != nil {
		if m.less(key, nd.min) && nd.left != nil {
			nd = nd.left
		} else if m.less(nd.max, key) && nd.right != nil {
			nd = nd.right
		} else {
			_, value, found = findInner(nd, key)
			return
		}
	}
	return zeroVal, false // string null value
}

// For compatibility
func (m *OSMap) Get(key string) (value string, found bool) { return m.Find(key) }

// Delete deletes the key-value returning true,
// or does nothing returning false
//      deleted := myMap.Delete(key).
func (m *OSMap) Delete(key string) (deleted bool) {
	if m.root != nil {
		if m.root, deleted = m.remove(m.root, key); m.root != nil {
			m.root.red = false
		}
	}
	if deleted {
		m.length--
	}
	return deleted
}

// Do calls the given function
// on every key-value in the Map in order.
func (m *OSMap) Do(fct1 func(string, string)) {
	do(m.root, fct1)
}

// Len returns the number of key-value pairs in the map.
func (m *OSMap) Len() int {
	return m.length
}

// ========================================================
func INNER_CORE_FOR_EXTRACTION() {
	// insert(args)
	// do(args)
	// remove(args)

}

func findInner(nd *node, key string) (idx int, val string, found bool) {

	for i := 0; i < len(nd.kv); i++ {
		if nd.kv[i].key == key {
			return i, nd.kv[i].val, true
		}
	}
	return noSuccessor, zeroVal, false

}

func (m *OSMap) insert(nd *node, key, value string, lvl int) (*node, bool) {
	inserted := false

	if nd == nil {
		// If the key was in the tree it would belong here
		newNd := node{}
		newNd.red = true
		newNd.min = key
		newNd.max = key
		newNd.kv = make([]kvt, 0, cFanout)
		newNd.kv = append(newNd.kv, kvt{key, value, noSuccessor})
		return &newNd, true
	}

	if isRed(nd.left) && isRed(nd.right) {
		colorFlip(nd)
	}

	if m.less(key, nd.min) && nd.left != nil {
		// printStatus("branch left", key, lvl)
		nd.left, inserted = m.insert(nd.left, key, value, lvl+1)
	} else if m.less(nd.max, key) && nd.right != nil {
		// printStatus("branch rght", key, lvl)
		nd.right, inserted = m.insert(nd.right, key, value, lvl+1)
	} else {
		// printStatus("remain", key, lvl)
		idx, _, contains := findInner(nd, key)
		if contains {
			nd.kv[idx].val = value // update
		} else {
			inserted = insertLinkedList(nd, m.less, key, value)
		}

	}

	if len(nd.kv) > cFanout-1 {
		// -1 | split one step *before* reaching capacity
		x := nd.right
		nd.right = m.split(nd, lvl)
		nd.right.right = x
	}

	//

	if isRed(nd.right) && !isRed(nd.left) {
		nd = rotateLeft(nd)
	}
	if isRed(nd.left) && isRed(nd.left.left) {
		nd = rotateRight(nd)
	}
	return nd, inserted
}

func (m *OSMap) split(nd *node, lvl int) *node {

	sortedK := make([]string, len(nd.kv))
	for i := 0; i < len(nd.kv); i++ {
		sortedK[i] = nd.kv[i].key
	}
	sort.Strings(sortedK)

	halfIdx := len(nd.kv) / 2
	splitkey := sortedK[halfIdx]

	ndStat[lvl][nd]++

	if lvl == 4 || lvl == 6 {
		kd := make([]string, 3)
		kd[0], kd[1], kd[2] = sortedK[0], splitkey, sortedK[len(sortedK)-1]
		for i := 0; i < len(kd); i++ {
			kd[i] = kd[i][:util.Min(len(kd[i]), 3)]
		}
		// fmt.Printf("splitting l%2v ac%3v %4q <  %4q < %4q \n", lvl, m.allocCntr, kd[0], kd[1], kd[2])
	}

	// kv1 := make([]kvt, 0, cFanout)
	// kv2 := make([]kvt, 0, cFanout)

	kv1 := m.reservoir[m.allocCntr]
	m.allocCntr++
	kv2 := m.reservoir[m.allocCntr]
	m.allocCntr++

	for i := 0; i < len(nd.kv); i++ {
		if m.less(nd.kv[i].key, splitkey) {
			kv1 = append(kv1, nd.kv[i])
		} else {
			kv2 = append(kv2, nd.kv[i])
		}
	}

	nd.min = sortedK[0]
	nd.max = sortedK[halfIdx-1]
	nd.kv = kv1

	newNd := node{}
	newNd.red = true
	newNd.min = splitkey
	newNd.max = sortedK[len(sortedK)-1]
	newNd.kv = kv2

	return &newNd

}

func isRed(nd *node) bool {
	return nd != nil && nd.red
}

func colorFlip(nd *node) {
	nd.red = !nd.red
	if nd.left != nil {
		nd.left.red = !nd.left.red
	}
	if nd.right != nil {
		nd.right.red = !nd.right.red
	}
}

func rotateLeft(root *node) *node {
	x := root.right
	root.right = x.left
	x.left = root
	x.red = root.red
	root.red = true
	return x
}

func rotateRight(root *node) *node {
	x := root.left
	root.left = x.right
	x.right = root
	x.red = root.red
	root.red = true
	return x
}

func do(nd *node, fct1 func(string, string)) {
	if nd != nil {

		do(nd.left, fct1)

		for i := 0; i < len(nd.kv); i++ {
			fct1(nd.kv[i].key, nd.kv[i].val)
		}

		do(nd.right, fct1)

	}
}

// We do not provide an exported First() method because this is an
// implementation detail.
func first(root *node) *node {
	for root.left != nil {
		root = root.left
	}
	return root
}

func (m *OSMap) remove(root *node, key string) (*node, bool) {
	deleted := false

	/*
		if m.less(key, root.key) {
			if root.left != nil {
				if !isRed(root.left) && !isRed(root.left.left) {
					root = moveRedLeft(root)
				}
				root.left, deleted = m.remove(root.left, key)
			}
		} else {
			if isRed(root.left) {
				root = rotateRight(root)
			}
			if !m.less(key, root.key) && !m.less(root.key, key) &&
				root.right == nil {
				return nil, true
			}
			if root.right != nil {
				if !isRed(root.right) && !isRed(root.right.left) {
					root = moveRedRight(root)
				}
				if !m.less(key, root.key) && !m.less(root.key, key) {
					smallest := first(root.right)
					root.key = smallest.key
					root.value = smallest.value
					root.right = deleteMinimum(root.right)
					deleted = true
				} else {
					root.right, deleted = m.remove(root.right, key)
				}
			}
		}
	*/
	return fixUp(root), deleted
}

func moveRedLeft(root *node) *node {
	colorFlip(root)
	if root.right != nil && isRed(root.right.left) {
		root.right = rotateRight(root.right)
		root = rotateLeft(root)
		colorFlip(root)
	}
	return root
}

func moveRedRight(root *node) *node {
	colorFlip(root)
	if root.left != nil && isRed(root.left.left) {
		root = rotateRight(root)
		colorFlip(root)
	}
	return root
}

func deleteMinimum(root *node) *node {
	if root.left == nil {
		return nil
	}
	if !isRed(root.left) && !isRed(root.left.left) {
		root = moveRedLeft(root)
	}
	root.left = deleteMinimum(root.left)
	return fixUp(root)
}

func fixUp(root *node) *node {
	if isRed(root.right) {
		root = rotateLeft(root)
	}
	if isRed(root.left) && isRed(root.left.left) {
		root = rotateRight(root)
	}
	if isRed(root.left) && isRed(root.right) {
		colorFlip(root)
	}
	return root
}

func printStatus(msg, key string, lvl int) {

	if cntr%1000 < 20 {
		kd := key
		if len(key) > 3 {
			kd = key[:3]
		}
		fmt.Printf("%-12v %5v l%2v\n", msg, kd, lvl)
	}
	cntr++
}

// is there a better way to convert a pointer to a hash
func pointerHash(p *string) (int64, int) {

	sp := fmt.Sprintf("%p", p) // string pointer

	pi, err := strconv.ParseInt(sp, 0, 64) // pointer integer
	if err != nil {
		fmt.Printf("Error converting pointer to int64:  %v - pointer was %v\n", err, sp)
		return 0, 0
	}

	mod := int(pi % (1000))

	cntr++
	if cntr%5000 == 0 {
		fmt.Printf("%v %v | ", pi, mod)
	}

	return pi, mod

}
