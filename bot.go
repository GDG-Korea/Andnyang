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
	server        string
	port          string
	nick          string
	user          string
	channel       string
	pass          string
	conn          net.Conn
}

func NewBot() *Bot {

	//server:  "irc.ozinger.org",
	return &Bot{
		server:  "kanade.irc.ozinger.org",
		port:    "6668",
		nick:    "안드봇",
		channel: "#gdgand",
		pass:    "",
		conn:    nil,
		user:    "gdgand",
	}
}

func (bot *Bot) Connect() (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", bot.server+":"+bot.port)
	if err != nil {
		log.Fatal("unable to connect to IRC server ", err)
	}
	bot.conn = conn
	log.Printf("Connected to IRC server %s (%s)\n", bot.server, bot.conn.RemoteAddr())
	return bot.conn, nil
}

func main() {
	ircbot := NewBot()
	conn, _ := ircbot.Connect()
	defer conn.Close()

	reader := bufio.NewReader(conn)
	tpReader := textproto.NewReader(reader)
	writer := bufio.NewWriter(conn)
	tpWriter := textproto.NewWriter(writer)

	userCommand := fmt.Sprintf("USER %s 8 * :%s\n", ircbot.user, ircbot.user)
	tpWriter.PrintfLine(userCommand)
	tpWriter.PrintfLine("NICK " + ircbot.nick)

	for {
		line, err := tpReader.ReadLine()
		if err != nil {
			break
		}

		arr := strings.Split(line, " ")
		if arr[0] == "PING" {
			token := arr[1]
			request := fmt.Sprintf("PONG %s", token)
			tpWriter.PrintfLine(request)
		} else if arr[0][0] == ':' && arr[1] == "001" {
			request := fmt.Sprintf("JOIN %s", ircbot.channel)
			tpWriter.PrintfLine(request)
		} else if arr[0][0] == ':' && arr[1] == "PRIVMSG" && arr[2] == ircbot.channel && arr[3][1] == '!' {
			fmt.Printf(">>> %s\n", line)
			request := fmt.Sprintf("PRIVMSG %s :%s", ircbot.channel, arr[3][2:])
			tpWriter.PrintfLine(request)
		} else {
			fmt.Printf(">>> %s\n", line)
		}
	}
}
