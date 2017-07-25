package slackboy

import "testing"

func TestPostToDefault(t *testing.T) {
	opt := Options{
		Env:         "production",
		WebhookURL:  "https://hooks.slack.com/services/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",
		DefaultTags: []string{"host: 127.0.0.1"},
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
