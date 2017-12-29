package humannameparser

import (
	"strings"
)

type HumanName struct {
	Name        string
	LeadingInit string
	FirstName   string
	MiddleName  string
	LastName    string
	Suffix      string

	Suffixes string
	Prefixes string
}

// ParseHumanName is a initial stub function to parse human names.
// The goal is to have a full implementation like a port of:
// https://github.com/jasonpriem/HumanNameParser.php
func ParseHumanName(n string) (*HumanName, error) {
	h := &HumanName{}
	parts := strings.Split(n, " ")
	if len(parts) == 2 {
		h.FirstName = parts[0]
		h.LastName = parts[1]
	} else if len(parts) == 3 {
		h.FirstName = parts[0]
		h.MiddleName = parts[1]
		h.LastName = parts[2]
	}
	return h, nil
}
