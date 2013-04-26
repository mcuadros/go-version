package version

import (
	"strings"
)

type Constraint struct {
	operator string
	version  string
}

func NewConstrain(operator, version string) *Constraint {
	constraint := new(Constraint)
	constraint.SetOperator(operator)
	constraint.SetVersion(version)

	return constraint
}

func NewConstrainFromString(name string) *Constraint {
	constraint := new(Constraint)

	return constraint
}

func (self *Constraint) SetOperator(operator string) {
	self.operator = operator
}

func (self *Constraint) GetOperator() string {
	return self.operator
}

func (self *Constraint) SetVersion(version string) {
	self.version = version
}

func (self *Constraint) GetVersion() string {
	return self.version
}

func (self *Constraint) Match(version string) bool {
	return Compare(version, self.version, self.operator)
}

func (self *Constraint) String() string {
	return strings.Trim(self.operator+" "+self.version, " ")
}
