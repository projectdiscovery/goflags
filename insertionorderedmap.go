package goflags

type InsertionOrderedMap struct {
	values map[string]*flagData
	keys   []string `yaml:"-"`
}

func (insertionOrderedMap *InsertionOrderedMap) forEach(fn func(key string, data *flagData)) {
	for _, key := range insertionOrderedMap.keys {
		fn(key, insertionOrderedMap.values[key])
	}
}

func (insertionOrderedMap *InsertionOrderedMap) Set(key string, value *flagData) {
	_, present := insertionOrderedMap.values[key]
	insertionOrderedMap.values[key] = value
	if !present {
		insertionOrderedMap.keys = append(insertionOrderedMap.keys, key)
	}
}
