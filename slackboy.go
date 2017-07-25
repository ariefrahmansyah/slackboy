package slackboy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

type Message struct {
	Channel         string
	Text            string
	AttachmentColor string
	AttachmentText  string
}

type messageMap map[int]*Message

type SlackboyTags map[string]string

type Options struct {
	Env  string
	Tags SlackboyTags

	WebhookURL     string
	DefaultChannel string
	SuccessChannel string
	InfoChannel    string
	WarningChannel string
	ErrorChannel   string
}

type SlackBoy struct {
	message messageMap
	opt     Options
}

func New(o Options) *SlackBoy {
	msgMap := messageMap{}

	msgMap[successType] = &Message{Channel: o.SuccessChannel, AttachmentColor: "good"}
	msgMap[infoType] = &Message{Channel: o.InfoChannel, AttachmentColor: "#3AA3E3"}
	msgMap[warningType] = &Message{Channel: o.WarningChannel, AttachmentColor: "warning"}
	msgMap[errorType] = &Message{Channel: o.ErrorChannel, AttachmentColor: "danger"}

	return &SlackBoy{message: msgMap, opt: o}
}

func (s *SlackBoy) getMessage(msgType int) *Message {
	if msg, ok := s.message[msgType]; ok {
		return msg
	}

	return &Message{}
}

func (s *SlackBoy) Success(text, snip string) {
	msg := s.getMessage(successType)
	msg.Text = text
	msg.AttachmentText = snip

	s.Post(msg)
}

func (s *SlackBoy) Info(text, snip string) {
	msg := s.getMessage(infoType)
	msg.Text = text
	msg.AttachmentText = snip

	s.Post(msg)
}

func (s *SlackBoy) Warning(text, snip string) {
	msg := s.getMessage(warningType)
	msg.Text = text
	msg.AttachmentText = snip

	s.Post(msg)
}

func (s *SlackBoy) Error(text, snip string) {
	msg := s.getMessage(errorType)
	msg.Text = text
	msg.AttachmentText = snip

	s.Post(msg)
}

func (s *SlackBoy) Post(msg *Message) {
	channel := msg.Channel
	if channel == "" {
		channel = s.opt.DefaultChannel
	}

	tags := []string{fmt.Sprintf("`env: %s`", s.opt.Env)}
	if len(s.opt.Tags) > 0 {
		for k, v := range s.opt.Tags {
			tags = append(tags, fmt.Sprintf("`%s: %s`", k, v))
		}
	}
	tagString := strings.Join(tags, " ")

	payload := map[string]interface{}{
		"channel":    channel,
		"link_names": 1,
		"attachments": []map[string]interface{}{
			map[string]interface{}{
				"title":     msg.Text,
				"color":     msg.AttachmentColor,
				"text":      msg.AttachmentText + "\n" + tagString,
				"mrkdwn_in": []string{"text"},
			},
		},
	}

	go s.post(payload)
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
