package icalparser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	calendarv1 "github.com/atomic-blend/backend/grpc/gen/calendar/v1"
	ical "github.com/emersion/go-ical"
	"github.com/teambition/rrule-go"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ParseICal parses iCalendar data and returns a slice of ical.Calendar
func ParseICal(data []byte) ([]ical.Calendar, error) {
	//convert to io.Reader
	reader := bytes.NewReader(data)

	calendars := []ical.Calendar{}

	dec := ical.NewDecoder(reader)
	for {
		cal, err := dec.Decode()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		calendars = append(calendars, *cal)
	}

	return calendars, nil
}

// ToCalendarPayload converts an ical.Calendar to a calendarv1.Calendar payload
func ToCalendarPayload(cal ical.Calendar) (*calendarv1.Calendar, error) {
	payload := &calendarv1.Calendar{}

	// Extract calendar properties
	if prodID, err := cal.Props.Text(ical.PropProductID); err == nil {
		payload.ProdId = prodID
	}
	if version, err := cal.Props.Text(ical.PropVersion); err == nil {
		payload.Version = version
	}
	if method, err := cal.Props.Text(ical.PropMethod); err == nil {
		payload.Method = method
	}

	// Extract events
	events := cal.Events()
	payload.Events = make([]*calendarv1.Event, len(events))
	for i, event := range events {
		evt, err := toCalendarEvent(event)
		if err != nil {
			return nil, err
		}
		payload.Events[i] = evt
	}

	// Extract timezones
	timezones := make([]*calendarv1.VTimezone, 0)
	for _, child := range cal.Children {
		if child.Name == ical.CompTimezone {
			tz, err := toVTimezone(child)
			if err != nil {
				return nil, err
			}
			timezones = append(timezones, tz)
		}
	}
	payload.Timezones = timezones

	return payload, nil
}

func toCalendarEvent(event ical.Event) (*calendarv1.Event, error) {
	evt := &calendarv1.Event{}

	if uid, err := event.Props.Text(ical.PropUID); err == nil {
		evt.Uid = uid
	}

	if dtStamp, err := event.Props.DateTime(ical.PropDateTimeStamp, nil); err == nil {
		evt.DtStamp = timestamppb.New(dtStamp)
	}

	if start, err := event.Props.DateTime(ical.PropDateTimeStart, nil); err == nil {
		evt.Start = toTimeSpec(start, event.Props.Get(ical.PropDateTimeStart))
	}

	if end, err := event.Props.DateTime(ical.PropDateTimeEnd, nil); err == nil {
		evt.End = toTimeSpec(end, event.Props.Get(ical.PropDateTimeEnd))
	}

	if durationProp := event.Props.Get(ical.PropDuration); durationProp != nil {
		if duration, err := durationProp.Duration(); err == nil {
			evt.Duration = durationpb.New(duration)
		}
	}

	if summary, err := event.Props.Text(ical.PropSummary); err == nil {
		evt.Summary = summary
	}

	if description, err := event.Props.Text(ical.PropDescription); err == nil {
		evt.Description = description
	}

	if location, err := event.Props.Text(ical.PropLocation); err == nil {
		evt.Location = location
	}

	if organizer := event.Props.Get(ical.PropOrganizer); organizer != nil {
		evt.Organizer = toCalAddress(organizer)
	}

	attendees := event.Props.Values(ical.PropAttendee)
	evt.Attendees = make([]*calendarv1.CalAddress, len(attendees))
	for i, attendee := range attendees {
		evt.Attendees[i] = toCalAddress(&attendee)
	}

	if rrule, err := event.Props.RecurrenceRule(); err == nil && rrule != nil {
		evt.Recurrence = toRRule(rrule)
	}

	exDates := event.Props.Values(ical.PropExceptionDates)
	evt.ExDates = make([]*calendarv1.TimeSpec, len(exDates))
	for i, exDate := range exDates {
		if dt, err := exDate.DateTime(nil); err == nil {
			evt.ExDates[i] = toTimeSpec(dt, &exDate)
		}
	}

	rDates := event.Props.Values(ical.PropRecurrenceDates)
	evt.RDates = make([]*calendarv1.TimeSpec, len(rDates))
	for i, rDate := range rDates {
		if dt, err := rDate.DateTime(nil); err == nil {
			evt.RDates[i] = toTimeSpec(dt, &rDate)
		}
	}

	if status, err := event.Props.Text(ical.PropStatus); err == nil {
		evt.Status = status
	}

	if sequenceProp := event.Props.Get(ical.PropSequence); sequenceProp != nil {
		if sequence, err := sequenceProp.Int(); err == nil {
			evt.Sequence = int32(sequence)
		}
	}

	if lastModified, err := event.Props.DateTime(ical.PropLastModified, nil); err == nil {
		evt.LastModified = timestamppb.New(lastModified)
	}

	if created, err := event.Props.DateTime(ical.PropCreated, nil); err == nil {
		evt.CreatedAt = timestamppb.New(created)
	}

	// UpdatedAt might be same as LastModified if not present
	if updated, err := event.Props.DateTime(ical.PropLastModified, nil); err == nil {
		evt.UpdatedAt = timestamppb.New(updated)
	}

	return evt, nil
}

func toTimeSpec(t time.Time, prop *ical.Prop) *calendarv1.TimeSpec {
	ts := &calendarv1.TimeSpec{
		Time: timestamppb.New(t),
	}
	if prop != nil {
		ts.Tzid = prop.Params.Get(ical.ParamTimezoneID)
	}
	return ts
}

func toCalAddress(prop *ical.Prop) *calendarv1.CalAddress {
	addr := &calendarv1.CalAddress{
		Email:    prop.Value,
		Cn:       prop.Params.Get(ical.ParamCommonName),
		Role:     prop.Params.Get(ical.ParamRole),
		PartStat: prop.Params.Get(ical.ParamParticipationStatus),
	}
	if rsvp := prop.Params.Get(ical.ParamRSVP); rsvp == "TRUE" {
		addr.Rsvp = true
	}
	return addr
}

func toRRule(roption *rrule.ROption) *calendarv1.RRule {
	rr := &calendarv1.RRule{
		Interval: int32(roption.Interval),
		Count:    int32(roption.Count),
		// ByDay:    roption.Byday, // TODO: find correct field name
	}
	// Map Freq int to string
	switch roption.Freq {
	case rrule.DAILY:
		rr.Freq = "DAILY"
	case rrule.WEEKLY:
		rr.Freq = "WEEKLY"
	case rrule.MONTHLY:
		rr.Freq = "MONTHLY"
	case rrule.YEARLY:
		rr.Freq = "YEARLY"
	}
	if !roption.Until.IsZero() {
		rr.Until = timestamppb.New(roption.Until)
	}
	rr.ByMonth = make([]int32, len(roption.Bymonth))
	for i, m := range roption.Bymonth {
		rr.ByMonth[i] = int32(m)
	}
	rr.BySetPos = make([]int32, len(roption.Bysetpos))
	for i, p := range roption.Bysetpos {
		rr.BySetPos[i] = int32(p)
	}
	return rr
}

func toVTimezone(comp *ical.Component) (*calendarv1.VTimezone, error) {
	tz := &calendarv1.VTimezone{}
	if tzid, err := comp.Props.Text(ical.PropTimezoneID); err == nil {
		tz.Tzid = tzid
	}
	// Offset from TZOFFSETTO, parsing +HHMM or -HHMM
	if offsetProp := comp.Props.Get(ical.PropTimezoneOffsetTo); offsetProp != nil {
		if offset, err := parseTimezoneOffset(offsetProp.Value); err == nil {
			tz.Offset = int32(offset)
		}
	}
	return tz, nil
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func parseTimezoneOffset(s string) (int, error) {
	if len(s) < 5 || len(s) > 6 {
		return 0, fmt.Errorf("invalid timezone offset format")
	}
	sign := 1
	if s[0] == '-' {
		sign = -1
	} else if s[0] != '+' {
		return 0, fmt.Errorf("invalid timezone offset format")
	}
	offsetStr := s[1:]
	var hours, minutes int
	if len(offsetStr) == 4 {
		// +HHMM
		if !isDigits(offsetStr) {
			return 0, fmt.Errorf("invalid timezone offset format")
		}
		hours = int(offsetStr[0]-'0')*10 + int(offsetStr[1]-'0')
		minutes = int(offsetStr[2]-'0')*10 + int(offsetStr[3]-'0')
	} else if len(offsetStr) == 5 && offsetStr[2] == ':' {
		// +HH:MM
		if !isDigits(offsetStr[:2]) || !isDigits(offsetStr[3:]) {
			return 0, fmt.Errorf("invalid timezone offset format")
		}
		hours = int(offsetStr[0]-'0')*10 + int(offsetStr[1]-'0')
		minutes = int(offsetStr[3]-'0')*10 + int(offsetStr[4]-'0')
	} else {
		return 0, fmt.Errorf("invalid timezone offset format")
	}
	if hours < 0 || hours > 23 || minutes < 0 || minutes > 59 {
		return 0, fmt.Errorf("invalid timezone offset format")
	}
	return sign * (hours*3600 + minutes*60), nil
}
