package nexar

import "strings"

type Tree struct {
	root *TreeNode
}

type TreeNode struct {
	staticChildren map[string]*TreeNode // exact matches
	paramChild     *TreeNode 
	paramName string
	handler        func(c *Context) *Context
}

func New() *Tree {
	return &Tree{
		root: &TreeNode{
			staticChildren: make(map[string]*TreeNode),
		},
	}
}

func (t *Tree) AddNode (paths []string, fn func(c *Context) *Context) {
	node := t.root
	for _, path := range paths {
		if path != "" && path[0] == ':' {
			if node.paramChild == nil {
				node.paramChild = &TreeNode{
					staticChildren: make(map[string]*TreeNode),
					paramName: path[1:],
				}
			}

			node = node.paramChild

			continue
		} 


		// static segment
		if node.staticChildren == nil {
			node.staticChildren = make(map[string]*TreeNode)
		}

		if node.staticChildren[path] == nil {
			node.staticChildren[path] = &TreeNode{
				staticChildren: make(map[string]*TreeNode),
			}
		}

		node = node.staticChildren[path]
	}


	node.handler = fn
}

func (t *Tree) FindNode(paths []string) (*TreeNode, map[string]string) {
	node := t.root
	params := make(map[string]string)

	for _, path := range paths {
		if child, ok := node.staticChildren[path]; ok {
			node = child
		} else if node.paramChild != nil {
			node = node.paramChild
			params[node.paramName] = path
		} else {
			return nil, nil
		}
	}

	return node, params
}

func (t *Tree) FindNodeByRoute(route string) (*TreeNode, map[string]string) {
	requestLineArr := strings.Split(route, "/")

	return t.FindNode(requestLineArr)
}