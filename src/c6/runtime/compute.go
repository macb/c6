package runtime

import "c6/ast"
import "fmt"

/*
Used for Incompatible unit, data type or unsupported operations

TODO: This is not used yet. our compute functions should return error if possible
*/
type ComputeError struct {
	Message string
	Left    ast.Value
	Right   ast.Value
}

func (self ComputeError) Error() string {
	return self.Message
}

/*
Value
*/
type ComputeFunction func(a ast.Value, b ast.Value) ast.Value

const ValueTypeNum = 7

var computableMatrix [ValueTypeNum][ValueTypeNum]bool = [ValueTypeNum][ValueTypeNum]bool{
	/* NumberValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* HexColorValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* RGBAColorValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},

	/* RGBColorValue */
	[ValueTypeNum]bool{false, false, false, false, false, false, false},
}

/**
Each row: [5]ComputeFunction{ NumberValue, HexColorValue, RGBAColorValue, RGBColorValue }
*/
var computeFunctionMatrix [5][5]ComputeFunction = [5][5]ComputeFunction{

	/* NumberValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* HexColorValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* RGBAColorValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},

	/* RGBColorValue */
	[5]ComputeFunction{nil, nil, nil, nil, nil},
}

func ComputeBoolean(op *ast.Op, a ast.Value, b ast.Value) ast.Value {
	if op == nil {
		panic("op can't be nil")
	}
	switch op.Type {
	case ast.T_LOGICAL_AND:
		switch ta := a.(type) {
		case *ast.Boolean:
			switch tb := b.(type) {
			case *ast.Boolean:
				return ast.NewBoolean(ta.Value && tb.Value)
			}
		}

	case ast.T_LOGICAL_OR:
		switch ta := a.(type) {
		case *ast.Boolean:
			switch tb := b.(type) {
			case *ast.Boolean:
				return ast.NewBoolean(ta.Value || tb.Value)
			}
		}
	}
	return nil
}

func Compute(op *ast.Op, a ast.Value, b ast.Value) ast.Value {
	if op == nil {
		panic("op can't be nil")
	}
	switch op.Type {
	case ast.T_PLUS:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberAddNumber(ta, tb)
			case *ast.HexColor:
				return HexColorAddNumber(tb, ta)
			}
		case *ast.HexColor:
			switch tb := b.(type) {
			case *ast.Number:
				return HexColorAddNumber(ta, tb)
			}
		case *ast.RGBColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBColorAddNumber(ta, tb)
			}
		case *ast.RGBAColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBAColorAddNumber(ta, tb)
			}
		}
	case ast.T_MINUS:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberSubNumber(ta, tb)
			}
		case *ast.HexColor:
			switch tb := b.(type) {
			case *ast.Number:
				return HexColorSubNumber(ta, tb)
			}

		case *ast.RGBColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBColorSubNumber(ta, tb)
			}

		case *ast.RGBAColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBAColorSubNumber(ta, tb)
			}
		}

	case ast.T_DIV:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberDivNumber(ta, tb)
			}
		case *ast.HexColor:
			switch tb := b.(type) {
			case *ast.Number:
				return HexColorDivNumber(ta, tb)
			}
		case *ast.RGBColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBColorDivNumber(ta, tb)
			}
		case *ast.RGBAColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBAColorDivNumber(ta, tb)
			}
		}

	case ast.T_MUL:
		switch ta := a.(type) {
		case *ast.Number:
			switch tb := b.(type) {
			case *ast.Number:
				return NumberMulNumber(ta, tb)
			}

		case *ast.HexColor:
			switch tb := b.(type) {
			case *ast.Number:
				return HexColorMulNumber(ta, tb)
			}

		case *ast.RGBColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBColorMulNumber(ta, tb)
			}

		case *ast.RGBAColor:
			switch tb := b.(type) {
			case *ast.Number:
				return RGBAColorMulNumber(ta, tb)
			}
		}
	}
	return nil
}

func IsConstantValue(val ast.Value) bool {
	switch val.(type) {
	case *ast.Number, *ast.HexColor, *ast.RGBColor, *ast.RGBAColor:
		return true
	}
	return false
}

func EvaluateBinaryExpressionInBooleanContext(expr *ast.BinaryExpression, symTable *SymTable) ast.Value {

	var lval ast.Value = nil
	var rval ast.Value = nil

	switch expr := expr.Left.(type) {
	case *ast.UnaryExpression:
		lval = EvaluateUnaryExpressionInBooleanContext(expr, symTable)

	case *ast.BinaryExpression:
		lval = EvaluateBinaryExpressionInBooleanContext(expr, symTable)

	// operands
	case *ast.Number, *ast.HexColor, *ast.RGBColor, *ast.RGBAColor:
		if bval, ok := lval.(ast.BooleanValue); ok {
			lval = ast.NewBoolean(bval.Boolean())
		} else {
			panic(fmt.Errorf("BooleanValue interface is not support for %+v", lval))
		}
	}

	switch expr := expr.Right.(type) {
	case *ast.UnaryExpression:
		rval = EvaluateUnaryExpressionInBooleanContext(expr, symTable)

	case *ast.BinaryExpression:
		rval = EvaluateBinaryExpressionInBooleanContext(expr, symTable)

	// operands
	case *ast.Number, *ast.HexColor, *ast.RGBColor, *ast.RGBAColor:
		if bval, ok := rval.(ast.BooleanValue); ok {
			rval = ast.NewBoolean(bval.Boolean())
		} else {
			panic(fmt.Errorf("BooleanValue interface is not support for %+v", rval))
		}
	}
	if lval != nil && rval != nil {
		return ComputeBoolean(expr.Op, lval, rval)
	}
	return nil
}

func EvaluateUnaryExpressionInBooleanContext(expr *ast.UnaryExpression, symTable *SymTable) ast.Value {
	var val ast.Value = nil

	switch t := expr.Expr.(type) {
	case *ast.BinaryExpression:
		val = EvaluateBinaryExpression(t, symTable)
	case *ast.UnaryExpression:
		val = EvaluateUnaryExpression(t, symTable)
	default:
		val = ast.Value(t)
	}

	switch expr.Op.Type {
	case ast.T_LOGICAL_NOT:
		if bval, ok := val.(ast.BooleanValue); ok {
			return ast.NewBoolean(bval.Boolean())
		} else {
			panic(fmt.Errorf("BooleanValue interface is not support for %+v", val))
		}
	}
	return val
}

/*
EvaluateExpression calls EvaluateBinaryExpression. except EvaluateExpression
prevents calculate css slash as division.  otherwise it's the same as
EvaluateBinaryExpression.
*/
func EvaluateExpression(expr ast.Expression, symTable *SymTable) ast.Value {

	switch t := expr.(type) {

	case *ast.BinaryExpression:
		if t.IsCssSlash() {
			// return string object without quote
			return ast.NewString(0, t.Left.(*ast.Number).String()+"/"+t.Right.(*ast.Number).String(), nil)
		}
		return EvaluateBinaryExpression(t, symTable)

	case *ast.UnaryExpression:
		return EvaluateUnaryExpression(t, symTable)

	case *ast.Number, *ast.HexColor, *ast.RGBColor, *ast.RGBAColor:
		return ast.Value(expr)
	}

	panic("Unsupported expression type")
	return nil
}

/**
EvaluateBinaryExpression recursively.
*/
func EvaluateBinaryExpression(expr *ast.BinaryExpression, symTable *SymTable) ast.Value {
	var lval ast.Value = nil
	var rval ast.Value = nil

	switch expr := expr.Left.(type) {

	case *ast.BinaryExpression:
		lval = EvaluateBinaryExpression(expr, symTable)

	case *ast.UnaryExpression:
		lval = EvaluateUnaryExpression(expr, symTable)

	case *ast.Number, *ast.HexColor, *ast.RGBColor, *ast.RGBAColor:
		lval = ast.Value(expr)

	}

	switch expr := expr.Right.(type) {

	case *ast.UnaryExpression:
		rval = EvaluateUnaryExpression(expr, symTable)

	case *ast.BinaryExpression:
		rval = EvaluateBinaryExpression(expr, symTable)

	case *ast.Number, *ast.HexColor, *ast.RGBColor, *ast.RGBAColor:
		rval = ast.Value(expr)

	}

	if lval != nil && rval != nil {
		return Compute(expr.Op, lval, rval)
	}
	return nil
}

func EvaluateUnaryExpression(expr *ast.UnaryExpression, symTable *SymTable) ast.Value {
	var val ast.Value = nil

	switch t := expr.Expr.(type) {
	case *ast.BinaryExpression:
		val = EvaluateBinaryExpression(t, symTable)
	case *ast.UnaryExpression:
		val = EvaluateUnaryExpression(t, symTable)
	default:
		val = ast.Value(t)
	}

	switch expr.Op.Type {
	case ast.T_MINUS:
		switch n := val.(type) {
		case *ast.Number:
			n.Value = -n.Value
		}
	}
	return val
}
