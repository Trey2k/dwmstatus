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
)

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

	buffer := "    "

	err := exec.Command("xsetroot", "-name", buffer+volume+" "+wifi+" "+bat+" "+date+buffer).Run()
	if err != nil {
		log.Fatalf("get xsetroot -name failed: %+v", err)
	}

}

func clickEvent() {
	for true {
		mleft := robotgo.AddEvent("mleft")
		if mleft {
			x, y := robotgo.GetMousePos()
			if y <= 20 && y >= 0 && x <= 1920 && x >= 1450 {
				if x <= 1589 && x >= 1553 {
					spawn("wicd-curses")
				} else if x <= 1539 && x >= 1493 {
					spawn("pulsemixer")
				} else {
					fmt.Println("You clicked the menu at ", fmt.Sprintf("X: %d, Y: %d", x, y))
				}
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
	buffer := ""
	if vol < 10 {
		buffer = "  "
	} else if vol < 100 {
		buffer = " "
	}

	if muted {
		return fmt.Sprintf("🔇 %d%%", vol) + buffer
	} else if vol < 30 {
		return fmt.Sprintf("🔈 %d%%", vol) + buffer
	} else if vol > 30 && vol < 60 {
		return fmt.Sprintf("🔉 %d%%", vol) + buffer
	} else {
		return fmt.Sprintf("🔊 %d%%", vol) + buffer
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

	switch content[0] {
	case '1':
		return "📶 🟢"
	default:
		return "📶 🔴"
	}

}
