package currencyutil

import (
	"errors"
	"sort"
	"strings"

	"github.com/grokify/mogo/type/maputil"
)

type Currencies []Currency

func (c Currencies) Codes() []string {
	codes := []string{}
	for _, ci := range c {
		codes = append(codes, ci.Code)
	}
	sort.Strings(codes)
	return codes
}

type Currency struct {
	Code    string
	Symbol  string
	Country string
	Name    string
}

func (c *Currency) TrimSpace() {
	c.Code = strings.TrimSpace(c.Code)
	c.Symbol = strings.TrimSpace(c.Symbol)
	c.Country = strings.TrimSpace(c.Country)
	c.Name = strings.TrimSpace(c.Name)
}

type CurrencySet struct {
	Map map[string]Currency
}

var ErrNoCurrencyCode = errors.New("currency code not set")

func NewCurrencySet(c ...Currency) (CurrencySet, error) {
	set := CurrencySet{Map: map[string]Currency{}}
	err := set.Add(c...)
	return set, err
}

func (set *CurrencySet) Add(c ...Currency) error {
	for _, ci := range c {
		if strings.TrimSpace(ci.Code) == "" {
			return ErrNoCurrencyCode
		}
		set.Map[ci.Code] = ci
	}
	return nil
}

func (set CurrencySet) Codes() []string {
	return maputil.Keys(set.Map)
}
