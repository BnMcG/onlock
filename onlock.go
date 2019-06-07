package main

import (
	"log"
	"os/exec"

	"github.com/godbus/dbus"
)

type commandToExecute struct {
	Path      string
	Arguments []string
}

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer conn.Close()

	match := dbus.WithMatchSender("org.freedesktop.ScreenSaver")
	conn.AddMatchSignal(match)

	dbusChannel := make(chan *dbus.Signal, 10)
	conn.Signal(dbusChannel)

	for m := range dbusChannel {
		if m.Name == "org.freedesktop.ScreenSaver.ActiveChanged" && m.Path == "/org/freedesktop/ScreenSaver" {
			active := m.Body[0].(bool)
			if active {
				onScreenLock()
			} else {
				onScreenUnlock()
			}
		}
	}
}

func onScreenLock() {
	commandsToExecute := []commandToExecute{
		{
			Path: "/usr/local/bin/ckb-next",
			Arguments: []string{
				"--profile",
				"Off",
			},
		},
	}

	for _, c := range commandsToExecute {
		exec.Command(c.Path, c.Arguments...).Run()
	}
}

func onScreenUnlock() {
	commandsToExecute := []commandToExecute{
		{
			Path: "/usr/local/bin/ckb-next",
			Arguments: []string{
				"--profile",
				"On",
			},
		},
	}

	for _, c := range commandsToExecute {
		exec.Command(c.Path, c.Arguments...).Run()
	}
}
