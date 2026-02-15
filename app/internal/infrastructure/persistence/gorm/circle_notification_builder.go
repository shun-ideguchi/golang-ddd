package gorm

import (
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/circle"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
)

type CircleDataModelBuilder struct {
	id      circle.CircleID
	name    circle.CircleName
	owner   user.User
	members []user.User
}

func (b *CircleDataModelBuilder) ID(id circle.CircleID) circle.CircleNotification {
	b.id = id
	return b
}

func (b *CircleDataModelBuilder) Name(name circle.CircleName) circle.CircleNotification {
	b.name = name
	return b
}

func (b *CircleDataModelBuilder) Owner(owner user.User) circle.CircleNotification {
	b.owner = owner
	return b
}

func (b *CircleDataModelBuilder) Members(members []user.User) circle.CircleNotification {
	b.members = members
	return b
}

func (b *CircleDataModelBuilder) Build() *Circle {
	members := make([]User, len(b.members))
	for i, member := range b.members {
		members[i] = User{
			ID:    member.ID().String(),
			Name:  member.Name().String(),
			Email: member.Email().String(),
		}
	}
	return &Circle{
		ID:      b.id.String(),
		Name:    b.name.String(),
		Owner:   b.owner.ID().String(),
		Members: members,
	}
}
