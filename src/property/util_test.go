package property

type PropertyStatusTestCase struct {
	Name           string
	IntValue       interface{}
	ExpectedPanic  bool
	ExpectedError  bool
	ExpectedStatus string
}
