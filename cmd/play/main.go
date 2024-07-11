package main

import (
	"fmt"
	"time"

	"github.com/expr-lang/expr"
)

func main() {
	env := map[string]interface{}{
		"lateness": time.Duration(0),
		"score":    80,
	}

	code := `score * (int(lateness) == 0 ? 1 : 0)`

	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}
