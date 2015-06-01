package osmap

// OSMap is a ordered string keyed map.
// It is eight to ten times slower than golang builtin map.
// Confirming my rule of thumb for hash-based vs. btree-based structures.

type OSMap struct {
	root   *node
	less   func(string, string) bool
	length int
}

type node struct {
	key, value  string
	red         bool
	left, right *node
}

// New returns an empty Map
func New() *OSMap {
	f := func(a, b string) bool {
		return a < b
	}
	return &OSMap{less: f}
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
	root := m.root
	for root != nil {
		if m.less(key, root.key) {
			root = root.left
		} else if m.less(root.key, key) {
			root = root.right
		} else {
			return root.value, true
		}
	}
	return "", false // string null value
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

// Do calls the given function on every key-value in the Map in order.
func (m *OSMap) Do(function func(string, string)) {
	do(m.root, function)
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

func (m *OSMap) insert(root *node, key, value string, lvl int) (*node, bool) {
	inserted := false
	if root == nil { // If the key was in the tree it would belong here
		if m.length%1000 == 0 && m.length > 10 {
			// fmt.Printf("l%v-%v - ", lvl, m.length)
		}
		return &node{key: key, value: value, red: true}, true
	}
	if isRed(root.left) && isRed(root.right) {
		colorFlip(root)
	}
	if m.less(key, root.key) {
		root.left, inserted = m.insert(root.left, key, value, lvl+1)
	} else if m.less(root.key, key) {
		root.right, inserted = m.insert(root.right, key, value, lvl+1)
	} else { // The key is already in the tree so just replace its value
		root.value = value
	}
	if isRed(root.right) && !isRed(root.left) {
		root = rotateLeft(root)
	}
	if isRed(root.left) && isRed(root.left.left) {
		root = rotateRight(root)
	}
	return root, inserted
}

func isRed(root *node) bool {
	return root != nil && root.red
}

func colorFlip(root *node) {
	root.red = !root.red
	if root.left != nil {
		root.left.red = !root.left.red
	}
	if root.right != nil {
		root.right.red = !root.right.red
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

func do(root *node, function func(string, string)) {
	if root != nil {
		do(root.left, function)
		function(root.key, root.value)
		do(root.right, function)
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
