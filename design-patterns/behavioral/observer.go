package main

import "fmt"

/*
OBSERVER
*/
type Observer interface {
	Update(video string)
}

/*
SUBJECT
*/
type Channel struct {
	subscribers []Observer
}

func (c *Channel) Subscribe(o Observer) {
	c.subscribers = append(c.subscribers, o)
}

func (c *Channel) Notify(video string) {
	for _, sub := range c.subscribers {
		sub.Update(video)
	}
}

/*
CONCRETE OBSERVER
*/
type User struct {
	name string
}

func (u *User) Update(video string) {
	fmt.Println(u.name, "notified about:", video)
}

func main() {
	channel := &Channel{}

	user1 := &User{name: "Alice"}
	user2 := &User{name: "Bob"}

	channel.Subscribe(user1)
	channel.Subscribe(user2)

	channel.Notify("Go Design Patterns")
}
