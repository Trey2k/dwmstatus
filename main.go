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

const maxY = 20
const minY = 0

func main() {
	go clickEvent()

	s := gocron.NewScheduler()
	s.Every(1).Seconds().Do(updateStatus)
	<-s.Start()
}

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

func clickEvent() {
	EvChan := hook.Start()
	defer hook.End()

	for ev := range EvChan {

		if ev.Kind == hook.MouseDown && ev.Button == hook.MouseMap["left"] {
			x, y := robotgo.GetMousePos()
			if y <= maxY && y >= minY {
				if x <= 1564 && x >= 1521 {
					spawn("wicd-curses")
				} else if x <= 1461 && x >= 1429 {
					spawn("pulsemixer")
				} else if x <= 1423 && x >= 1403 {
					muted, _ := volume.GetMuted()

					if muted {
						volume.Unmute()
					} else {
						volume.Mute()
					}

					updateStatus()
				} else if x <= 1483 && x >= 1466 {
					volume.IncreaseVolume(10)
					updateStatus()
				} else if x <= 1508 && x >= 1487 {
					volume.IncreaseVolume(-10)
					updateStatus()
				}
				fmt.Println("You clicked the menu at ", fmt.Sprintf("X: %d, Y: %d", x, y))

			}
		}
	}
}

func spawn(cmd string) {
	err := exec.Command("alacritty", "--command", cmd).Start()
	if err != nil {
		log.Fatalf("launching %s failed: %+v", cmd, err)
	}
}

func respToString(str []byte) string {
	return strings.Replace(string(str), string(rune(10)), "", -1)
}

func getVolume() string {
	vol, err := volume.GetVolume()
	if err != nil {
		log.Fatalf("get volume failed: %+v", err)
	}

	muted, err := volume.GetMuted()
	if err != nil {
		log.Fatalf("get volume failed: %+v", err)
	}
	buffer := " "
	if vol < 10 {
		buffer = "   "
	} else if vol < 100 {
		buffer = "  "
	}

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

func getDate() string {

	return "📅 " + time.Now().Format("(Mon) 01/02/06 03:04:05 PM")
}

func getBat() string {
	content, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/capacity")
	if err != nil {
		log.Fatalf("get batery percentage failed: %+v", err)
	}

	levelStr := respToString(content)

	level, err := strconv.Atoi(levelStr)
	if err != nil {
		log.Fatalf("get batery percentage failed: %+v", err)
	}

	buffer := ""
	if level < 10 {
		buffer = "  "
	} else if level < 100 {
		buffer = " "
	}

	return "🔋 " + levelStr + "%" + buffer
}

func getWifi() string {
	content, err := ioutil.ReadFile("/sys/class/net/wlan0/carrier")
	if err != nil {
		log.Fatalf("get wifi status failed: %+v", err)
	}

	if content[0] == '1' {
		return "📶 🟢"
	} else {
		return "📶 🔴"
	}

}
