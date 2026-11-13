package lab

import (
	"reflect"
	"strings"
	"testing"
)

func TestReviewKeepsCategoriesAndSource(t *testing.T) {
	input := Meeting{Title: "Review", AsOf: "2026-11-13", Items: []Entry{
		{ID: "q", Kind: "question", Summary: "Expiry at capacity?", Source: "Minutes §1"},
		{ID: "d", Kind: "decision", Summary: "Adopt the retention rule", Source: "Minutes §2"},
		{ID: "n", Kind: "note", Summary: "Limit reported", Source: "Minutes §3"},
		{ID: "a", Kind: "action", Summary: "Publish notes", Source: "Minutes §4", Owner: "Gregg"},
	}}
	got, err := Review(input)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Questions, input.Items[0:1]) ||
		!reflect.DeepEqual(got.Decisions, input.Items[1:2]) ||
		!reflect.DeepEqual(got.Notes, input.Items[2:3]) {
		t.Fatalf("unexpected register: %+v", got)
	}
	want := Action{Entry: input.Items[3], Status: "open", Findings: []string{"Missing deadline"}}
	if len(got.Actions) != 1 || !reflect.DeepEqual(got.Actions[0], want) {
		t.Fatalf("actions = %+v, want %+v", got.Actions, want)
	}
}

func TestActionStatusAndFindings(t *testing.T) {
	for _, tc := range []struct {
		name, due, evidence, status string
		done                        bool
		findings                    []string
	}{
		{name: "due today", due: "2026-11-13", status: "open", findings: []string{}},
		{name: "past due", due: "2026-11-12", status: "overdue", findings: []string{}},
		{name: "claimed complete", due: "2026-11-12", done: true, evidence: " ", status: "unverified", findings: []string{"Missing completion evidence"}},
		{name: "evidence recorded", due: "2026-11-12", done: true, evidence: "Release 42", status: "closed", findings: []string{}},
		{name: "no deadline", status: "open", findings: []string{"Missing deadline"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			entry := Entry{ID: "a", Kind: "action", Summary: "Publish", Source: "Minutes", Owner: "Team", Due: tc.due, Done: tc.done, Evidence: tc.evidence}
			report, err := Review(Meeting{Title: "Review", AsOf: "2026-11-13", Items: []Entry{entry}})
			if err != nil {
				t.Fatal(err)
			}
			action := report.Actions[0]
			if action.Status != tc.status || !reflect.DeepEqual(action.Findings, tc.findings) {
				t.Fatalf("got %+v, want status %q and findings %v", action, tc.status, tc.findings)
			}
		})
	}
}

func TestReviewRejectsMalformedEntries(t *testing.T) {
	valid := Entry{ID: "a", Kind: "action", Summary: "Publish", Source: "Minutes"}
	for _, tc := range []struct {
		name, date, want string
		items            []Entry
	}{
		{"invalid date", "2026-02-30", "as_of", []Entry{valid}},
		{"duplicate ID", "2026-11-13", "duplicate", []Entry{valid, valid}},
		{"invalid deadline", "2026-11-13", "invalid deadline", []Entry{{ID: "a", Kind: "action", Summary: "Publish", Source: "Minutes", Due: "tomorrow"}}},
		{"ambiguous kind", "2026-11-13", "unknown kind", []Entry{{ID: "a", Kind: "proposal", Summary: "Publish", Source: "Minutes"}}},
		{"missing source", "2026-11-13", "summary and source", []Entry{{ID: "a", Kind: "note", Summary: "Publish"}}},
		{"mixed fields", "2026-11-13", "action fields", []Entry{{ID: "a", Kind: "question", Summary: "Publish?", Source: "Minutes", Owner: "Team"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Review(Meeting{Title: "Review", AsOf: tc.date, Items: tc.items})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("got %v, want error containing %q", err, tc.want)
			}
		})
	}
}

func TestMissingOwnerAndDeadline(t *testing.T) {
	report, err := Review(Meeting{Title: "Review", AsOf: "2026-11-13", Items: []Entry{
		{ID: "a", Kind: "action", Summary: "Publish", Source: "Minutes", Owner: " "},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if got := report.Actions[0].Findings; !reflect.DeepEqual(got, []string{"Missing owner", "Missing deadline"}) {
		t.Fatalf("unexpected findings: %v", got)
	}
}
