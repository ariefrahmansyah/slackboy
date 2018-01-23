# slackboy [![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/ariefrahmansyah/slackboy) [![CircleCI](https://circleci.com/gh/ariefrahmansyah/slackboy/tree/master.png?style=shield)](https://circleci.com/gh/ariefrahmansyah/slackboy/tree/master) [![Coverage Status](https://coveralls.io/repos/github/ariefrahmansyah/slackboy/badge.svg?branch=master)](https://coveralls.io/github/ariefrahmansyah/slackboy?branch=master) [![GoReportCard](https://goreportcard.com/badge/github.com/ariefrahmansyah/slackboy)](https://goreportcard.com/report/github.com/ariefrahmansyah/slackboy)

SlackBoy will help you send Slack webhooks message to specify channel depending on message type (success, info, warning, or error).

## Usage
```go
import "github.com/ariefrahmansyah/slackboy"

slackboyOpt := slackboy.Options{
        Env:         opt.Env,
        DefaultTags: opt.DefaultTags,
        WebhookURL:  opt.WebhookURL,
}
slackBoy := slackboy.New(slackboyOpt)

slackBoy.Success("Success 1", "Success description 1")
```

## Synchronous or Asynchronous
As default, SlackBoy will always send webhook asynchronously using go routines. To disable it, set option to `Synchronous` to `true`.
