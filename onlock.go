package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/godbus/dbus"
)

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
	fmt.Printf("Screen locked!\n")
	exec.Command("/usr/local/bin/ckb-next", "--profile", "Off").Run()
}

func onScreenUnlock() {
	fmt.Printf("Screen unlocked!\n")
	exec.Command("/usr/local/bin/ckb-next", "--profile", "On").Run()
}
