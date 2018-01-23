package slackboy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

var httpClient = http.Client{
	Timeout: 5 * time.Second,
}

const (
	successType = iota
	infoType
	warningType
	errorType
)

// Message is the webhook message contains channel data, text, snippet, tags and attachment color.
type Message struct {
	Channel         string
	Text            string
	Snippet         string
	AttachmentColor string
	Tags            []string
}

type messageMap map[int]*Message

// Options to create new slackboy object.
type Options struct {
	Env         string
	DefaultTags []string
	Synchronous bool

	WebhookURL     string
	DefaultChannel string
	SuccessChannel string
	InfoChannel    string
	WarningChannel string
	ErrorChannel   string
}

// SlackBoy object.
type SlackBoy struct {
	message messageMap
	opt     Options
}

// New returns new *SlackBoy object.
func New(o Options) *SlackBoy {
	msgMap := messageMap{}

	msgMap[successType] = &Message{Channel: o.SuccessChannel, AttachmentColor: "good"}
	msgMap[infoType] = &Message{Channel: o.InfoChannel, AttachmentColor: "#3AA3E3"}
	msgMap[warningType] = &Message{Channel: o.WarningChannel, AttachmentColor: "warning"}
	msgMap[errorType] = &Message{Channel: o.ErrorChannel, AttachmentColor: "danger"}

	return &SlackBoy{message: msgMap, opt: o}
}

func (s *SlackBoy) getMessageType(msgType int) *Message {
	if msg, ok := s.message[msgType]; ok {
		return msg
	}

	return &Message{}
}

// Success sends webhook with green attachment color.
func (s *SlackBoy) Success(text, snip string, tags ...string) {
	msg := s.getMessageType(successType)
	msg.Text = text
	msg.Snippet = snip
	msg.Tags = tags

	s.Post(msg)
}

// Info sends webhook with blue attachment color.
func (s *SlackBoy) Info(text, snip string) {
	msg := s.getMessageType(infoType)
	msg.Text = text
	msg.Snippet = snip

	s.Post(msg)
}

// Warning sends webhook with orange attachment color.
func (s *SlackBoy) Warning(text, snip string) {
	msg := s.getMessageType(warningType)
	msg.Text = text
	msg.Snippet = snip

	s.Post(msg)
}

// Error sends webhook with red attachment color.
func (s *SlackBoy) Error(text, snip string) {
	msg := s.getMessageType(errorType)
	msg.Text = text
	msg.Snippet = snip

	s.Post(msg)
}

// Post posts specified message.
func (s *SlackBoy) Post(msg *Message) {
	channel := msg.Channel
	if channel == "" {
		channel = s.opt.DefaultChannel
	}

	tagsString := s.getTags(*msg)

	payload := map[string]interface{}{
		"channel":    channel,
		"link_names": 1,
		"attachments": []map[string]interface{}{
			map[string]interface{}{
				"title":     msg.Text,
				"color":     msg.AttachmentColor,
				"text":      msg.Snippet + "\n" + tagsString,
				"mrkdwn_in": []string{"text"},
			},
		},
	}

	if !s.opt.Synchronous {
		go s.post(payload)
	} else {
		s.post(payload)
	}
}

func (s *SlackBoy) post(payload map[string]interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := httpClient.Post(s.opt.WebhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(data))
	}

	return
}

func (s *SlackBoy) getTags(msg Message) string {
	tags := []string{}
	tagsString := ""

	if s.opt.Env != "" {
		tags = append(tags, fmt.Sprintf("`env: %s`", s.opt.Env))
	}

	if len(s.opt.DefaultTags) > 0 {
		for _, v := range s.opt.DefaultTags {
			tags = append(tags, fmt.Sprintf("`%s`", v))
		}
	}

	if len(msg.Tags) > 0 {
		for _, v := range msg.Tags {
			tags = append(tags, fmt.Sprintf("`%s`", v))
		}
	}

	if len(tags) > 0 {
		SortTags(tags)
		tagsString = strings.Join(tags, " ")
	}

	return tagsString
}

type tagsSorter []string

func (tags tagsSorter) Len() int           { return len(tags) }
func (tags tagsSorter) Swap(i, j int)      { tags[i], tags[j] = tags[j], tags[i] }
func (tags tagsSorter) Less(i, j int) bool { return tags[i] < tags[j] }

// SortTags sorts payment tags by total price descending
func SortTags(tags []string) {
	sort.Sort(tagsSorter(tags))
}
