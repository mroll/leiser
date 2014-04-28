package main

import (
	"strconv"
	"time"
	"fmt"
	"regexp"
	"math/rand"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"github.com/thoj/go-ircevent"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

var roomName = "#unix"
var weatherURL = "http://api.wunderground.com/api/06ae7ac7474e21f1/geolookup/conditions/q/"

func trimN(target string, i int) (string) {
	return target[:len(target) - i]
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	db := mysql.New("tcp", "", "127.0.0.1:3306", "leiser", "123botp@ss", "leiserDB")

	err := db.Connect()
	if err != nil {
		panic(err)
	}

	fmt.Println("Made database connection")
	taglines, _, _ := db.Query("select count(*) from taglines")

	botName := "leiser"

	var weatherRegexp = regexp.MustCompile("(?:leiser weather ([0-9]{5}))")
	var lookupRegexp = regexp.MustCompile("(?:^" + botName + " lookup ([a-zA-Z0-9]+))")
	var defineRegexp = regexp.MustCompile("^" + botName + " define ([a-zA-Z0-9]+): (.*)")
	var removeRegexp = regexp.MustCompile("^" + botName + " remove ([a-zA-Z0-9]+)")
	var dicerollRegexp = regexp.MustCompile("^" + botName + " roll me ([0-9]+)")
	var helpRegexp = regexp.MustCompile("^leiser (what can you do|help|man)")

	conn := irc.IRC(botName, botName)
	err  = conn.Connect("irc.esper.net:6667")

	conn.AddCallback("001", func(e *irc.Event) {
		conn.Join(roomName)
	})

	conn.AddCallback("JOIN", func(e *irc.Event) {
		conn.Privmsg(roomName, "Hi everybody")
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		matched, _ := regexp.MatchString("^[hH](ello|i|ey) " + botName + ".*", e.Message())
		if matched {
			conn.Privmsg(roomName, "Hi, " + e.Nick)
		}
	})

	// Help
	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		matched := helpRegexp.FindStringSubmatch(e.Message())
		if len(matched) == 2 {
			conn.Privmsg(roomName, "Define terms: leiser define <term>: <definition>")
			conn.Privmsg(roomName, "Lookup terms: leiser lookup <term>")
			conn.Privmsg(roomName, "Weather: leiser weather <zipcode>")
			conn.Privmsg(roomName, "Random Quote: leiser quote me")
			conn.Privmsg(roomName, "Dice Roll: leiser roll me <# of sides>")
		}
	})

	// Quote me
	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		matched, _ := regexp.MatchString("^" + botName + " quote me", e.Message())
		if matched {
			randID := rand.Intn(taglines[0].Int(0))
			rows, _, err := db.Query("select body from taglines where id = %d", randID)
			if err != nil {
				panic(err)
			}

			conn.Privmsg(roomName, rows[0].Str(0))
		}
	})

	// Define
	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		matched := defineRegexp.FindStringSubmatch(e.Message())
		if len(matched) == 3 {
			object := matched[1]
			def := matched[2]
			

			_, err := db.Start("insert into factoids (object, definition) values ('%s', '%s')", object, def)
			if err != nil {
				panic(err)
			}
			conn.Privmsg(roomName,object + ": " + def)
		}

	})

	// Remove factoid
	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		matched := removeRegexp.FindStringSubmatch(e.Message())
		if len(matched) == 2 {
			object := matched[1]
			_, err := db.Start("delete from factoids where object = '" + object + "'")
			if err != nil {
				conn.Privmsg(roomName, "Failed to remove " + object)
				return
			}

			conn.Privmsg(roomName, "Removed " + object + " from factoid database")
		}
	})

	// Lookup
	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		matched := lookupRegexp.FindStringSubmatch(e.Message())
		if len(matched) == 2 {
			object := matched[1]

			rows, _, _ := db.Query("select definition from factoids where object = '%s'", object)
			if len(rows) == 0 {
				conn.Privmsg(roomName, "Sorry " + e.Nick + ", I don't know that one")
				return
			}

			conn.Privmsg(roomName, rows[0].Str(0))

		}
	})

	// Weather
	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		matched := weatherRegexp.FindStringSubmatch(e.Message())
		if len(matched) == 2 {
			zip := matched[1]
			resp, err := http.Get(weatherURL + zip + ".json")
			if err != nil {
				conn.Privmsg(roomName, e.Nick + ", couldn't connect. Try again maybe?")
				return
			}

			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			var dat map[string]interface{}

			if err := json.Unmarshal(b, &dat); err != nil {
				fmt.Println(err)
				return
			}
			if len(dat) != 3 {
				conn.Privmsg(roomName, e.Nick + ", couldn't find anything with that zipcode")
				return
			}

			w, _ := json.Marshal(dat["current_observation"])
			if err := json.Unmarshal(w, &dat); err != nil {
				fmt.Println(err)
				return
			}

			conn.Privmsg(roomName, dat["temperature_string"].(string) + " " + dat["weather"].(string))
		}
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		matched := dicerollRegexp.FindStringSubmatch(e.Message())
		if len(matched) == 2 {
			sides, err := strconv.ParseInt(matched[1], 0, 64)
			if err != nil {
				conn.Privmsg(roomName, "Error parsing your input " + e.Nick)
				return
			}

			result := rand.Intn(int(sides))
			fmt.Println(result)
			conn.Privmsg(roomName, strconv.Itoa(result))
		}
	})

	conn.Loop()
}
