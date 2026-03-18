package tui

import (
	"testing"

	"github.com/roeyazroel/linear-tui/internal/linearapi"
)

func TestGroupByAssignee(t *testing.T) {
	issues := []linearapi.Issue{
		{ID: "1", Assignee: "Alice", AssigneeID: "a1", Title: "Task 1"},
		{ID: "2", Assignee: "Bob", AssigneeID: "b1", Title: "Task 2"},
		{ID: "3", Assignee: "Alice", AssigneeID: "a1", Title: "Task 3"},
		{ID: "4", Assignee: "", AssigneeID: "", Title: "Task 4"},
		{ID: "5", Assignee: "Bob", AssigneeID: "b1", Title: "Task 5"},
	}

	groups := groupByAssignee(issues)

	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}

	// Alphabetical: Alice, Bob, then Unassigned
	if groups[0].Key != "a1" {
		t.Errorf("expected first group key=a1, got %s", groups[0].Key)
	}
	if len(groups[0].Issues) != 2 {
		t.Errorf("expected Alice to have 2 issues, got %d", len(groups[0].Issues))
	}

	if groups[1].Key != "b1" {
		t.Errorf("expected second group key=b1, got %s", groups[1].Key)
	}
	if len(groups[1].Issues) != 2 {
		t.Errorf("expected Bob to have 2 issues, got %d", len(groups[1].Issues))
	}

	if groups[2].Key != "_unassigned" {
		t.Errorf("expected third group key=_unassigned, got %s", groups[2].Key)
	}
	if len(groups[2].Issues) != 1 {
		t.Errorf("expected Unassigned to have 1 issue, got %d", len(groups[2].Issues))
	}
}

func TestGroupByStatus(t *testing.T) {
	issues := []linearapi.Issue{
		{ID: "1", State: "In Progress", StateID: "s1"},
		{ID: "2", State: "Done", StateID: "s2"},
		{ID: "3", State: "In Progress", StateID: "s1"},
	}

	groups := groupByStatus(issues)

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if len(groups[0].Issues) != 2 {
		t.Errorf("expected first group to have 2 issues, got %d", len(groups[0].Issues))
	}
}

func TestGroupByPriority(t *testing.T) {
	issues := []linearapi.Issue{
		{ID: "1", Priority: 1}, // Urgent
		{ID: "2", Priority: 3}, // Normal
		{ID: "3", Priority: 1}, // Urgent
		{ID: "4", Priority: 0}, // No priority
	}

	groups := groupByPriority(issues)

	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	// Order: Urgent, Normal, No Priority
	if groups[0].Key != "p1" {
		t.Errorf("expected first group=p1 (Urgent), got %s", groups[0].Key)
	}
	if groups[1].Key != "p3" {
		t.Errorf("expected second group=p3 (Normal), got %s", groups[1].Key)
	}
	if groups[2].Key != "p0" {
		t.Errorf("expected third group=p0 (No Priority), got %s", groups[2].Key)
	}
}

func TestFilterActiveIssues(t *testing.T) {
	issues := []linearapi.Issue{
		{ID: "1", State: "In Progress"},
		{ID: "2", State: "Done"},
		{ID: "3", State: "Todo"},
		{ID: "4", State: "Canceled"},
		{ID: "5", State: "Completed"},
		{ID: "6", State: "Backlog"},
	}

	active := filterActiveIssues(issues)

	if len(active) != 3 {
		t.Fatalf("expected 3 active issues, got %d", len(active))
	}
	expectedIDs := map[string]bool{"1": true, "3": true, "6": true}
	for _, issue := range active {
		if !expectedIDs[issue.ID] {
			t.Errorf("unexpected active issue ID=%s State=%s", issue.ID, issue.State)
		}
	}
}

func TestIsCompletedState(t *testing.T) {
	tests := []struct {
		state    string
		expected bool
	}{
		{"Done", true},
		{"Completed", true},
		{"Canceled", true},
		{"In Progress", false},
		{"Todo", false},
		{"Backlog", false},
		{"done", true},
		{"CANCELED", true},
	}
	for _, tt := range tests {
		if got := isCompletedState(tt.state); got != tt.expected {
			t.Errorf("isCompletedState(%q) = %v, want %v", tt.state, got, tt.expected)
		}
	}
}

func TestGroupByNoneFallsBackToMyOther(t *testing.T) {
	issues := []linearapi.Issue{
		{ID: "1", AssigneeID: "me", Title: "My task"},
		{ID: "2", AssigneeID: "other", Title: "Other task"},
	}

	groups := groupIssues(issues, GroupByNone, "me")

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups (My/Other), got %d", len(groups))
	}
	if groups[0].Key != "_my" {
		t.Errorf("expected first group=_my, got %s", groups[0].Key)
	}
	if groups[1].Key != "_other" {
		t.Errorf("expected second group=_other, got %s", groups[1].Key)
	}
}
