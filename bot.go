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
	channel  string
	pass     string
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	tpReader *textproto.Reader
	tpWriter *textproto.Writer
}

func NewBot() *Bot {
	return &Bot{
		server:  "irc.ozinger.org",
		port:    "6668",
		nick:    "안드냥2",
		channel: "#gdgand",
		pass:    "",
		conn:    nil,
		user:    "gdgand",
	}
}

func (bot *Bot) Connect() (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", bot.server+":"+bot.port)
	if err != nil {
		log.Fatal("unable to connect to IRC server", err)
	}
	bot.conn = conn
	log.Printf("Connected to IRC server %s(%s)\n", bot.server, bot.conn.RemoteAddr())
	bot.reader = bufio.NewReader(bot.conn)
	bot.tpReader = textproto.NewReader(bot.reader)
	bot.writer = bufio.NewWriter(bot.conn)
	bot.tpWriter = textproto.NewWriter(bot.writer)

	return bot.conn, nil
}

func ParseLine(line string) []string {
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

type Channel struct {
	writer  *textproto.Writer
	channel string
}

func (b *Bot) NewChannel() *Channel {
	return &Channel {
		b.tpWriter,
		b.channel,
	}
}

func (c *Channel) Talk(msg string) {
	text := fmt.Sprintf("PRIVMSG %s :%s", c.channel, msg)
	c.writer.PrintfLine(text)
}

func main() {
	ircbot := NewBot()
	conn, _ := ircbot.Connect()
	defer conn.Close()

	//channel := Channel{ircbot.tpWriter, ircbot.channel}
	channel := ircbot.NewChannel()

	userCommand := fmt.Sprintf("USER %s 8 * :%s\n", ircbot.user, ircbot.user)
	ircbot.tpWriter.PrintfLine(userCommand)
	ircbot.tpWriter.PrintfLine("NICK " + ircbot.nick)

	for {
		line, err := ircbot.tpReader.ReadLine()
		if err != nil {
			break
		}

		arr := ParseLine(line)

		if arr[0] == "PING" {
			token := arr[1]
			request := fmt.Sprintf("PONG %s", token)
			ircbot.tpWriter.PrintfLine(request)
		} else if arr[0][0] == ':' && arr[1] == "001" {
			request := fmt.Sprintf("JOIN %s", ircbot.channel)
			ircbot.tpWriter.PrintfLine(request)
		} else if arr[0][0] == ':' && arr[1] == "PRIVMSG" && arr[2] == ircbot.channel && arr[3][1] == '!' {
			fmt.Printf(">>> %s\n", line)
			channel.Talk(arr[3][2:])
		} else if arr[0][0] == ':' && arr[1] == "JOIN" && arr[2][1:] == ircbot.channel {
			fmt.Printf(">>> %s\n", line)
			nameLine := arr[0][1:]
			nameArr := strings.Split(nameLine, "!")
			name := nameArr[0]
			if name == ircbot.nick {
				channel.Talk("오랜만이에요. :) 모두 안녕하세요.")
			} else {
				text := fmt.Sprintf("안녕하세요. %s님 ^^", name)
				channel.Talk(text)
			}
			request := fmt.Sprintf("MODE %s +o %s", ircbot.channel, name)
			ircbot.tpWriter.PrintfLine(request)
		} else {
			fmt.Printf(">>> %s\n", line)
		}
	}
}
