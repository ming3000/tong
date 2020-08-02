package tong

import "net/http"

// RouteInfo is detail info of request.
type RouteInfo struct {
	Method string
	Path   string
	Name   string
}

// Router is for request matching.
type Router struct {
	root   *treeNode
	routes map[string]*RouteInfo
}

// NewRouter returns a new Router instance.
func NewRouter() *Router {
	return &Router{
		root:   newTreeNode(),
		routes: map[string]*RouteInfo{},
	}
}

// Add registers a route for method and path with matching handler.
func (r *Router) Add(method, path string, h HandlerFunc) {
	r.root.insert(method, fixPath(path), h)
}

// Find a handler registered for method and path.
func (r *Router) Find(method, path string, ctx *Context) {
	h := r.root.search(method, fixPath(path))
	ctx.handler = h
}

type methodHandler struct {
	get  HandlerFunc
	post HandlerFunc
}

type treeNode struct {
	next          [128]*treeNode
	methodHandler *methodHandler
}

/** Initialize your data structure here. */
func newTreeNode() *treeNode {
	return &treeNode{next: [128]*treeNode{}, methodHandler: new(methodHandler)}
}

/** inserts a path & handler into the trie. */
func (t *treeNode) insert(method, path string, hand HandlerFunc) {
	cur := t
	for _, v := range path {
		if cur.next[v] == nil {
			cur.next[v] = newTreeNode()
		} // if>>
		cur = cur.next[v]
	} // for>
	if cur != nil {
		cur.addHandler(method, hand)
	} // if>
}

/** returns handler if the path is in the trie. */
func (t *treeNode) search(method, path string) HandlerFunc {
	cur := t
	for _, v := range path {
		if cur.next[v] == nil {
			return NotFoundHandler
		} // if>>
		cur = cur.next[v]
	} // for>
	if cur != nil {
		return cur.findHandler(method)
	} else {
		return NotFoundHandler
	} // else>
}

func (t *treeNode) addHandler(method string, h HandlerFunc) {
	switch method {
	case http.MethodGet:
		t.methodHandler.get = h
	case http.MethodPost:
		t.methodHandler.post = h
	}
}

func (t *treeNode) findHandler(method string) HandlerFunc {
	switch method {
	case http.MethodGet:
		return t.methodHandler.get
	case http.MethodPost:
		return t.methodHandler.post
	default:
		return NotFoundHandler
	}
}
