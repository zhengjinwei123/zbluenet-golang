package znet

type NetIdAllocator struct {
	id int64
}

func NewNetIdAllocator() *NetIdAllocator {
	return &NetIdAllocator{
		id: 0,
	}
}

func (this *NetIdAllocator) NextId() int64 {
	if this.id < 0 {
		this.id = 0
	}
	this.id ++

	return this.id
}
