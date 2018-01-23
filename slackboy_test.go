package slackboy

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

var webhookSuccess *httptest.Server
var webhookFail *httptest.Server

func testMain(m *testing.M) int {
	webhookSuccess = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer webhookSuccess.Close()

	webhookFail = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("not ok"))
	}))
	defer webhookFail.Close()

	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func TestNew(t *testing.T) {
	type args struct {
		o Options
	}
	tests := []struct {
		name string
		args args
		want *SlackBoy
	}{
		{
			"Empty options",
			args{
				Options{},
			},
			&SlackBoy{
				message: messageMap{
					successType: &Message{AttachmentColor: "good"},
					infoType:    &Message{AttachmentColor: "#3AA3E3"},
					warningType: &Message{AttachmentColor: "warning"},
					errorType:   &Message{AttachmentColor: "danger"},
				},
				opt: Options{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlackBoy_getMessageType(t *testing.T) {
	type fields struct {
		message messageMap
		opt     Options
	}
	type args struct {
		msgType int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Message
	}{
		{
			"message type exist",
			fields{
				message: messageMap{
					successType: &Message{AttachmentColor: "good"},
				},
			},
			args{
				successType,
			},
			&Message{AttachmentColor: "good"},
		},
		{
			"message type not exist",
			fields{
				message: messageMap{},
			},
			args{
				successType,
			},
			&Message{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SlackBoy{
				message: tt.fields.message,
				opt:     tt.fields.opt,
			}
			if got := s.getMessageType(tt.args.msgType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlackBoy.getMessageType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostDefault(t *testing.T) {
	opt := Options{
		Env:         "production",
		WebhookURL:  webhookSuccess.URL,
		DefaultTags: []string{"host: 127.0.0.1", "app: slackboy"},
	}
	slackBoy := New(opt)

	slackBoy.Success("Success 1", "Success description 1")
	slackBoy.Info("Info 1", "Info description 1")
	slackBoy.Warning("Warning 1", "Warning description 1")
	slackBoy.Error("Error 1", "Error description 1")
}

func TestPostSuccessWithTags(t *testing.T) {
	opt := Options{
		Env:            "production",
		WebhookURL:     webhookSuccess.URL,
		SuccessChannel: "success",
		InfoChannel:    "info",
		WarningChannel: "warning",
		ErrorChannel:   "error",
		DefaultTags:    []string{"host: 127.0.0.1"},
	}
	slackBoy := New(opt)

	slackBoy.Success("Success 1", "Success description 1", []string{"user: @ariefrahmansyah"}...)
}

func TestPostInfoWithLink(t *testing.T) {
	opt := Options{
		Env:            "production",
		WebhookURL:     webhookSuccess.URL,
		SuccessChannel: "success",
		InfoChannel:    "info",
		WarningChannel: "warning",
		ErrorChannel:   "error",
		DefaultTags:    []string{"host: 127.0.0.1"},
	}
	slackBoy := New(opt)

	slackBoy.Info("Info 1", "Link to Google: https://www.google.com")
}

func TestPostFail(t *testing.T) {
	opt := Options{
		Env:         "production",
		WebhookURL:  webhookFail.URL,
		DefaultTags: []string{"host: 127.0.0.1", "app: slackboy"},
	}
	slackBoy := New(opt)

	slackBoy.Success("Success 1", "Success description 1")
}

func TestPostAsync(t *testing.T) {
	opt := Options{
		Env:         "production",
		WebhookURL:  webhookSuccess.URL,
		Synchronous: true,
		DefaultTags: []string{"host: 127.0.0.1", "app: slackboy"},
	}
	slackBoy := New(opt)

	slackBoy.Success("Success 1", "Success description 1")

	time.Sleep(1 * time.Second)
}

func TestSlackBoy_getTags(t *testing.T) {
	type fields struct {
		message messageMap
		opt     Options
	}
	type args struct {
		msg Message
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"empty",
			fields{
				message: messageMap{},
				opt: Options{
					Env:         "",
					DefaultTags: []string{},
				},
			},
			args{
				msg: Message{
					Tags: []string{},
				},
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SlackBoy{
				message: tt.fields.message,
				opt:     tt.fields.opt,
			}
			if got := s.getTags(tt.args.msg); got != tt.want {
				t.Errorf("SlackBoy.getTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortTags(t *testing.T) {
	type args struct {
		tags []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"1",
			args{
				[]string{"env: production", "app: slackboy"},
			},
			[]string{"app: slackboy", "env: production"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SortTags(tt.args.tags)

			if got := tt.args.tags; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
