package utree

type TreeNode[T any] interface {
	GetId() int

	GetPid() int

	GetKey() string

	GetChildren() []TreeNode[T]

	GetData() T
}

// 根据key查找相应节点
func FindNode[T any](key string, nodes []TreeNode[T]) TreeNode[T] {
	for _, node := range nodes {
		if node.GetKey() == key {
			return node
		}
		if found := FindNode(key, node.GetChildren()); found != nil {
			return found
		}
	}
	return nil
}

// 根据key查找相应父节点
func FindParent[T any](key string, nodes []TreeNode[T]) TreeNode[T] {
	for _, node := range nodes {
		for _, child := range node.GetChildren() {
			if child.GetKey() == key {
				return node
			}
		}
		if parent := FindParent(key, node.GetChildren()); parent != nil {
			return parent
		}
	}
	return nil
}

// 根据key检索路径
func FindNodePath[T any](key string, nodes []TreeNode[T]) []TreeNode[T] {
	var path []TreeNode[T]
	var dfs func([]TreeNode[T]) bool

	dfs = func(curNodes []TreeNode[T]) bool {
		for _, node := range curNodes {
			path = append(path, node)
			if node.GetKey() == key {
				return true
			}
			if dfs(node.GetChildren()) {
				return true
			}
			path = path[:len(path)-1]
		}
		return false
	}

	if dfs(nodes) {
		return path
	}
	return nil
}

// 根据key检索对应Node和其parent
func FindNodeAndParent[T any](key string, nodes []TreeNode[T]) (node TreeNode[T], parent TreeNode[T], found bool) {
	var dfs func(TreeNode[T], []TreeNode[T]) bool
	dfs = func(p TreeNode[T], curNodes []TreeNode[T]) bool {
		for _, n := range curNodes {
			if n.GetKey() == key {
				node, parent, found = n, p, true
				return true
			}
			if dfs(n, n.GetChildren()) {
				return true
			}
		}
		return false
	}

	dfs(nil, nodes)

	return
}

// go无法支持协变, 将子类手动转为父类
func ConvertToTreeNode[T any, N TreeNode[T]](nodes []N) []TreeNode[T] {
	result := make([]TreeNode[T], len(nodes))

	for i, n := range nodes {
		result[i] = n
	}

	return result
}

// 构建key2node的索引
func BuildNodeIndex[T any](nodes []TreeNode[T]) map[string]TreeNode[T] {
	key2node := make(map[string]TreeNode[T])
	var dfs func([]TreeNode[T])

	dfs = func(curNodes []TreeNode[T]) {
		for _, node := range curNodes {
			key2node[node.GetKey()] = node
			dfs(node.GetChildren())
		}
	}

	dfs(nodes)
	return key2node
}

// 构建key到parentNode的索引
func BuildParentIndex[T any](nodes []TreeNode[T]) map[string]TreeNode[T] {
	key2parentNode := make(map[string]TreeNode[T])
	var dfs func([]TreeNode[T])

	dfs = func(curNodes []TreeNode[T]) {
		for _, node := range curNodes {
			for _, child := range node.GetChildren() {
				key2parentNode[child.GetKey()] = node
			}
			dfs(node.GetChildren())
		}
	}

	dfs(nodes)

	return key2parentNode
}

// 根据key构建路径索引
func BuildPathIndex[T any](nodes []TreeNode[T]) map[string][]TreeNode[T] {
	pathIndex := make(map[string][]TreeNode[T])
	var curPath []TreeNode[T]

	var dfs func([]TreeNode[T])
	dfs = func(curNodes []TreeNode[T]) {
		for _, node := range curNodes {
			curPath = append(curPath, node)
			pathIndex[node.GetKey()] = append([]TreeNode[T]{}, curPath...)

			dfs(node.GetChildren())

			curPath = curPath[:len(curPath)-1]
		}
	}

	dfs(nodes)

	return pathIndex
}

// 查找兄弟节点
func FindSiblings[T any](key string, nodes []TreeNode[T]) []TreeNode[T] {
	parent := FindParent(key, nodes)

	if parent == nil {
		return nodes
	}

	return parent.GetChildren()
}

type WalkFunc[T any] func(node, parent TreeNode[T], depth int)

// 遍历树并携带父节点信息
func Walk[T any](nodes []TreeNode[T], fn WalkFunc[T]) {
	walkNodes[T](nodes, nil, 0, fn)
}

func walkNodes[T any](nodes []TreeNode[T], parent TreeNode[T], depth int, fn WalkFunc[T]) {
	for _, node := range nodes {
		fn(node, parent, depth)
		walkNodes(node.GetChildren(), node, depth+1, fn)
	}
}

// 不携带parent的快速遍历
type WalkFastFunc[T any] func(node TreeNode[T], depth int)

func WalkFast[T any](nodes []TreeNode[T], fn WalkFastFunc[T]) {
	walkNodesFast[T](nodes, 0, fn)
}

func walkNodesFast[T any](nodes []TreeNode[T], depth int, fn WalkFastFunc[T]) {
	for _, node := range nodes {
		fn(node, depth)
		walkNodesFast(node.GetChildren(), depth+1, fn)
	}
}

// 遍历单个父节点及其子树
func WalkSubtree[T any](parent TreeNode[T], fn WalkFastFunc[T]) {
	fn(parent, 0)
	walkNodesFast(parent.GetChildren(), 1, fn)
}

// 遍历节点子树
func WalkChildren[T any](parent TreeNode[T], fn WalkFastFunc[T]) {
	for _, child := range parent.GetChildren() {
		fn(child, 0)
	}
}

type WalkControlFunc[T any] func(node, parent TreeNode[T], depth int) bool

// 可中断的遍历
func WalkWithControl[T any](nodes []TreeNode[T], wcFn WalkControlFunc[T]) {
	var walk func([]TreeNode[T], TreeNode[T], int) bool
	walk = func(nodes []TreeNode[T], parent TreeNode[T], depth int) bool {
		for _, node := range nodes {
			if !wcFn(node, parent, depth) {
				return false
			}

			if !walk(node.GetChildren(), node, depth+1) {
				return false
			}
		}

		return true
	}

	walk(nodes, nil, 0)
}

type FlatMapper[T any, R any] func(TreeNode[T]) R

// 将树平铺成切片
func Flatmap[T any, R any](nodes []TreeNode[T], mapper FlatMapper[T, R]) []R {
	var mappedNodes []R
	WalkFast(nodes, func(node TreeNode[T], depth int) {
		mappedNodes = append(mappedNodes, mapper(node))
	})

	return mappedNodes
}

type TreeNodeBuilder[T any] interface {
	TreeNode[T]

	SetChildren([]TreeNode[T])
}

// 将切片数据转化为树
func BuildTree[T any](flatNodes []TreeNodeBuilder[T]) []TreeNode[T] {
	nodeMap := make(map[int]TreeNodeBuilder[T])

	for _, node := range flatNodes {
		nodeMap[node.GetId()] = node
	}

	var roots []TreeNode[T]
	for _, node := range flatNodes {
		pid := node.GetPid()
		if pid == 0 {
			roots = append(roots, node)
			continue
		}

		parent, exists := nodeMap[pid]
		if exists {
			parent.SetChildren(append(parent.GetChildren(), node))
		}
	}

	return roots
}
