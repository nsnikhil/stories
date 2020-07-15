package main

type command struct {
	name string
	run  func()
}

func (c command) execute() {
	c.run()
}

func newCommand(name string, run func()) command {
	return command{
		name: name,
		run:  run,
	}
}
