package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
)

type Bot struct {
	server        string
	port          string
	nick          string
	user          string
	channel       string
	pass          string
	pread, pwrite chan string
	conn          net.Conn
}

func NewBot() *Bot {

		//server:  "irc.ozinger.org",
	return &Bot{
		server: "kanade.irc.ozinger.org",
		port:    "6668",
		nick:    "gdgandbot",
		channel: "gdgand",
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
	//	fmt.Fprintf(conn, "USER %s 8 * :%s\n", ircbot.nick, ircbot.nick)
	fmt.Fprintf(conn, "USER bot 0 * :fishing\n")
	conn.Write([]byte("NICK " + ircbot.nick + "\n"))
	conn.Write([]byte("JOIN " + ircbot.channel + "\n"))
	defer conn.Close()

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	for {
		line, err := tp.ReadLine()
		if err != nil {
			break
		}
		fmt.Printf("%s\n", line)
	}
}
