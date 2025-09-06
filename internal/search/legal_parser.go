package search

import (
	"regexp"
	"strings"
)

// Named struct for events
type LegalEvent struct {
	Date  string `json:"Date"`
	Event string `json:"Event"`
}

type LegalDocument struct {
	CaseNumber  string       `json:"CaseNumber"`
	Petitioners []string     `json:"Petitioners"`
	Respondents []string     `json:"Respondents"`
	Events      []LegalEvent `json:"Events"`
	Synopsis    string       `json:"Synopsis"`
}

func ParseLegalDocument(text string) *LegalDocument {
	doc := &LegalDocument{}

	// --- Case number ---
	reCase := regexp.MustCompile(`(?i)Complaint No[^\s]*[_/:]?\s*([\w\d/-]+)`)
	if m := reCase.FindStringSubmatch(text); len(m) > 1 {
		doc.CaseNumber = m[1]
	}

	// --- Petitioners ---
	rePet := regexp.MustCompile(`(?i)PETITIONERS\s*\n(.+?)\nAND`)
	if m := rePet.FindStringSubmatch(text); len(m) > 1 {
		doc.Petitioners = splitLines(m[1])
	}

	// --- Respondents ---
	reResp := regexp.MustCompile(`(?i)AND\s+(.+?)\nSYNOPSIS`)
	if m := reResp.FindStringSubmatch(text); len(m) > 1 {
		doc.Respondents = splitLines(m[1])
	}

	// --- Events: "dd/mm/yyyy | event description" ---
	reEvents := regexp.MustCompile(`(?m)(\d{2}/\d{2}/\d{4})\s*\|\s*(.+)`)
	for _, m := range reEvents.FindAllStringSubmatch(text, -1) {
		doc.Events = append(doc.Events, LegalEvent{
			Date:  m[1],
			Event: m[2],
		})
	}

	// --- Synopsis: everything after "SYNOPSIS" ---
	reSyn := regexp.MustCompile(`(?i)SYNOPSIS\s*\n([\s\S]+)`)
	if m := reSyn.FindStringSubmatch(text); len(m) > 1 {
		doc.Synopsis = strings.TrimSpace(m[1])
	}

	return doc
}

// Split lines, trim spaces, and remove empty lines
func splitLines(s string) []string {
	lines := strings.Split(s, "\n")
	var result []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			result = append(result, l)
		}
	}
	return result
}
