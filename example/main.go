package main

import (
	"fmt"
	"github.com/lincaiyong/ql"
)

func main() {
	// https://codeql.github.com/docs/codeql-language-guides/analyzing-data-flow-in-go/
	q := `from Function osOpen, CallExpr call
where
  osOpen.hasQualifiedName("os", "Open") and
  call.getTarget() = osOpen
select call.getArgument(0)
`
	ret := ql.Query(q)
	fmt.Println(ret)
}
