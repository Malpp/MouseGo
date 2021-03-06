package main

import (
	"github.com/go-vgo/robotgo"
	"github.com/googollee/go-socket.io"
	"github.com/tarm/serial"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		so.On("click", func() {
			robotgo.MouseClick()
		})

		so.On("rclick", func() {
			robotgo.MouseClick("right", false)
		})

		so.On("move", func(data string) {
			xy := strings.Split(data, " ")
			x, _ := strconv.ParseFloat(xy[0], 64)
			y, _ := strconv.ParseFloat(xy[1], 64)
			mouseX, mouseY := robotgo.GetMousePos()
			robotgo.MoveMouse(mouseX+int(x), mouseY+int(y))
		})

		so.On("scroll", func(scrollData string) {
			direction := "up"
			scrollFloat, _ := strconv.ParseFloat(scrollData, 64)
			scroll := int(scrollFloat)
			if scroll < 0 {
				direction = "down"
			}
			robotgo.ScrollMouse(Abs(scroll), direction)
		})

		so.On("screen", func(scrollData string) {
			go SendCloseMonitorCommand()
		})
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("static")))
	log.Println("Serving at localhost...")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func SendCloseMonitorCommand() {
	c := &serial.Config{Name: "COM8", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Printf("%+v", err)
	}
	defer s.Close()
	time.Sleep(time.Second)
	s.Write([]byte("a"))
	log.Println("Sent screen command")
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
