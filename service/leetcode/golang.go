package leetcode

type Node struct {
	Id       int
	ParentId int
	Name     string
	Children []*Node
}

func list2Tree(list []*Node) *Node {
	dict := make(map[int]*Node)
	for _, item := range list {
		dict[item.Id] = item
	}
	for _, item := range list {
		parentId := item.ParentId
		parent := dict[parentId]
		if parent == nil {
			continue
		}
		parent.Children = append(parent.Children, item)
	}
	return dict[0]
}
