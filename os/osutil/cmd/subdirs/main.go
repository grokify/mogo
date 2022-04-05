package main

import (
	"fmt"

	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/os/osutil"
)

func main() {
	err := osutil.VisitPath(osutil.GoPathSrc(), nil, true, false, false, func(dir string) error {
		fmt.Println(dir)
		return nil
	})
	logutil.FatalErr(err)
	fmt.Println("DONE")
}
