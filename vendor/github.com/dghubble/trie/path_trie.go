package trie

// PathTrie is a trie of paths with string keys and interface{} values.

// PathTrie is a trie of string keys and interface{} values. Internal nodes
// have nil values so stored nil values cannot be distinguished and are
// excluded from walks. By default, PathTrie will segment keys by forward
// slashes with PathSegmenter (e.g. "/a/b/c" -> "/a", "/b", "/c"). A custom
// StringSegmenter may be used to customize how strings are segmented into
// nodes. A classic trie might segment keys by rune (i.e. unicode points).
type PathTrie struct {
	Segmenter StringSegmenter // key segmenter, must not cause heap allocs
	Value     interface{}
	Children  map[string]*PathTrie
}

// PathTrieConfig for building a path trie with different segmenter
type PathTrieConfig struct {
	Segmenter StringSegmenter
}

// NewPathTrie allocates and returns a new *PathTrie.
func NewPathTrie() *PathTrie {
	return &PathTrie{
		Segmenter: PathSegmenter,
	}
}

// NewPathTrieWithConfig allocates and returns a new *PathTrie with the given *PathTrieConfig
func NewPathTrieWithConfig(config *PathTrieConfig) *PathTrie {
	segmenter := PathSegmenter
	if config != nil && config.Segmenter != nil {
		segmenter = config.Segmenter
	}

	return &PathTrie{
		Segmenter: segmenter,
	}
}

// Get returns the value stored at the given key. Returns nil for internal
// nodes or for nodes with a value of nil.
func (trie *PathTrie) Get(key string) interface{} {
	node := trie
	for part, i := trie.Segmenter(key, 0); part != ""; part, i = trie.Segmenter(key, i) {
		node = node.Children[part]
		if node == nil {
			return nil
		}
	}
	return node.Value
}

// Put inserts the value into the trie at the given key, replacing any
// existing items. It returns true if the put adds a new value, false
// if it replaces an existing value.
// Note that internal nodes have nil values so a stored nil value will not
// be distinguishable and will not be included in Walks.
func (trie *PathTrie) Put(key string, value interface{}) bool {
	node := trie
	for part, i := trie.Segmenter(key, 0); part != ""; part, i = trie.Segmenter(key, i) {
		child, _ := node.Children[part]
		if child == nil {
			if node.Children == nil {
				node.Children = map[string]*PathTrie{}
			}
			child = NewPathTrie()
			node.Children[part] = child
		}
		node = child
	}
	// does node have an existing value?
	isNewVal := node.Value == nil
	node.Value = value
	return isNewVal
}

// Delete removes the value associated with the given key. Returns true if a
// node was found for the given key. If the node or any of its ancestors
// becomes childless as a result, it is removed from the trie.
func (trie *PathTrie) Delete(key string) bool {
	var path []nodeStr // record ancestors to check later
	node := trie
	for part, i := trie.Segmenter(key, 0); part != ""; part, i = trie.Segmenter(key, i) {
		path = append(path, nodeStr{part: part, node: node})
		node = node.Children[part]
		if node == nil {
			// node does not exist
			return false
		}
	}
	// delete the node value
	node.Value = nil
	// if leaf, remove it from its parent's children map. Repeat for ancestor path.
	if node.isLeaf() {
		// iterate backwards over path
		for i := len(path) - 1; i >= 0; i-- {
			parent := path[i].node
			part := path[i].part
			delete(parent.Children, part)
			if !parent.isLeaf() {
				// parent has other children, stop
				break
			}
			parent.Children = nil
			if parent.Value != nil {
				// parent has a value, stop
				break
			}
		}
	}
	return true // node (internal or not) existed and its value was nil'd
}

// Walk iterates over each key/value stored in the trie and calls the given
// walker function with the key and value. If the walker function returns
// an error, the walk is aborted.
// The traversal is depth first with no guaranteed order.
func (trie *PathTrie) Walk(walker WalkFunc) error {
	return trie.walk("", walker)
}

// WalkPath iterates over each key/value in the path in trie from the root to
// the node at the given key, calling the given walker function for each
// key/value. If the walker function returns an error, the walk is aborted.
func (trie *PathTrie) WalkPath(key string, walker WalkFunc) error {
	// Get root value if one exists.
	if trie.Value != nil {
		if err := walker("", trie.Value); err != nil {
			return err
		}
	}
	for part, i := trie.Segmenter(key, 0); ; part, i = trie.Segmenter(key, i) {
		if trie = trie.Children[part]; trie == nil {
			return nil
		}
		if trie.Value != nil {
			var k string
			if i == -1 {
				k = key
			} else {
				k = key[0:i]
			}
			if err := walker(k, trie.Value); err != nil {
				return err
			}
		}
		if i == -1 {
			break
		}
	}
	return nil
}

// PathTrie node and the part string key of the child the path descends into.
type nodeStr struct {
	node *PathTrie
	part string
}

func (trie *PathTrie) walk(key string, walker WalkFunc) error {
	if trie.Value != nil {
		if err := walker(key, trie.Value); err != nil {
			return err
		}
	}
	for part, child := range trie.Children {
		if err := child.walk(key+part, walker); err != nil {
			return err
		}
	}
	return nil
}

func (trie *PathTrie) isLeaf() bool {
	return len(trie.Children) == 0
}
