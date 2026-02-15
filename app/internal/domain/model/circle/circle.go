package circle

import (
	"fmt"

	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
)

type Circle struct {
	id      CircleID
	name    CircleName
	owner   user.User
	members []user.User
}

func NewCircle(
	circleID CircleID,
	circleName CircleName,
	owner user.User,
	members []user.User,
) (*Circle, error) {
	return &Circle{
		id:      circleID,
		name:    circleName,
		owner:   owner,
		members: members,
	}, nil
}

func (c *Circle) CircleID() CircleID {
	return c.id
}

func (c *Circle) CircleName() CircleName {
	return c.name
}

func (c *Circle) Owner() user.User {
	return c.owner
}

func (c *Circle) Members() []user.User {
	return c.members
}

func (c *Circle) Notify(n CircleNotification) {
	n.ID(c.id).Name(c.name).Owner(c.owner).Members(c.members)
}

func (c *Circle) IsFull() bool {
	return c.CountMembers() >= 30
}

func (c *Circle) CountMembers() int {
	return len(c.members) + 1
}

func (c *Circle) Join(member *user.User) error {
	if member == nil {
		return fmt.Errorf("member is nil: %s", c.id.String())
	}
	if c.IsFull() {
		return fmt.Errorf("circle is full: %s", c.id.String())
	}
	c.members = append(c.members, *member)
	return nil
}
