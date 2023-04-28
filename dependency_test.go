package depkit

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestInterfaceFunc func(value bool) bool

type TestModuleInterface interface {
	GetName() string
}

type ServiceTestSuite struct {
	suite.Suite
}

type TestModule struct {
}

func (t *TestModule) GetName() string {
	return "test"
}

type TestGenericModule[T error] interface {
}

type TestGenericModuleFunc[T error] func() TestGenericModule[T]

func (test *ServiceTestSuite) SetupTest() {
	Reset()
	test.Equalf(0, len(instance.dependencies), "Expected services to be empty")
}

func (test *ServiceTestSuite) TestRegister() {
	test.Equal(0, len(instance.dependencies))
	test.Equal([]string(nil), Dependencies())

	Register[TestModuleInterface](&TestModule{})

	test.Equal(1, len(instance.dependencies))
	test.Equal([]string{"github.com/kitstack/depkit/TestModuleInterface"}, Dependencies())

	UnRegister[TestModuleInterface]()

	test.Panics(func() { Get[TestModuleInterface]() })
}

func (test *ServiceTestSuite) TestRegisterWithoutInterfaceType() {
	test.Panics(func() { Register[TestModule](TestModule{}) })
}

func (test *ServiceTestSuite) TestAlreadyRegistered() {
	// Reset
	Reset()

	Register[TestModuleInterface](&TestModule{})
	Register[TestModuleInterface](&TestModule{})

	test.Equal(1, len(instance.dependencies))
}

func (test *ServiceTestSuite) TestGet() {
	Register[TestModuleInterface](&TestModule{})
	module := Get[TestModuleInterface]()

	test.Equalf("test", module.GetName(), "Expected module name to be 'test'")
	test.Equal([]string{"github.com/kitstack/depkit/TestModuleInterface"}, Dependencies())
}

func (test *ServiceTestSuite) TestGetAsync() {
	var exec bool
	GetAfterRegister[TestModuleInterface](func(t TestModuleInterface) {
		exec = true
	})
	Register[TestModuleInterface](&TestModule{})

	test.Equal(true, exec)
}

func (test *ServiceTestSuite) TestGetAsyncAlreadyRegistered() {
	var exec bool

	Register[TestModuleInterface](&TestModule{})
	GetAfterRegister[TestModuleInterface](func(t TestModuleInterface) {
		exec = true
	})

	test.Equal(true, exec)
}

func (test *ServiceTestSuite) TestGetWithError() {
	test.Panics(func() { Get[TestModuleInterface]() })
}

func (test *ServiceTestSuite) TestRegisterFunc() {
	Register[TestInterfaceFunc](func(value bool) bool {
		return true
	})

	test.Equal(true, Get[TestInterfaceFunc]()(true))
}

func (test *ServiceTestSuite) TestRegisterFuncWithGeneric() {
	Register[TestGenericModuleFunc[error]](func() TestGenericModule[error] {
		return nil
	})

	Get[TestGenericModuleFunc[error]]()
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func BenchmarkGet(b *testing.B) {
	Reset()

	Register[TestModuleInterface](&TestModule{})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Get[TestModuleInterface]()
	}
}

func BenchmarkRegister(b *testing.B) {
	Reset()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Register[TestModuleInterface](&TestModule{})
	}
}

func BenchmarkGetAfterRegister(b *testing.B) {
	Reset()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GetAfterRegister[TestModuleInterface](func(t TestModuleInterface) {})
	}

	Register[TestModuleInterface](&TestModule{})
}
