package tokenizer

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/net/html/atom"
)

type AtomSet struct {
	Atoms map[string]atom.Atom
}

func NewAtomSet(htmlAtoms ...atom.Atom) AtomSet {
	atomSet := AtomSet{
		Atoms: map[string]atom.Atom{}}
	if len(htmlAtoms) > 0 {
		atomSet.Add(htmlAtoms...)
	}
	return atomSet
}

func NewAtomSetStringMust(tagNames ...string) AtomSet {
	atoms, err := NewAtomSetString(tagNames...)
	if err != nil {
		panic(err)
	}
	return atoms
}

func NewAtomSetString(tagNames ...string) (AtomSet, error) {
	if len(tagNames) == 0 {
		return NewAtomSet(), nil
	}
	atoms := []atom.Atom{}
	unmatchedNames := []string{}
	for _, tagName := range tagNames {
		tagAtom := AtomLookupString(tagName)
		if strings.TrimSpace(tagAtom.String()) == "" {
			unmatchedNames = append(unmatchedNames, tagName)
		} else {
			atoms = append(atoms, tagAtom)
		}
	}
	atomsSet := NewAtomSet(atoms...)
	if len(unmatchedNames) > 0 {
		return atomsSet, fmt.Errorf("unmatchedTagNames [%s]", strings.Join(unmatchedNames, ","))
	}
	return atomsSet, nil
}

func AtomLookupString(tagName string) atom.Atom {
	return atom.Lookup([]byte(strings.ToLower(strings.TrimSpace(tagName))))
}

func (set AtomSet) Len() int {
	return len(set.Atoms)
}

func (set AtomSet) Add(htmlAtoms ...atom.Atom) {
	for _, htmlAtom := range htmlAtoms {
		set.Atoms[htmlAtom.String()] = htmlAtom
	}
}

func (set AtomSet) Exists(htmlAtom atom.Atom) bool {
	if _, ok := set.Atoms[htmlAtom.String()]; ok {
		return true
	}
	return false
}

func (set AtomSet) Names() []string {
	names := []string{}
	for _, htmlAtom := range set.Atoms {
		names = append(names, htmlAtom.String())
	}
	sort.Strings(names)
	return names
}
