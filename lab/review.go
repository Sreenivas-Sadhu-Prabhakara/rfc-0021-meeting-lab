package lab

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Meeting struct {
	Title string  `json:"title"`
	AsOf  string  `json:"as_of"`
	Items []Entry `json:"items"`
}

type Entry struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	Owner    string `json:"owner,omitempty"`
	Due      string `json:"due,omitempty"`
	Done     bool   `json:"done,omitempty"`
	Evidence string `json:"evidence,omitempty"`
}

type Action struct {
	Entry
	Status   string   `json:"status"`
	Findings []string `json:"findings"`
}

type Report struct {
	Title     string   `json:"title"`
	AsOf      string   `json:"as_of"`
	Decisions []Entry  `json:"decisions"`
	Questions []Entry  `json:"questions"`
	Notes     []Entry  `json:"notes"`
	Actions   []Action `json:"actions"`
}

var entryID = regexp.MustCompile(`^[a-z][a-z0-9-]{0,31}$`)

// Review classifies entries as provided and checks whether actions can be followed up.
func Review(meeting Meeting) (Report, error) {
	if strings.TrimSpace(meeting.Title) == "" {
		return Report{}, fmt.Errorf("meeting title is required")
	}
	asOf, err := time.Parse(time.DateOnly, meeting.AsOf)
	if err != nil {
		return Report{}, fmt.Errorf("as_of must be a valid YYYY-MM-DD date")
	}
	if len(meeting.Items) == 0 || len(meeting.Items) > 200 {
		return Report{}, fmt.Errorf("expected 1 to 200 entries")
	}

	report := Report{
		Title: meeting.Title, AsOf: meeting.AsOf,
		Decisions: []Entry{}, Questions: []Entry{}, Notes: []Entry{}, Actions: []Action{},
	}
	seen := make(map[string]bool, len(meeting.Items))
	for _, entry := range meeting.Items {
		if !entryID.MatchString(entry.ID) || seen[entry.ID] {
			return Report{}, fmt.Errorf("invalid or duplicate entry ID %q", entry.ID)
		}
		seen[entry.ID] = true
		if strings.TrimSpace(entry.Summary) == "" || strings.TrimSpace(entry.Source) == "" {
			return Report{}, fmt.Errorf("entry %s needs a summary and source", entry.ID)
		}
		if entry.Kind != "action" && (entry.Owner != "" || entry.Due != "" || entry.Done || entry.Evidence != "") {
			return Report{}, fmt.Errorf("entry %s has action fields but is not an action", entry.ID)
		}
		switch entry.Kind {
		case "decision":
			report.Decisions = append(report.Decisions, entry)
		case "question":
			report.Questions = append(report.Questions, entry)
		case "note":
			report.Notes = append(report.Notes, entry)
		case "action":
			action, err := reviewAction(entry, asOf)
			if err != nil {
				return Report{}, err
			}
			report.Actions = append(report.Actions, action)
		default:
			return Report{}, fmt.Errorf("entry %s has unknown kind %q", entry.ID, entry.Kind)
		}
	}
	return report, nil
}

func reviewAction(entry Entry, asOf time.Time) (Action, error) {
	action := Action{Entry: entry, Status: "open", Findings: []string{}}
	if strings.TrimSpace(entry.Owner) == "" {
		action.Findings = append(action.Findings, "Missing owner")
	}

	if entry.Due == "" {
		action.Findings = append(action.Findings, "Missing deadline")
	} else {
		due, err := time.Parse(time.DateOnly, entry.Due)
		if err != nil {
			return Action{}, fmt.Errorf("entry %s has an invalid deadline", entry.ID)
		}
		if due.Before(asOf) {
			action.Status = "overdue"
		}
	}

	if entry.Done {
		action.Status = "closed"
		if strings.TrimSpace(entry.Evidence) == "" {
			action.Status = "unverified"
			action.Findings = append(action.Findings, "Missing completion evidence")
		}
	}
	return action, nil
}
