package znet

type NetIdAllocator struct {
	id uint32
}

func NewNetIdAllocator(init_id uint32) *NetIdAllocator {
	return &NetIdAllocator{
		id: 0,
	}
}

func (this *NetIdAllocator) NextId() uint32 {
	if this.id < 0 {
		this.id = 0
	}
	this.id ++

	return this.id
}
