package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/roeyazroel/linear-tui/internal/linearapi"
)

// GroupByField represents how issues are grouped in the issues column.
type GroupByField string

const (
	GroupByNone     GroupByField = "none"
	GroupByAssignee GroupByField = "assignee"
	GroupByStatus   GroupByField = "status"
	GroupByPriority GroupByField = "priority"
)

// IssueGroup is a named collection of issues for display in a single table.
type IssueGroup struct {
	Key    string             // unique key for the group (e.g. user ID, state name)
	Title  string             // display title for the table border
	Issues []linearapi.Issue  // issues belonging to this group
}

// groupIssues splits issues into groups based on the selected GroupByField.
// When groupBy is GroupByNone, issues are split into "My Issues" / "Other Issues"
// using the legacy assignee-based split.
func groupIssues(issues []linearapi.Issue, groupBy GroupByField, currentUserID string) []IssueGroup {
	switch groupBy {
	case GroupByAssignee:
		return groupByAssignee(issues)
	case GroupByStatus:
		return groupByStatus(issues)
	case GroupByPriority:
		return groupByPriority(issues)
	default:
		return groupByMyOther(issues, currentUserID)
	}
}

// groupByMyOther is the legacy grouping: "My Issues" and "Other Issues".
func groupByMyOther(issues []linearapi.Issue, currentUserID string) []IssueGroup {
	my, other := splitIssuesByAssignee(issues, currentUserID)
	var groups []IssueGroup
	if len(my) > 0 {
		groups = append(groups, IssueGroup{
			Key:    "_my",
			Title:  "My Issues",
			Issues: my,
		})
	}
	groups = append(groups, IssueGroup{
		Key:    "_other",
		Title:  "Other Issues",
		Issues: other,
	})
	return groups
}

// groupByAssignee groups issues by their assignee name.
func groupByAssignee(issues []linearapi.Issue) []IssueGroup {
	type groupEntry struct {
		name   string
		issues []linearapi.Issue
	}

	byAssignee := make(map[string]*groupEntry)
	var order []string

	for _, issue := range issues {
		key := issue.AssigneeID
		name := issue.Assignee
		if key == "" {
			key = "_unassigned"
			name = "Unassigned"
		}
		entry, exists := byAssignee[key]
		if !exists {
			entry = &groupEntry{name: name}
			byAssignee[key] = entry
			order = append(order, key)
		}
		entry.issues = append(entry.issues, issue)
	}

	// Sort groups alphabetically by name, but put "Unassigned" last
	sort.Slice(order, func(i, j int) bool {
		if order[i] == "_unassigned" {
			return false
		}
		if order[j] == "_unassigned" {
			return true
		}
		return byAssignee[order[i]].name < byAssignee[order[j]].name
	})

	groups := make([]IssueGroup, 0, len(order))
	for _, key := range order {
		entry := byAssignee[key]
		groups = append(groups, IssueGroup{
			Key:    key,
			Title:  fmt.Sprintf("%s (%d)", entry.name, len(entry.issues)),
			Issues: entry.issues,
		})
	}
	return groups
}

// groupByStatus groups issues by their workflow state.
func groupByStatus(issues []linearapi.Issue) []IssueGroup {
	type groupEntry struct {
		name   string
		issues []linearapi.Issue
	}

	byStatus := make(map[string]*groupEntry)
	var order []string

	for _, issue := range issues {
		key := issue.StateID
		name := issue.State
		if key == "" {
			key = "_unknown"
			name = "Unknown"
		}
		entry, exists := byStatus[key]
		if !exists {
			entry = &groupEntry{name: name}
			byStatus[key] = entry
			order = append(order, key)
		}
		entry.issues = append(entry.issues, issue)
	}

	groups := make([]IssueGroup, 0, len(order))
	for _, key := range order {
		entry := byStatus[key]
		groups = append(groups, IssueGroup{
			Key:    key,
			Title:  fmt.Sprintf("%s (%d)", entry.name, len(entry.issues)),
			Issues: entry.issues,
		})
	}
	return groups
}

// isCompletedState returns true if the state name indicates a completed or canceled state.
// Uses the same heuristics as state coloring in issues_table.go.
func isCompletedState(stateName string) bool {
	lower := strings.ToLower(stateName)
	return strings.Contains(lower, "done") ||
		strings.Contains(lower, "complete") ||
		strings.Contains(lower, "cancel")
}

// filterActiveIssues returns only issues whose state is not completed or canceled.
func filterActiveIssues(issues []linearapi.Issue) []linearapi.Issue {
	active := make([]linearapi.Issue, 0, len(issues))
	for _, issue := range issues {
		if !isCompletedState(issue.State) {
			active = append(active, issue)
		}
	}
	return active
}

// groupByPriority groups issues by priority level.
func groupByPriority(issues []linearapi.Issue) []IssueGroup {
	priorityNames := map[int]string{
		1: "Urgent",
		2: "High",
		3: "Normal",
		4: "Low",
		0: "No Priority",
	}
	// Display order: Urgent, High, Normal, Low, No Priority
	priorityOrder := []int{1, 2, 3, 4, 0}

	buckets := make(map[int][]linearapi.Issue)
	for _, issue := range issues {
		buckets[issue.Priority] = append(buckets[issue.Priority], issue)
	}

	var groups []IssueGroup
	for _, p := range priorityOrder {
		if items, ok := buckets[p]; ok && len(items) > 0 {
			groups = append(groups, IssueGroup{
				Key:    fmt.Sprintf("p%d", p),
				Title:  fmt.Sprintf("%s (%d)", priorityNames[p], len(items)),
				Issues: items,
			})
		}
	}
	return groups
}
