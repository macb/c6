package ast

type Variable struct {
	Name  string
	Value Expression
	Token *Token
}

func (self Variable) CanBeNode() {}

func (self *Variable) SetValue(val Expression) {
	self.Value = val
}

func (self Variable) String() string {
	return self.Name
}

func NewVariable(name string) *Variable {
	return &Variable{name, nil, nil}
}

func NewVariableWithToken(token *Token) *Variable {
	return &Variable{token.Str, nil, token}
}
