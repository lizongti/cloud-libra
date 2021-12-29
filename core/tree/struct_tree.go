package tree

type StructTree struct {
	data interface{}
}

func NewStructTree(data interface{}) *StructTree {
	return &StructTree{
		data: data,
	}
}

func (st *StructTree) Data() interface{} {
	return st.Data
}

func (st *StructTree) Get(path []string) interface{} {
	return nil
}

func (st *StructTree) get(source interface{}, path []string) interface{} {
	return nil
}
