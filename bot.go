package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
)

type Bot struct {
	server   string
	port     string
	nick     string
	user     string
	pass     string
	conn     net.Conn
	tpReader *textproto.Reader
	tpWriter *textproto.Writer
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
	conn, err := net.Dial("tcp", bot.server+":"+bot.port)
	if err != nil {
		log.Fatal("unable to connect to IRC server", err)
	}
	bot.conn = conn

	log.Printf("Connected to IRC server %s(%s)\n", bot.server, bot.conn.RemoteAddr())
	reader := bufio.NewReader(bot.conn)
	writer := bufio.NewWriter(bot.conn)
	bot.tpReader = textproto.NewReader(reader)
	bot.tpWriter = textproto.NewWriter(writer)

	userCommand := fmt.Sprintf("USER %s 8 * :%s\n", bot.user, bot.user)
	bot.tpWriter.PrintfLine(userCommand)
	bot.tpWriter.PrintfLine("NICK " + bot.nick)
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
	c.bot.tpWriter.PrintfLine(text)
}

func main() {
	ircbot := NewBot()
	ircbot.Connect()
	defer ircbot.Close()
	channel := ircbot.NewChannel("#gdgand")

	for {
		line, err := ircbot.tpReader.ReadLine()
		if err != nil {
			break
		}

		arr := TokenizeLine(line)

		// Each line is started with ':' except 'PING' messages.
		if arr[0] == "PING" {
			token := arr[1]
			request := fmt.Sprintf("PONG %s", token)
			ircbot.tpWriter.PrintfLine(request)
			continue
		} else if arr[0][0] != ':' {
			fmt.Printf("Something is wrong!\n")
			fmt.Printf(">>> %s\n", line)
			continue
		}

		// We will send join message after we will get the first notification. Otherwise, the message we sent are not processed by the server.
		systemMessageNo := arr[1]
		if systemMessageNo == "001" {
			request := fmt.Sprintf("JOIN %s", channel.channel)
			ircbot.tpWriter.PrintfLine(request)
			continue
		}

		command := arr[1]
		name := strings.Split(arr[0][1:], "!")[0]
		if command == "PRIVMSG" && arr[2] == channel.channel && arr[3][1] == '!' {
			fmt.Printf(">>> %s\n", line)
			channel.Talk(arr[3][2:])
		} else if command == "JOIN" && arr[2][1:] == channel.channel {
			fmt.Printf(">>> %s\n", line)
			if name == ircbot.nick {
				channel.Talk("오랜만이에요. :) 모두 안녕하세요.")
			} else {
				text := fmt.Sprintf("안녕하세요. %s님 ^^", name)
				channel.Talk(text)
			}
			request := fmt.Sprintf("MODE %s +o %s", channel.channel, name)
			ircbot.tpWriter.PrintfLine(request)
		} else {
			fmt.Printf(">>> %s\n", line)
		}
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
