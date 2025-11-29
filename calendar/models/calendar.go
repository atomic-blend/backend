package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Calendar represents an RFC 5545 VCALENDAR object.
type Calendar struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    *primitive.ObjectID `json:"userId" bson:"user_id,omitempty"`
	ProdID    string              `json:"prodId,omitempty" bson:"prod_id,omitempty"` // PRODID (required)
	Version   string              `json:"version" bson:"version"`                    // VERSION (must be "2.0")
	Method    string              `json:"method,omitempty" bson:"method,omitempty"`  // METHOD (optional, used for iTIP: REQUEST/PUBLISH/etc.)
	Events    []Event             `json:"events,omitempty" bson:"events,omitempty"`  // VEVENT components
	Timezones []VTimezone         `json:"timezones,omitempty" bson:"timezones,omitempty"`
	CreatedAt *primitive.DateTime `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *primitive.DateTime `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

// Event represents an RFC 5545 VEVENT object.
type Event struct {
	UID          string              // UID (required)
	DTStamp      time.Time           `json:"dtStamp" bson:"dt_stamp"`                               // DTSTAMP (required)
	Start        *TimeSpec           `json:"start,omitempty" bson:"start,omitempty"`                // DTSTART (recommended)
	End          *TimeSpec           `json:"end,omitempty" bson:"end,omitempty"`                    // DTEND (optional)
	Duration     *time.Duration      `json:"duration,omitempty" bson:"duration,omitempty"`          // DURATION (optional, alternative to DTEND)
	Summary      string              `json:"summary,omitempty" bson:"summary,omitempty"`            // SUMMARY
	Description  string              `json:"description,omitempty" bson:"description,omitempty"`    // DESCRIPTION
	Location     *string             `json:"location,omitempty" bson:"location,omitempty"`          // LOCATION
	Organizer    *CalAddress         `json:"organizer,omitempty" bson:"organizer,omitempty"`        // ORGANIZER
	Attendees    []CalAddress        `json:"attendees,omitempty" bson:"attendees,omitempty"`        // ATTENDEE
	Recurrence   *RRule              `json:"recurrence,omitempty" bson:"recurrence,omitempty"`      // RRULE
	ExDates      []TimeSpec          `json:"exDates,omitempty" bson:"ex_dates,omitempty"`           // EXDATE
	RDates       []TimeSpec          `json:"rDates,omitempty" bson:"r_dates,omitempty"`             // RDATE
	Status       string              `json:"status,omitempty" bson:"status,omitempty"`              // STATUS: CONFIRMED / TENTATIVE / CANCELLED
	Sequence     int                 `json:"sequence,omitempty" bson:"sequence,omitempty"`          // SEQUENCE (increments on update)
	LastModified *time.Time          `json:"lastModified,omitempty" bson:"last_modified,omitempty"` // LAST-MODIFIED
	CreatedAt    *primitive.DateTime `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    *primitive.DateTime `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}

// TimeSpec represents a date or date-time with an optional timezone.
// RFC 5545 allows floating times, UTC times, or TZID-bound times.
type TimeSpec struct {
	Time time.Time `json:"time" bson:"time"`
	TZID string    `json:"tzid,omitempty" bson:"tzid,omitempty"` // empty means floating or UTC depending on Time.Location
}

// CalAddress represents ORGANIZER and ATTENDEE fields.
type CalAddress struct {
	Email    string `json:"email" bson:"email"`                            // mailto: address without the "mailto:" prefix
	CN       string `json:"cn,omitempty" bson:"cn,omitempty"`              // Common Name
	Role     string `json:"role,omitempty" bson:"role,omitempty"`          // ROLE=REQ-PARTICIPANT / OPTIONAL / NON-PARTICIPANT
	PartStat string `json:"partStat,omitempty" bson:"part_stat,omitempty"` // PARTSTAT=ACCEPTED/TENTATIVE/DECLINED
	RSVP     bool   `json:"rsvp,omitempty" bson:"rsvp,omitempty"`          // RSVP=TRUE/FALSE
}

// RRule represents the RFC 5545 recurrence rule.
type RRule struct {
	Freq     string     `json:"freq" bson:"freq"`                             // FREQ=DAILY/WEEKLY/MONTHLY/YEARLY
	Interval int        `json:"interval,omitempty" bson:"interval,omitempty"` // INTERVAL=n
	Until    *time.Time `json:"until,omitempty" bson:"until,omitempty"`
	Count    int        `json:"count,omitempty" bson:"count,omitempty"`
	ByDay    []string   `json:"byDay,omitempty" bson:"by_day,omitempty"` // e.g. MO,WE,FR
	ByMonth  []int      `json:"byMonth,omitempty" bson:"by_month,omitempty"`
	BySetPos []int      `json:"bySetPos,omitempty" bson:"by_set_pos,omitempty"`
}

// VTimezone represents VTIMEZONE blocks.
// Minimal version; can be extended if needed.
type VTimezone struct {
	TZID   string `json:"tzid" bson:"tzid"`
	Offset int    `json:"offset,omitempty" bson:"offset,omitempty"` // offset from UTC in seconds (simplified model)
}
