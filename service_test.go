package servicekit

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

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

func (serviceTest *ServiceTestSuite) SetupTest() {
	Reset()
	serviceTest.Equalf(0, len(instance.services), "Expected services to be empty")
}

func (serviceTest *ServiceTestSuite) TestRegister() {
	serviceTest.Equal(0, len(instance.services))
	Register[TestModuleInterface](&TestModule{})
	serviceTest.Equal(1, len(instance.services))
}

func (serviceTest *ServiceTestSuite) TestRegisterWithoutInterfaceType() {
	serviceTest.Panics(func() { Register[TestModule](TestModule{}) })
}

func (serviceTest *ServiceTestSuite) TestAlreadyRegistered() {
	// Reset
	Reset()

	Register[TestModuleInterface](&TestModule{})
	Register[TestModuleInterface](&TestModule{})

	serviceTest.Equal(1, len(instance.services))
}

func (serviceTest *ServiceTestSuite) TestGet() {
	Register[TestModuleInterface](&TestModule{})
	module := Get[TestModuleInterface]()

	serviceTest.Equalf("test", module.GetName(), "Expected module name to be 'test'")
}

func (serviceTest *ServiceTestSuite) TestGetAsync() {
	var exec bool
	GetAfterRegister[TestModuleInterface](func(t TestModuleInterface) {
		exec = true
	})
	Register[TestModuleInterface](&TestModule{})

	serviceTest.Equal(true, exec)
}

func (serviceTest *ServiceTestSuite) TestGetAsyncAlreadyRegistered() {
	var exec bool

	Register[TestModuleInterface](&TestModule{})
	GetAfterRegister[TestModuleInterface](func(t TestModuleInterface) {
		exec = true
	})

	serviceTest.Equal(true, exec)
}

func (serviceTest *ServiceTestSuite) TestGetWithError() {
	serviceTest.Panics(func() { Get[TestModuleInterface]() })
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func BenchmarkServiceGet(b *testing.B) {
	Reset()

	Register[TestModuleInterface](&TestModule{})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Get[TestModuleInterface]()
	}
}

func BenchmarkServiceGetAfterRegister(b *testing.B) {
	Reset()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GetAfterRegister[TestModuleInterface](func(t TestModuleInterface) {})
	}

	Register[TestModuleInterface](&TestModule{})
}
