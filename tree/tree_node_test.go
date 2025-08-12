package utree

import (
	"fmt"
	uio "myutils/io"
	uparse "myutils/parse"
	utype "myutils/type"
	"strings"
	"testing"
)

type PathInfo struct {
	PomPath    string
	IsSelected bool // 标为是否选中
}

type RecursionNode struct {
	Id       int
	Pid      int
	Title    string
	Key      string
	Data     *PathInfo
	Children []*RecursionNode
}

func (rsn *RecursionNode) GetId() int {
	return rsn.Id
}

func (rsn *RecursionNode) GetPid() int {
	return rsn.Pid
}

func (rsn *RecursionNode) GetKey() string {
	return rsn.Key
}

func (rsn *RecursionNode) GetChildren() []TreeNode[*PathInfo] {
	return ConvertToTreeNode[*PathInfo, *RecursionNode](rsn.Children)
}

func (rsn *RecursionNode) GetData() *PathInfo {
	return rsn.Data
}

var testNodes = []*RecursionNode{
	{
		Id:    1,
		Pid:   0,
		Title: "I am:0",
		Key:   "0",
		Data: &PathInfo{
			PomPath:    "/0",
			IsSelected: false,
		},
		Children: []*RecursionNode{
			{
				Id:    2,
				Pid:   1,
				Title: "I am:00",
				Key:   "0-0",
				Data: &PathInfo{
					PomPath:    "/0/0",
					IsSelected: false,
				},
				Children: []*RecursionNode{},
			},
			{
				Id:    3,
				Pid:   1,
				Title: "I am:02",
				Key:   "0-1",
				Data: &PathInfo{
					PomPath:    "/0/1",
					IsSelected: false,
				},
				Children: []*RecursionNode{},
			},
		},
	},
	{
		Id:    4,
		Pid:   0,
		Title: "I am:1",
		Key:   "1",
		Data: &PathInfo{
			PomPath:    "/1",
			IsSelected: false,
		},
		Children: []*RecursionNode{
			{
				Id:    5,
				Pid:   4,
				Title: "I am:10",
				Key:   "1-0",
				Data: &PathInfo{
					PomPath:    "/1/0",
					IsSelected: false,
				},
				Children: []*RecursionNode{
					{
						Id:    6,
						Pid:   5,
						Title: "I am:100",
						Key:   "1-0-0",
						Data: &PathInfo{
							PomPath:    "/1/0/0",
							IsSelected: false,
						},
						Children: []*RecursionNode{},
					},
					{
						Id:    7,
						Pid:   5,
						Title: "I am:101",
						Key:   "1-0-1",
						Data: &PathInfo{
							PomPath:    "/1/0/1",
							IsSelected: false,
						},
						Children: []*RecursionNode{
							{
								Id:    8,
								Pid:   7,
								Title: "I am:1010",
								Key:   "1-0-1-0",
								Data: &PathInfo{
									PomPath:    "/1/0/1/0",
									IsSelected: false,
								},
								Children: []*RecursionNode{
									{
										Id:    9,
										Pid:   8,
										Title: "I am:10100",
										Key:   "1-0-1-0-0",
										Data: &PathInfo{
											PomPath:    "/1/0/1/0/0",
											IsSelected: false,
										},
										Children: []*RecursionNode{},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Id:    10,
		Pid:   0,
		Title: "I am:2",
		Key:   "2",
		Data: &PathInfo{
			PomPath:    "/2",
			IsSelected: false,
		},
		Children: make([]*RecursionNode, 0),
	},
}
var treeProxy = NewTreeProxy[*PathInfo, *RecursionNode](testNodes)

func TestFindNode(t *testing.T) {
	node := treeProxy.FindNode("1-0-1-0")
	fmt.Printf("node:%v\n", node)

	node = treeProxy.FindNode("2")
	fmt.Printf("node:%v\n", node)
}

func TestTreeProxy_WalkFast(t *testing.T) {
	treeProxy.WalkFast(func(node TreeNode[*PathInfo], depth int) {
		fmt.Printf("node:%v, depth:%d, isRoot:%t, isLeaf:%t\n", node, depth, depth == 0, len(node.GetChildren()) == 0)
	})
}

func TestTreeProxy_Walk(t *testing.T) {
	treeProxy.Walk(func(node, parent TreeNode[*PathInfo], depth int) {
		fmt.Printf("node:%v, parent:%v, depth:%d, isRoot:%t, isLeaf:%t\n", node, parent, depth, depth == 0, len(node.GetChildren()) == 0)
	})
}

func TestTreeProxy_WalkSubTree(t *testing.T) {
	treeProxy.WalkSubtree(treeProxy.FindParent("1-0-0"), func(node TreeNode[*PathInfo], depth int) {
		fmt.Printf("node:%v, depth:%d\n", node, depth)
	})
}

func TestTreeProxy_WalkChildren(t *testing.T) {
	treeProxy.WalkChildren(treeProxy.FindParent("1-0-0"), func(node TreeNode[*PathInfo], depth int) {
		fmt.Printf("node:%v, depth:%d\n", node, depth)
	})
}

func TestTreeProxy_FindNodes(t *testing.T) {
	nodes := treeProxy.FindNodes([]string{"1-0-1-0", "2", "1-0-1-0-0"})
	for _, node := range nodes {
		fmt.Printf("found node:%v\n", node)
	}
}

func TestTreeProxy_FindParents(t *testing.T) {
	parents := treeProxy.FindParents([]string{"1-0-1-0", "2", "1-0-1-0-0"}, true)
	for _, parent := range parents {
		fmt.Printf("found parent:%v\n", parent)
	}
}

func TestMarkAsSelected(t *testing.T) {
	selectedKeys := []string{"1-0-0", "1-0-1-0"}

	treeProxy.Walk(func(node, parent TreeNode[*PathInfo], depth int) {
		for _, key := range selectedKeys {
			if node.GetKey() == key {
				pi := node.GetData()
				pi.IsSelected = true
				break
			}
		}
	})

	treeProxy.Walk(func(node, parent TreeNode[*PathInfo], depth int) {
		fmt.Printf("node:%v, parent:%v, depth:%d, isRoot:%t, nodePathInfo:%+v\n", node, parent, depth, depth == 0, *node.GetData())
	})
}

type TestRecord struct {
	Id       int
	Pid      int
	Title    string
	Path     string
	Type     int
	Icon     string
	IsSelect bool
}

func TestFlatMapTree(t *testing.T) {
	selectedKeys := []string{"1-0-0", "1-0-1-0"}

	treeProxy.Walk(func(node, parent TreeNode[*PathInfo], depth int) {
		for _, key := range selectedKeys {
			if node.GetKey() == key {
				pi := node.GetData()
				pi.IsSelected = true
				break
			}
		}
	})

	records := Flatmap[*PathInfo, TestRecord](treeProxy.Nodes, func(node TreeNode[*PathInfo]) TestRecord {
		recNode := node.(*RecursionNode)
		return TestRecord{
			Id:       node.GetId(),
			Pid:      node.GetPid(),
			Title:    recNode.Title,
			Path:     recNode.GetData().PomPath,
			Type:     1,
			Icon:     "",
			IsSelect: recNode.GetData().IsSelected,
		}
	})

	var appender strings.Builder
	for idx, rec := range records {
		nodeJson, err := uparse.FmtJson(&rec)
		if err != nil {
			panic(err)
		}

		appender.Write(nodeJson)
		if idx < len(records)-1 {
			appender.WriteString("\n")
		}
	}
	err := uio.WriteFile("records.txt", []byte(appender.String()))
	if err != nil {
		panic(err)
	}
}

type TestNodeData struct {
	Path     string
	Type     int
	IsSelect bool
	Icon     string
}

type TestNode struct {
	Id       int
	Pid      int
	Key      string
	Title    string
	Data     *TestNodeData
	Children []TreeNode[*TestNodeData]
}

func (tn *TestNode) GetId() int {
	return tn.Id
}

func (tn *TestNode) GetPid() int {
	return tn.Pid
}

func (tn *TestNode) GetKey() string {
	return tn.Key
}

func (tn *TestNode) GetChildren() []TreeNode[*TestNodeData] {
	return tn.Children
}

func (tn *TestNode) GetData() *TestNodeData {
	return tn.Data
}

func (tn *TestNode) SetChildren(children []TreeNode[*TestNodeData]) {
	tn.Children = children
}

func prepareNodes() []TreeNodeBuilder[*TestNodeData] {
	nodes, err := uio.ReadLinesAndParse[TreeNodeBuilder[*TestNodeData]]("records.txt", func(nodeJson string) (TreeNodeBuilder[*TestNodeData], error) {
		var record TestRecord
		if err := uparse.ParseJson([]byte(nodeJson), &record); err != nil {
			return nil, err
		}

		node := TestNode{
			Id:    record.Id,
			Pid:   record.Pid,
			Title: "I am:" + utype.I2s(record.Id),
			Key:   utype.I2s(record.Id),
			Data: &TestNodeData{
				Path:     record.Path,
				Type:     record.Type,
				IsSelect: record.IsSelect,
				Icon:     record.Icon,
			},
			Children: make([]TreeNode[*TestNodeData], 0),
		}

		return &node, nil
	})

	if err != nil {
		panic(err)
	}

	return nodes

}

func TestBuildTree(t *testing.T) {
	nodes := prepareNodes()

	treeNodes := BuildTree[*TestNodeData](nodes)
	Walk(treeNodes, func(node, parent TreeNode[*TestNodeData], depth int) {
		fmt.Printf("node:%+v, parent:%+v, data:%v\n", node, parent, *node.GetData())
	})
}

func TestMarkAs(t *testing.T) {
	nodes := prepareNodes()

	treeNodes := BuildTree[*TestNodeData](nodes)

	tp := NewTreeProxy[*TestNodeData, TreeNode[*TestNodeData]](treeNodes)

	tp.MarkAsPlus([]string{"2", "10", "8"}, func(node TreeNode[*TestNodeData]) {
		node.GetData().IsSelect = true
	})

	tp.Walk(func(node, parent TreeNode[*TestNodeData], depth int) {
		if node.GetData().IsSelect {
			fmt.Printf("node:%+v, test:%+v\n", node, node.GetData())
		}
	})

	records := Flatmap[*TestNodeData, TestRecord](tp.Nodes, func(node TreeNode[*TestNodeData]) TestRecord {
		recNode := node.(*TestNode)
		return TestRecord{
			Id:       node.GetId(),
			Pid:      node.GetPid(),
			Title:    recNode.Title,
			Path:     recNode.GetData().Path,
			Type:     1,
			Icon:     "",
			IsSelect: recNode.GetData().IsSelect,
		}
	})

	var appender strings.Builder
	for idx, rec := range records {
		nodeJson, err := uparse.FmtJson(&rec)
		if err != nil {
			panic(err)
		}

		appender.Write(nodeJson)
		if idx < len(records)-1 {
			appender.WriteString("\n")
		}
	}
	err := uio.WriteFile("records-1.txt", []byte(appender.String()))
	if err != nil {
		panic(err)
	}
}

func TestFindSiblings(t *testing.T) {
	nodes := prepareNodes()

	treeNodes := BuildTree[*TestNodeData](nodes)

	tp := NewTreeProxy[*TestNodeData, TreeNode[*TestNodeData]](treeNodes)

	sibNodes := tp.FindSiblings("6")

	println(sibNodes)

}
