package slackboy

import (
	"reflect"
	"testing"
)

func TestPostToDefault(t *testing.T) {
	opt := Options{
		Env:         "production",
		WebhookURL:  "https://hooks.slack.com/services/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",
		DefaultTags: []string{"host: 127.0.0.1", "app: slackboy"},
	}
	slackBoy := New(opt)

	slackBoy.Success("Success 1", "Success description 1")
	slackBoy.Info("Info 1", "Info description 1")
}

func TestPostToDefault2(t *testing.T) {
	opt := Options{
		Env:         "production",
		WebhookURL:  "https://hooks.slack.com/services/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",
		DefaultTags: []string{"host: 127.0.0.1"},
	}
	slackBoy := New(opt)

	slackBoy.Success("Success 1", "Success description 1")
	slackBoy.Info("Info 1", "Info description 1")

	opt2 := Options{
		Env:         "production",
		WebhookURL:  "https://hooks.slack.com/services/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",
		DefaultTags: []string{"host: 127.0.0.1", "user: @ariefrahmansyah"},
	}
	slackBoy2 := New(opt2)

	slackBoy2.Warning("Warning 1", "Warning description 1")
	slackBoy2.Error("Error 1", "Error description 1")
}

func TestPostToSuccess(t *testing.T) {
	opt := Options{
		Env:            "production",
		WebhookURL:     "https://hooks.slack.com/services/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",
		SuccessChannel: "success",
		InfoChannel:    "info",
		WarningChannel: "warning",
		ErrorChannel:   "error",
		DefaultTags:    []string{"host: 127.0.0.1"},
	}
	slackBoy := New(opt)

	slackBoy.Success("Success 1", "Success description 1")
}

func TestPostToSuccessWithTags(t *testing.T) {
	opt := Options{
		Env:            "production",
		WebhookURL:     "https://hooks.slack.com/services/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",
		SuccessChannel: "success",
		InfoChannel:    "info",
		WarningChannel: "warning",
		ErrorChannel:   "error",
		DefaultTags:    []string{"host: 127.0.0.1"},
	}
	slackBoy := New(opt)

	slackBoy.Success("Success 1", "Success description 1", []string{"user: @ariefrahmansyah"}...)
}

func TestPostToInfo(t *testing.T) {
	opt := Options{
		Env:            "production",
		WebhookURL:     "https://hooks.slack.com/services/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",
		SuccessChannel: "success",
		InfoChannel:    "info",
		WarningChannel: "warning",
		ErrorChannel:   "error",
		DefaultTags:    []string{"host: 127.0.0.1"},
	}
	slackBoy := New(opt)

	slackBoy.Info("Info 1", "Info description 1")
}

func TestPostToInfoWithLink(t *testing.T) {
	opt := Options{
		Env:            "production",
		WebhookURL:     "https://hooks.slack.com/services/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",
		SuccessChannel: "success",
		InfoChannel:    "info",
		WarningChannel: "warning",
		ErrorChannel:   "error",
		DefaultTags:    []string{"host: 127.0.0.1"},
	}
	slackBoy := New(opt)

	slackBoy.Info("Info 1", "Link to Google: https://www.google.com")
}

func TestPostToWarning(t *testing.T) {
	opt := Options{
		Env:            "production",
		WebhookURL:     "https://hooks.slack.com/services/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",
		SuccessChannel: "success",
		InfoChannel:    "info",
		WarningChannel: "warning",
		ErrorChannel:   "error",
	}
	slackBoy := New(opt)

	slackBoy.Warning("Warning 1", "Warning description 1")
}

func TestPostToError(t *testing.T) {
	opt := Options{
		Env:            "production",
		WebhookURL:     "https://hooks.slack.com/services/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",
		SuccessChannel: "success",
		InfoChannel:    "info",
		WarningChannel: "warning",
		ErrorChannel:   "error",
	}
	slackBoy := New(opt)

	slackBoy.Error("Error 1", "Error description 1")
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

func TestSlackBoy_GetTags(t *testing.T) {
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
			if got := s.GetTags(tt.args.msg); got != tt.want {
				t.Errorf("SlackBoy.GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
