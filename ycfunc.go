package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/essentialkaos/ek/v13/log"
	"github.com/essentialkaos/ek/v13/req"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// VER is current function version
const VER = "0.0.1"

// TIMER_EVENT_TYPE is timer event type
const TIMER_EVENT_TYPE = "yandex.cloud.events.serverless.triggers.TimerMessage"

// ////////////////////////////////////////////////////////////////////////////////// //

type Trigger struct {
	Messages []*Message `json:"messages"`
}

type Message struct {
	Metadata *Metadata `json:"event_metadata"`
	Details  *Details  `json:"details"`
}

type Metadata struct {
	EventType string `json:"event_type"`
	CreatedAt string `json:"created_at"`
}

type Details struct {
	TriggerID string `json:"trigger_id"`
	Payload   string `json:"payload"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// main is used for check for compilation errors
func main() {
	return
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Timer is handler for timer trigger
func Timer(ctx context.Context, trigger *Trigger) error {
	req.SetUserAgent("YCFunction|funky", VER)

	log.Global.UseJSON = true
	log.Global.WithCaller = true

	defer log.Flush()

	log.Info("Got trigger event")

	if !validatePayload(trigger) {
		return fmt.Errorf("Error while trigger event validation")
	}

	var lf log.Fields

	invURL := trigger.Payload()
	async := os.Getenv("ASYNC") != ""

	r := req.Request{
		URL:         invURL,
		AutoDiscard: true,
		Headers: req.Headers{
			"X-Funky-Version": VER,
			"Trigger-Source":  "funky",
			"Trigger-Id":      trigger.ID(),
			"Trigger-Ts":      trigger.Timestamp(),
		},
	}

	if async {
		r.Query = req.Query{"integration": "async"}
	}

	lf.Add(log.F{"url", invURL}, log.F{"async", async})
	log.Info("Sending HTTP request", lf)

	resp, err := r.Get()

	if err != nil {
		log.Error("Error while sending request: %v", err, lf)
		return fmt.Errorf("Error while sending request")
	}

	lf.Add(log.F{"status-code", resp.StatusCode})

	if resp.StatusCode != req.STATUS_OK {
		log.Error("Server returned non-ok status code %d", lf)
		return fmt.Errorf("Error while processing response")
	}

	log.Info("Request successfully sent", lf)

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Payload returns payload from trigger data
func (d *Trigger) Payload() string {
	return d.Messages[0].Details.Payload
}

// ID returns trigger ID
func (d *Trigger) ID() string {
	return d.Messages[0].Details.TriggerID
}

// Timestamp returns trigger timestamp
func (d *Trigger) Timestamp() string {
	created, _ := time.Parse(time.RFC3339Nano, d.Messages[0].Metadata.CreatedAt)
	return strconv.FormatInt(created.UnixMilli(), 10)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// validatePayload validates trigger payload
func validatePayload(trigger *Trigger) bool {
	switch {
	case trigger == nil:
		log.Error("Trigger data is nil")
		return false

	case len(trigger.Messages) == 0:
		log.Error("No messages in trigger event")
		return false

	case trigger.Messages[0].Metadata == nil:
		log.Error("No metadata in message #0")
		return false

	case trigger.Messages[0].Metadata.EventType != TIMER_EVENT_TYPE:
		log.Error(
			"Unsupported event type",
			log.F{"event-type", trigger.Messages[0].Metadata.EventType},
		)
		return false

	case trigger.Messages[0].Details == nil:
		log.Error("No details in message #0")
		return false

	case trigger.Messages[0].Details.Payload == "":
		log.Error("Payload is empty")
		return false
	}

	_, err := url.Parse(trigger.Messages[0].Details.Payload)

	if err != nil {
		log.Error("Can't parse URL in payload: %v", err)
		return false
	}

	return true
}
