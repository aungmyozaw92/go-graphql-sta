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