package classroom

import "github.com/expr-lang/expr"

const DEFAULT_EXPR = `passed / (passed+failed)`

type Assignments []AssignmentSpec

func (a AssignmentSpec) Score(passed, failed, cover float64) (float64, error) {
	env := map[string]any{
		"passed": passed,
		"failed": failed,
		"cover":  cover,
	}
	if a.Expr == "" {
		a.Expr = DEFAULT_EXPR
	}
	pgm, err := expr.Compile(a.Expr, expr.Env(env))
	if err != nil {
		return 0.0, err
	}
	result, err := expr.Run(pgm, env)
	if err != nil {
		return 0.0, err
	}
	return result.(float64), nil
}

type AssignmentSpec struct {
	Name      string
	Path      string
	Course    string
	ExtrasSrc string
	ExtrasDst string
	Expr      string
}
