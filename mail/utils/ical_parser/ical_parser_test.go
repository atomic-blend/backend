package icalparser

import (
	"testing"

	ical "github.com/emersion/go-ical"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestParseICal(t *testing.T) {
	icalData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Example Corp//Example//EN
BEGIN:VEVENT
UID:1234567890@example.com
DTSTAMP:20231202T120000Z
DTSTART:20231202T120000Z
DTEND:20231202T130000Z
SUMMARY:Test Event
DESCRIPTION:This is a test event
LOCATION:Test Location
END:VEVENT
END:VCALENDAR`

	calendars, err := ParseICal([]byte(icalData))
	if err != nil {
		t.Fatalf("ParseICal failed: %v", err)
	}

	if len(calendars) != 1 {
		t.Fatalf("Expected 1 calendar, got %d", len(calendars))
	}

	cal := calendars[0]
	events := cal.Events()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	event := events[0]
	summary, _ := event.Props.Text(ical.PropSummary)
	if summary != "Test Event" {
		t.Errorf("Expected summary 'Test Event', got '%s'", summary)
	}
}

func TestToCalendarPayload(t *testing.T) {
	// Create a mock ical.Calendar
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//Example//EN")

	event := ical.NewEvent()
	event.Props.SetText(ical.PropUID, "test-uid")
	event.Props.SetDateTime(ical.PropDateTimeStamp, timestamppb.Now().AsTime())
	event.Props.SetDateTime(ical.PropDateTimeStart, timestamppb.Now().AsTime())
	event.Props.SetText(ical.PropSummary, "Test Event")

	cal.Children = append(cal.Children, event.Component)

	payload, err := ToCalendarPayload(*cal)
	if err != nil {
		t.Fatalf("ToCalendarPayload failed: %v", err)
	}

	if payload.Version != "2.0" {
		t.Errorf("Expected version '2.0', got '%s'", payload.Version)
	}

	if payload.ProdId != "-//Example//EN" {
		t.Errorf("Expected prodId '-//Example//EN', got '%s'", payload.ProdId)
	}

	if len(payload.Events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(payload.Events))
	}

	evt := payload.Events[0]
	if evt.Uid != "test-uid" {
		t.Errorf("Expected UID 'test-uid', got '%s'", evt.Uid)
	}

	if evt.Summary != "Test Event" {
		t.Errorf("Expected summary 'Test Event', got '%s'", evt.Summary)
	}
}

func TestParseICalMultipleCalendars(t *testing.T) {
	icalData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Example//EN
BEGIN:VEVENT
UID:event1
SUMMARY:Event 1
END:VEVENT
END:VCALENDAR
BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Another//EN
BEGIN:VEVENT
UID:event2
SUMMARY:Event 2
END:VEVENT
END:VCALENDAR`

	calendars, err := ParseICal([]byte(icalData))
	if err != nil {
		t.Fatalf("ParseICal failed: %v", err)
	}

	if len(calendars) != 2 {
		t.Fatalf("Expected 2 calendars, got %d", len(calendars))
	}

	if calendars[0].Events()[0].Props.Get(ical.PropUID).Value != "event1" {
		t.Errorf("Expected first event UID 'event1'")
	}

	if calendars[1].Events()[0].Props.Get(ical.PropUID).Value != "event2" {
		t.Errorf("Expected second event UID 'event2'")
	}
}

func TestParseICalInvalidData(t *testing.T) {
	invalidData := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Example//EN
BEGIN:VEVENT
UID:event1
SUMMARY:Event 1
END:VEVENT
END:VCALENDAR
INVALID DATA HERE`

	_, err := ParseICal([]byte(invalidData))
	if err == nil {
		t.Fatalf("Expected ParseICal to fail on invalid data")
	}
}

func TestToCalendarPayloadWithAttendees(t *testing.T) {
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")

	event := ical.NewEvent()
	event.Props.SetText(ical.PropUID, "attendee-test")
	event.Props.SetDateTime(ical.PropDateTimeStart, timestamppb.Now().AsTime())

	attendee := ical.NewProp(ical.PropAttendee)
	attendee.Value = "mailto:test@example.com"
	attendee.Params.Set(ical.ParamCommonName, "Test User")
	attendee.Params.Set(ical.ParamParticipationStatus, "ACCEPTED")
	event.Props.Add(attendee)

	cal.Children = append(cal.Children, event.Component)

	payload, err := ToCalendarPayload(*cal)
	if err != nil {
		t.Fatalf("ToCalendarPayload failed: %v", err)
	}

	if len(payload.Events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(payload.Events))
	}

	evt := payload.Events[0]
	if len(evt.Attendees) != 1 {
		t.Fatalf("Expected 1 attendee, got %d", len(evt.Attendees))
	}

	attendeePayload := evt.Attendees[0]
	if attendeePayload.Email != "mailto:test@example.com" {
		t.Errorf("Expected email 'mailto:test@example.com', got '%s'", attendeePayload.Email)
	}

	if attendeePayload.Cn != "Test User" {
		t.Errorf("Expected CN 'Test User', got '%s'", attendeePayload.Cn)
	}

	if attendeePayload.PartStat != "ACCEPTED" {
		t.Errorf("Expected partstat 'ACCEPTED', got '%s'", attendeePayload.PartStat)
	}
}

func TestToCalendarPayloadWithTimezone(t *testing.T) {
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")

	tz := ical.NewComponent(ical.CompTimezone)
	tz.Props.SetText(ical.PropTimezoneID, "America/New_York")
	tz.Props.SetText(ical.PropTimezoneOffsetTo, "+0500")

	cal.Children = append(cal.Children, tz)

	payload, err := ToCalendarPayload(*cal)
	if err != nil {
		t.Fatalf("ToCalendarPayload failed: %v", err)
	}

	if len(payload.Timezones) != 1 {
		t.Fatalf("Expected 1 timezone, got %d", len(payload.Timezones))
	}

	tzPayload := payload.Timezones[0]
	if tzPayload.Tzid != "America/New_York" {
		t.Errorf("Expected TZID 'America/New_York', got '%s'", tzPayload.Tzid)
	}

	if tzPayload.Offset != 18000 { // +0500 = 5*3600 seconds
		t.Errorf("Expected offset 18000, got %d", tzPayload.Offset)
	}
}

func TestToCalendarPayloadEmptyCalendar(t *testing.T) {
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")

	payload, err := ToCalendarPayload(*cal)
	if err != nil {
		t.Fatalf("ToCalendarPayload failed: %v", err)
	}

	if payload.Version != "2.0" {
		t.Errorf("Expected version '2.0', got '%s'", payload.Version)
	}

	if len(payload.Events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(payload.Events))
	}

	if len(payload.Timezones) != 0 {
		t.Errorf("Expected 0 timezones, got %d", len(payload.Timezones))
	}
}
