package main

import (
	"github.com/thoj/go-ircevent"
	"fmt"
	"regexp"
	"math/rand"
	"os"
	"bufio"
)

var roomName = "#unix"

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
	defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	return lines, scanner.Err()
}

func main() {
	botName := "leiser"
	conn := irc.IRC(botName, botName)
	err  := conn.Connect("irc.esper.net:6667")

	quoteArray, err := readLines("random.txt")
	if err != nil {
		fmt.Println("Failed connecting")
		return
	}

	conn.AddCallback("001", func(e *irc.Event) {
		conn.Join(roomName)
	})

	conn.AddCallback("JOIN", func(e *irc.Event) {
		conn.Privmsg(roomName, "Hi, I'm leiser. I'm a bot")
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		matched, _ := regexp.MatchString("(hello|hi|hey) " + botName + ".*", e.Message())
		if matched {
			conn.Privmsg(roomName, "Hi, " + e.Nick)
		}
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		matched, _ := regexp.MatchString(botName + " quote me", e.Message())
		if matched {
			index := rand.Intn(len(quoteArray))
			quote := quoteArray[index]

			conn.Privmsg(roomName, quote)
		}
	})
	
	conn.Loop()
}
