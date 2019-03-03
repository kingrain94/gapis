// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	calendar "google.golang.org/api/calendar/v3"
)

func init() {
	registerDemo("calendar", calendar.CalendarScope, calendarMain)
}

// calendarMain is an example that demonstrates calling the Calendar API.
// Its purpose is to test out the ability to get maps of struct objects.
//
// Example usage:
//   go build -o go-api-demo *.go
//   go-api-demo -clientid="my-clientid" -secret="my-secret" calendar
func calendarMain(client *http.Client, argv []string) {
	if len(argv) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: calendar")
		return
	}

	svc, err := calendar.New(client)
	if err != nil {
		log.Fatalf("unable to create calendar service: %v", err)
	}

	fmt.Println()
	switch argv[0] {
	case "--show-owner":
		fmt.Println("*********** Calendar Details ***********")
		err = showOwnerDetails(svc)
		if err != nil {
			log.Fatal("failed to show owner details")
		}
		fmt.Println("****************************************")

	case "--show-events":
		fmt.Println("*********** All Events ***********")
		if len(argv) < 2 {
			log.Fatal("pass calendar id behind!")
		}
		err = showEvents(svc, argv[1])
		if err != nil {
			log.Fatal(errors.Wrap(err, errorf("failed to show events for '%s'", argv[1])).Error())
		}
		fmt.Println("***********************************")

	case "--show-daily-updated-events":
		fmt.Println("*********** Daily Updated Events ***********")
		if len(argv) < 2 {
			log.Fatal("pass calendar id behind!")
		}
		err = showDailyUpdatedEvents(svc, argv[1])
		if err != nil {
			log.Fatal(errors.Wrap(err, errorf("failed to show events for '%s'", argv[1])).Error())
		}
		fmt.Println("********************************************")
	}
	fmt.Println()
}

func showOwnerDetails(svc *calendar.Service) error {
	c, err := svc.Colors.Get().Do()
	if err != nil {
		return errors.Wrap(err, errorf("unable to retrieve calendar colors"))
	}

	fmt.Printf("kind of colors: %v", c.Kind)
	fmt.Printf("colors last updated: %v", c.Updated)

	for k, v := range c.Calendar {
		fmt.Printf("calendar[%v]: background=%v, foreground=%v", k, v.Background, v.Foreground)
	}

	for k, v := range c.Event {
		fmt.Printf("event[%v]: background=%v, foreground=%v", k, v.Background, v.Foreground)
	}

	listRes, err := svc.CalendarList.List().Fields("items/id").Do()
	if err != nil {
		return errors.Wrap(err, errorf("unable to retrieve list of calendars"))
	}
	for _, v := range listRes.Items {
		fmt.Printf("calendar id: %v\n", v.Id)
	}

	return nil
}

func showEvents(svc *calendar.Service, calID string) error {
	res, err := svc.Events.List(calID).Fields("items(updated,summary)", "summary", "nextPageToken").Do()
	if err != nil {
		return errors.Wrap(err, errorf("unable to retrieve calendar events list"))
	}
	for _, v := range res.Items {
		fmt.Printf("calendar id %q event: %v: %q\n", calID, v.Updated, v.Summary)
	}
	fmt.Printf("calendar id %q Summary: %v\n", calID, res.Summary)
	fmt.Printf("calendar id %q next page token: %v\n", calID, res.NextPageToken)
	return nil
}

func showDailyUpdatedEvents(svc *calendar.Service, calID string) error {
	yesterday := time.Now().AddDate(0, 0, -1)
	res, err := svc.Events.List(calID).TimeMin(yesterday.Format(time.RFC3339)).Fields("items(id,updated,summary,visibility)", "summary").Do()
	if err != nil {
		return errors.Wrap(err, errorf("unable to retrieve calendar events list"))
	}

	if len(res.Items) == 0 {
		fmt.Printf("no daily updated events\n")
		return nil
	} else {
		fmt.Printf("daily updated events of '%s'\n", calID)
	}

	for _, v := range res.Items {
		updated, err := time.Parse(time.RFC3339, v.Updated)
		if err != nil {
			return errors.Wrap(err, errorf("failed to convert time to RFC3339"))
		}
		if updated.After(yesterday) {
			summary := v.Summary
			if summary == "" && v.Visibility == "private" {
				summary = "busy - private"
			}
			fmt.Printf("-------------\nevent id: %s\nsummary: %s\nupdated: %v\n", v.Id, summary, v.Updated)
		}
	}
	fmt.Printf("-------------\n")
	return nil
}
