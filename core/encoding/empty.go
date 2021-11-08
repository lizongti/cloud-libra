package encoding

// var emptyChain = NewChain()

type emptyEncoding int

var empty emptyEncoding

func Empty() Encoding {
	return empty
}

func (emptyEncoding) String() string {
	return ""
}
func (emptyEncoding) Marshal(interface{}) ([]byte, error) {
	return nil, nil
}
func (emptyEncoding) Unmarshal([]byte, interface{}) error {
	return nil
}
func (emptyEncoding) Reverse() Encoding {
	return nil
}
