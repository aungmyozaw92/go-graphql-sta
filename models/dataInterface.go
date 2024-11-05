package models

type Identifier interface {
	GetId() int
}

// interface for dataloader result
type Data interface {
	Identifier
	GetDefault(int) Data
}

func (user User) GetId() int {
	return user.ID
}

func (unit Unit) GetId() int {
	return unit.ID
}


func (category Category) GetId() int {
	return category.ID
}


func (p Product) GetId() int {
	return p.ID
}

// loader loading more than one model by one id
type RelatedData interface {
	GetReferenceId() int
}

func (i Image) GetReferenceId() int {
	return i.ReferenceID
}