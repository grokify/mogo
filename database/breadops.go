package database

import (
	"strings"

	"github.com/grokify/gotilla/type/stringsutil"
)

// BreadOps uses the BREAD acronym to store identifiers for each
// operation.
type BreadOps struct {
	Name   string
	Browse []string
	Read   []string
	Edit   []string
	Add    []string
	Delete []string
}

func NewBreadOps(name string) BreadOps {
	return BreadOps{
		Name:   name,
		Browse: []string{},
		Read:   []string{},
		Edit:   []string{},
		Add:    []string{},
		Delete: []string{}}
}

func (bo *BreadOps) TrimSpace(dedupe, sort bool) {
	bo.Name = strings.TrimSpace(bo.Name)
	if bo.Browse != nil {
		bo.Browse = stringsutil.SliceCondenseSpace(bo.Browse, dedupe, sort)
	}
	if bo.Read != nil {
		bo.Read = stringsutil.SliceCondenseSpace(bo.Read, dedupe, sort)
	}
	if bo.Edit != nil {
		bo.Edit = stringsutil.SliceCondenseSpace(bo.Edit, dedupe, sort)
	}
	if bo.Add != nil {
		bo.Add = stringsutil.SliceCondenseSpace(bo.Add, dedupe, sort)
	}
	if bo.Delete != nil {
		bo.Delete = stringsutil.SliceCondenseSpace(bo.Delete, dedupe, sort)
	}
}
