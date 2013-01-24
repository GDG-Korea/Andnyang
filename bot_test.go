package main

import (
	"strings"
	"testing"
)

func TestTokenizeLine(t *testing.T) {
	testData := []struct {
		src, expect string
	}{
		{
			":name!user@host PRIVMSG #channel :hello world",
			":name!user@host,PRIVMSG,#channel,:hello world",
		}, {
			":name!user@host MODE #channel +o someone",
			":name!user@host,MODE,#channel,+o someone",
		}, {
			"PING 23234",
			"PING,23234,,",
		},
	}

	for _, s := range testData {
		result := strings.Join(TokenizeLine(s.src), ",")
		if result != s.expect {
			t.Errorf("Expected \"%\" but, got \"%s\"\n", result, s.expect)
		}
	}
}
