package merkledag

type NodeBuilder struct {
	links

	// data is the raw node contents
	data []byte
}

func (nb *NodeBuilder) SetData(buf []byte) {
	// perform any validation here
	nb.data = buf
}
