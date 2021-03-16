package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/itchyny/volume-go"
	"github.com/jasonlvhit/gocron"
	hook "github.com/robotn/gohook"
)

// min and max values for the Y of menu
const maxY = 20
const minY = 0

func main() {
	updateStatus()  // updateStatus at startup so menu appears right awway
	go clickEvent() // running clickEvent in go routine

	// creating cron to run updateStatus every secound on the secound
	s := gocron.NewScheduler()
	s.Every(1).Seconds().Do(updateStatus)
	<-s.Start()
}

// updateStatus() will update the status bar
func updateStatus() {

	date := getDate()
	bat := getBat()
	wifi := getWifi()
	volume := getVolume()
	seperator := " ┃ "

	statusStr := seperator + volume + seperator + wifi + seperator + bat + seperator + date + seperator

	err := exec.Command("xsetroot", "-name", statusStr).Run()
	if err != nil {
		log.Fatalf("get xsetroot -name failed: %+v", err)
	}

}

// clickEvent() Always checks for clicks in the menu and handles them
func clickEvent() {
	EvChan := hook.Start()
	defer hook.End()

	for ev := range EvChan {

		if ev.Kind == hook.MouseDown && ev.Button == hook.MouseMap["left"] {
			x, y := robotgo.GetMousePos()
			if y <= maxY && y >= minY {
				if x <= 1564 && x >= 1521 { // wifi button
					spawn("wicd-curses")
				} else if x <= 1461 && x >= 1429 { // audio button
					spawn("pulsemixer")
				} else if x <= 1423 && x >= 1403 { // mute toggle
					muted, _ := volume.GetMuted()

					if muted {
						volume.Unmute()
					} else {
						volume.Mute()
					}

					updateStatus()
				} else if x <= 1483 && x >= 1466 { // volume up
					volume.IncreaseVolume(10)
					updateStatus()
				} else if x <= 1508 && x >= 1487 { // volume down
					volume.IncreaseVolume(-10)
					updateStatus()
				}
				//fmt.Println("You clicked the menu at ", fmt.Sprintf("X: %d, Y: %d", x, y))
			}
		}
	}
}

// spawn(cmd string) Spawn an Alacritty terminal abd rub the given command
func spawn(cmd string) {
	exec.Command("alacritty", "--command", cmd).Start()
}

// getVolume() return the volume in percntage with emojis for buttons
func getVolume() string {
	vol, err := volume.GetVolume()
	if err != nil {
		return err.Error()
	}

	muted, err := volume.GetMuted()
	if err != nil {
		return err.Error()
	}
	// buffer so always takes up same amount of space only works with mnospace fonts
	buffer := " "
	if vol < 10 {
		buffer = "   "
	} else if vol < 100 {
		buffer = "  "
	}

	// Update icon for colume levels and muted
	if muted {
		return fmt.Sprintf("🔇 %d%%%s🔺 🔻", vol, buffer)
	} else if vol < 30 {
		return fmt.Sprintf("🔈 %d%%%s🔺 🔻", vol, buffer)
	} else if vol > 30 && vol < 60 {
		return fmt.Sprintf("🔉 %d%%%s🔺 🔻", vol, buffer)
	} else {
		return fmt.Sprintf("🔊 %d%%%s🔺 🔻", vol, buffer)
	}
}

// getDate() returns a formatted date string with emoji
func getDate() string {

	return "📅 " + time.Now().Format("(Mon) 01/02/06 03:04:05 PM")
}

// getBat() returns the current battery percentage with emojis
func getBat() string {
	content, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/capacity")
	if err != nil {
		return err.Error()
	}

	levelStr := strings.Replace(string(content), "\n", "", -1)

	level, err := strconv.Atoi(levelStr)
	if err != nil {
		return err.Error()
	}

	// buffer so always takes up same amount of space, only works with monospace fonts
	buffer := ""
	if level < 10 {
		buffer = "  "
	} else if level < 100 {
		buffer = " "
	}

	return "🔋 " + levelStr + "%" + buffer
}

// getWifi() returns a red circle for no connection and green for connection
func getWifi() string {
	content, err := ioutil.ReadFile("/sys/class/net/wlan0/carrier")
	if err != nil {
		return err.Error()
	}

	if content[0] == '1' {
		return "📶 🟢"
	} else {
		return "📶 🔴"
	}

}
