package main

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"github.com/dalinaum/gdgevent/event"
	"log"
	"net/textproto"
	"strings"
	"time"
)

type Bot struct {
	server string
	port   string
	nick   string
	user   string
	pass   string
	conn   *textproto.Conn
}

func NewBot() *Bot {
	return &Bot{
		server: "irc.ozinger.org",
		port:   "6668",
		nick:   "안드냥",
		pass:   "",
		conn:   nil,
		user:   "gdgand",
	}
}

func (bot *Bot) Connect() {
	conn, err := textproto.Dial("tcp", bot.server+":"+bot.port)
	if err != nil {
		log.Fatal("unable to connect to IRC server", err)
	}
	bot.conn = conn
	log.Printf("Connected to IRC server %s:%s\n", bot.server, bot.port)

	userCommand := fmt.Sprintf("USER %s 8 * :%s\n", bot.user, bot.user)
	bot.conn.PrintfLine(userCommand)
	bot.conn.PrintfLine("NICK " + bot.nick)
}

func (bot *Bot) Pong(token string) {
	request := fmt.Sprintf("PONG %s", token)
	bot.conn.PrintfLine(request)
}

func (bot *Bot) Close() {
	bot.conn.Close()
}

type Channel struct {
	bot     *Bot
	channel string
}

func (b *Bot) NewChannel(channel string) *Channel {
	return &Channel{
		b,
		channel,
	}
}

func (c *Channel) Talk(msg string) {
	text := fmt.Sprintf("PRIVMSG %s :%s", c.channel, msg)
	c.bot.conn.PrintfLine(text)
}

func (c *Channel) Op(user string) {
	request := fmt.Sprintf("MODE %s +o %s", c.channel, user)
	c.bot.conn.PrintfLine(request)
}

func (c *Channel) Log(db *sql.DB, nick string, message string) {
	fmt.Print("나이스 사장님 샷!")
	_, error := db.Exec("INSERT INTO ANDNYANG_LOG(date, channel, nick, message) VALUES (($1), ($2), ($3), ($4))", time.Now().UTC(), c.channel, nick, message)
	if error != nil {
		log.Print(error)
	}
}

type GDGChannel struct {
	*Channel
	chapterId string
}

func (b *Bot) NewGDGChannel(channel, chapterId string) *GDGChannel {
	return &GDGChannel{
		b.NewChannel(channel),
		chapterId,
	}
}

func (c *GDGChannel) Activities() {
	start := time.Unix(0, 0)
	for _, e := range event.GetGDGEvents(c.chapterId, start, start) {
		event := fmt.Sprintf("[%s] %s (https://developers.google.com%s)", e.GetStart(), e.Title, e.Link)
		c.Talk(event)
	}
}

func main() {
	ircbot := NewBot()
	ircbot.Connect()
	defer ircbot.Close()
	channels := [...]*GDGChannel{
		ircbot.NewGDGChannel("#gdgand", "115078626730671785458"),
		ircbot.NewGDGChannel("#gdgwomen", "108196114606467432743"),
	}

	db, error := sql.Open("postgres", "user=postgres password=gdg dbname=andnyang sslmode=disable")
	if error != nil {
		log.Print(error)
	}

	defer db.Close()

	for {
		line, err := ircbot.conn.ReadLine()
		if err != nil {
			break
		}

		arr := TokenizeLine(line)

		// Each line is started with ':' except 'PING' messages.
		if arr[0] == "PING" {
			token := arr[1]
			ircbot.Pong(token)
			continue
		} else if arr[0][0] != ':' {
			fmt.Printf("Something is wrong!\n")
			fmt.Printf(">>> %s\n", line)
			continue
		}

		// We will send join message after we will get the first notification. Otherwise, the message we sent are not processed by the server.
		systemMessageNo := arr[1]
		if systemMessageNo == "001" {
			for _, channel := range channels {
				request := fmt.Sprintf("JOIN %s", channel.channel)
				ircbot.conn.PrintfLine(request)
			}
			continue
		}

		command := arr[1]
		name := strings.Split(arr[0][1:], "!")[0]
		for _, channel := range channels {
			if command == "PRIVMSG" && arr[2] == channel.channel && arr[3][1:] == "!action" {
				channel.Activities()
			} else if command == "PRIVMSG" && arr[2] == channel.channel && arr[3][1] == '!' {
				fmt.Printf(">>> %s\n", line)
				channel.Talk(arr[3][2:])
			} else if command == "PRIVMSG" && arr[2] == channel.channel {
				channel.Log(db, name, arr[3][1:])
			} else if command == "JOIN" && arr[2][1:] == channel.channel {
				fmt.Printf(">>> %s\n", line)
				if name == ircbot.nick {
					channel.Talk("오랜만이에요. :) 모두 안녕하세요.")
				} else {
					text := fmt.Sprintf("안녕하세요. %s님 ^^", name)
					channel.Talk(text)
					channel.Op(name)
				}
			}
		}
		fmt.Printf(">>> %s\n", line)
	}
}

func TokenizeLine(line string) []string {
	// Each line will break into four strings. First three strings will be splitted by whitespace, all rest will be the fourth string.

	output := make([]string, 4)
	oi := 0
	n := 0
	space := false

	for i, c := range line {
		if oi == 3 {
			continue
		} else if space == false && c == ' ' && n < i {
			output[oi] = line[n:i]
			space = true
			n = i + 1
			oi++
		} else if space == true && c == ' ' {
			n = i + 1
			continue
		} else if space == true && c != ' ' {
			space = false
		}
	}

	output[oi] = line[n:]
	return output
}
