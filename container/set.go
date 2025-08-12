package uconer

type SimpleSet[T comparable] struct {
	setMap map[T]struct{}
}

func NewSetWithSlice[T comparable](elements []T) *SimpleSet[T] {
	m := make(map[T]struct{}, len(elements))
	for _, t := range elements {
		m[t] = struct{}{}
	}

	return &SimpleSet[T]{
		setMap: m,
	}
}

func NewSet[T comparable]() *SimpleSet[T] {
	return &SimpleSet[T]{
		setMap: make(map[T]struct{}),
	}
}

func (ss *SimpleSet[T]) Add(ele T) {
	ss.setMap[ele] = struct{}{}
}

func (ss *SimpleSet[T]) Added(ele T) bool {
	if _, ok := ss.setMap[ele]; !ok {
		ss.setMap[ele] = struct{}{}
		return true
	}

	return false
}

func (ss *SimpleSet[T]) Remove(ele T) {
	delete(ss.setMap, ele)
}

func (ss *SimpleSet[T]) Removed(ele T) bool {
	if _, ok := ss.setMap[ele]; ok {
		delete(ss.setMap, ele)
		return true
	}

	return false
}

func (ss *SimpleSet[T]) Contains(ele T) bool {
	_, ok := ss.setMap[ele]
	return ok
}

func (ss *SimpleSet[T]) Size() int {
	return len(ss.setMap)
}

func (ss *SimpleSet[T]) Items() []T {
	items := make([]T, 0, len(ss.setMap))

	for ele, _ := range ss.setMap {
		items = append(items, ele)
	}

	return items
}

type IteFunc[T comparable] func(T) bool

func (ss *SimpleSet[T]) Range(itf IteFunc[T]) {
	for ele, _ := range ss.setMap {
		if !itf(ele) {
			break
		}
	}
}
