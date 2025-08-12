package utree

import (
	uconer "github.com/hxxshidage/myutils/container"
	"sync"
)

type TreeProxy[T any] struct {
	Nodes               []TreeNode[T]
	nodesIndex          map[string]TreeNode[T]
	builtNodesIndexOnce sync.Once

	parentsIndex          map[string]TreeNode[T]
	builtParentsIndexOnce sync.Once
}

func NewTreeProxy[T any, N TreeNode[T]](nodes []N) *TreeProxy[T] {
	realNodes := make([]TreeNode[T], len(nodes))
	for i := range nodes {
		realNodes[i] = nodes[i]
	}
	return &TreeProxy[T]{Nodes: realNodes}
}

func (tp *TreeProxy[T]) FindNode(key string) TreeNode[T] {
	return FindNode(key, tp.Nodes)
}

func (tp *TreeProxy[T]) FindParent(key string) TreeNode[T] {
	return FindParent(key, tp.Nodes)
}

func (tp *TreeProxy[T]) FindNodePath(key string) []TreeNode[T] {
	return FindNodePath(key, tp.Nodes)
}

func (tp *TreeProxy[T]) FindNodeAndParent(key string) (TreeNode[T], TreeNode[T], bool) {
	return FindNodeAndParent(key, tp.Nodes)
}

func (tp *TreeProxy[T]) FindNodes(keys []string) []TreeNode[T] {
	if len(keys) == 0 {
		return nil
	} else if len(keys) == 1 {
		return []TreeNode[T]{tp.FindNode(keys[0])}
	} else {
		tp.builtNodesIndexOnce.Do(func() {
			tp.nodesIndex = BuildNodeIndex[T](tp.Nodes)
		})

		foundNodes := make([]TreeNode[T], len(keys))
		for idx, key := range keys {
			foundNodes[idx] = tp.nodesIndex[key]
		}

		return foundNodes
	}
}

func (tp *TreeProxy[T]) FindParents(keys []string, parentRoot bool) []TreeNode[T] {
	if len(keys) == 0 {
		return nil
	}
	if len(keys) == 1 {
		parent := tp.FindParent(keys[0])
		if !parentRoot || parent == nil {
			return []TreeNode[T]{parent}
		}

		return []TreeNode[T]{tp.findRootParent(parent)}
	}

	tp.builtParentsIndexOnce.Do(func() {
		tp.parentsIndex = BuildParentIndex[T](tp.Nodes)
	})

	foundParents := make([]TreeNode[T], len(keys))
	for idx, key := range keys {
		if parent, exists := tp.parentsIndex[key]; exists {
			if parentRoot {
				foundParents[idx] = tp.findRootParent(parent)
			} else {
				foundParents[idx] = parent
			}
		}
	}

	return foundParents
}

// findRootParent 递归查找根父节点
func (tp *TreeProxy[T]) findRootParent(node TreeNode[T]) TreeNode[T] {
	for {
		parent := tp.parentsIndex[node.GetKey()]
		if parent == nil /*|| parent.GetPid() == 0*/ {
			return node
		}
		node = parent
	}
}

func (tp *TreeProxy[T]) FindSiblings(key string) []TreeNode[T] {
	return FindSiblings[T](key, tp.Nodes)
}

func (tp *TreeProxy[T]) IsParent(key string) bool {
	pn := FindParent[T](key, tp.Nodes)
	return nil != pn && pn.GetPid() == 0
}

func (tp *TreeProxy[T]) IsLeaf(key string) bool {
	node := tp.FindNode(key)
	return nil != node && len(node.GetChildren()) == 0
}

func (tp *TreeProxy[T]) IsLeafNode(node TreeNode[T]) bool {
	return nil != node && len(node.GetChildren()) == 0
}

func (tp *TreeProxy[T]) IsParentNode(node TreeNode[T]) bool {
	return node.GetPid() == 0
}

func (tp *TreeProxy[T]) Walk(fn WalkFunc[T]) {
	Walk[T](tp.Nodes, fn)
}

func (tp *TreeProxy[T]) WalkFast(fn WalkFastFunc[T]) {
	walkNodesFast(tp.Nodes, 0, fn)
}

func (tp *TreeProxy[T]) WalkSubtree(parent TreeNode[T], fn WalkFastFunc[T]) {
	WalkSubtree[T](parent, fn)
}

// 遍历单个父节点的直接子树
func (tp *TreeProxy[T]) WalkChildren(parent TreeNode[T], fn WalkFastFunc[T]) {
	WalkChildren[T](parent, fn)
}

func (tp *TreeProxy[T]) WalkWithControl(wcFn WalkControlFunc[T]) {
	WalkWithControl[T](tp.Nodes, wcFn)
}

type MarkFunc[T any] func(TreeNode[T])

// 标记, 当keys中包含父节点, 需要标记其下的所有子节点
func (tp *TreeProxy[T]) MarkAs(keys []string, mf MarkFunc[T]) {
	// 构建索引
	tp.builtNodesIndexOnce.Do(func() {
		tp.nodesIndex = BuildNodeIndex[T](tp.Nodes)
	})

	marked := uconer.NewSet[string]()

	var markChildren func(TreeNode[T])
	markChildren = func(node TreeNode[T]) {
		if marked.Contains(node.GetKey()) {
			return
		}
		marked.Add(node.GetKey())

		mf(node)

		for _, child := range node.GetChildren() {
			markChildren(child)
		}
	}

	for _, key := range keys {
		if node, exists := tp.nodesIndex[key]; exists {
			markChildren(node)
		}
	}
}

func (tp *TreeProxy[T]) MarkAsPlus(keys []string, mf MarkFunc[T]) {
	// 构建索引
	tp.builtNodesIndexOnce.Do(func() {
		tp.nodesIndex = BuildNodeIndex[T](tp.Nodes)
	})

	marked := uconer.NewSetWithSlice[string](keys)

	tp.WalkWithControl(func(node, parent TreeNode[T], depth int) bool {
		// 判断当前节点是否需要处理（自己是目标节点或是目标节点的后代）
		shouldMark := marked.Contains(node.GetKey())
		if parent != nil {
			parentMarked := marked.Contains(parent.GetKey())
			shouldMark = shouldMark || parentMarked
		}

		if shouldMark {
			mf(node)
			// 自动将当前节点加入marked，使其子节点也被处理
			marked.Add(node.GetKey())
		}

		// 继续遍历所有节点
		return true
	})
}
