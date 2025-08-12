package uconer

// 基于切片和map实现的有序map
// notes. 不具备并发安全, 操作不具备任何原子性
type OrderedMap struct {
	keys    []string
	content map[string]any
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		keys:    make([]string, 0),
		content: make(map[string]any),
	}
}

func NewOrderedMapWithCap(capacity int) *OrderedMap {
	return &OrderedMap{
		keys:    make([]string, 0, capacity),
		content: make(map[string]any, capacity),
	}
}

func (om *OrderedMap) Set(key string, value any) {
	if _, exists := om.content[key]; !exists {
		om.keys = append(om.keys, key)
	}

	om.content[key] = value
}

func (om *OrderedMap) Get(key string) (any, bool) {
	val, exists := om.content[key]
	return val, exists
}

func (om *OrderedMap) SetIfAbsent(key string, val any) bool {
	if _, ok := om.Get(key); !ok {
		om.Set(key, val)
		return true
	}

	return false
}

type ComputeFunc func(key string) any

func (om *OrderedMap) ComputeIfAbsent(key string, fun ComputeFunc) any {
	if val, ok := om.Get(key); !ok {
		val = fun(key)
		om.Set(key, val)
		return val
	} else {
		return val
	}
}

func (om *OrderedMap) Delete(key string) {
	delete(om.content, key)

	for i, k := range om.keys {
		if k == key {
			om.keys = append(om.keys[:i], om.keys[i+1:]...)
			break
		}
	}
}

func (om *OrderedMap) Keys() []string {
	return om.keys
}

func (om *OrderedMap) Values() []any {
	values := make([]any, len(om.keys))

	for i, key := range om.keys {
		values[i] = om.content[key]
	}

	return values
}

func (om *OrderedMap) Len() int {
	return len(om.keys)
}

func (om *OrderedMap) Clear() {
	om.keys = make([]string, 0)
	om.content = make(map[string]any)
}

func (om *OrderedMap) ToMap() map[string]any {
	m := make(map[string]any, len(om.content))

	for k, v := range om.content {
		m[k] = v
	}

	return m
}

type TraverseFunc func(idx int, key string, val any) bool

func (om *OrderedMap) Range(fun TraverseFunc) {
	for idx, key := range om.keys {
		if !fun(idx, key, om.content[key]) {
			break
		}
	}
}
