package tokenizer

import (
	"sort"

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
