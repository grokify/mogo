package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/grokify/mogo/data/currencyutil"
	"github.com/grokify/mogo/data/currencyutil/internal"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/os/osutil"
	"github.com/grokify/mogo/type/stringsutil"
)

func main() {
	//file := "../../currencies.tsv"
	//d, err := os.ReadFile(file)
	//logutil.FatalErr(err)
	lines := strings.Split(internal.CurrenciesDataRaw, "\n")
	fmtutil.PrintJSON(lines)

	currs := []currencyutil.Currency{}
	for _, line := range lines {
		c, err := procLine(line)
		if err != nil {
			log.Fatal(err)
		}
		if c != nil {
			currs = append(currs, *c)
		}
	}

	outputConstants("currency_codes.go", currs)

	fmt.Println("DONE")
}

func outputConstants(filename string, currs currencyutil.Currencies) error {
	set, err := currencyutil.NewCurrencySet(currs...)
	if err != nil {
		return err
	}
	fw, err := osutil.NewFileWriter(filename)
	if err != nil {
		return err
	}
	_, err = fw.WriteString(true, "package currencyutil", "const (")
	if err != nil {
		return err
	}
	codes := set.Codes()
	for _, code := range codes {
		if !stringsutil.IsUpper(code) || len(code) != 3 {
			continue
		}
		curr, ok := set.Map[code]
		if !ok {
			panic("code not found")
		}
		fw.WriteStringf(true, "Currency%s = \"%s\" // %s %s", code, code, curr.Country, curr.Name)
	}
	_, err = fw.WriteString(true, ")")
	if err != nil {
		return err
	}

	_, err = fw.WriteString(true,
		"func CurrencyCodesAll() []string {",
		"return []string{")
	if err != nil {
		return err
	}
	for i, code := range codes {
		if !stringsutil.IsUpper(code) || len(code) != 3 {
			continue
		}
		curr, ok := set.Map[code]
		if !ok {
			panic("code not found")
		}
		_, err = fw.WriteString(true, "Currency"+code+", // "+curr.Country+" "+curr.Name)
		if err != nil {
			return err
		}
		if i == len(codes)-1 {
			_, err = fw.WriteString(true, "}}")
			if err != nil {
				return err
			}
		}
	}
	return fw.Close()
}

func procLine(line string) (*currencyutil.Currency, error) {
	parts := stringsutil.SliceCondenseSpace(strings.Split(line, "\t"), false, false)
	fmtutil.PrintJSON(parts)
	if len(parts) == 0 {
		return nil, nil
	}
	if len(parts) < 4 {
		return nil, fmt.Errorf("wrong number of parts [%d]\n", len(parts))
		// panic(fmt.Sprintf("wrong number of parts [%d]\n", len(parts)))
	}
	c := &currencyutil.Currency{
		Code:    parts[0],
		Symbol:  parts[1],
		Country: parts[2],
		Name:    parts[3],
	}
	c.TrimSpace()
	fmtutil.PrintJSON(c)
	return c, nil
}
