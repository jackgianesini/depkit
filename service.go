package depkit

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

var instance *Service

type Service struct {
	sync.Mutex
	services  map[string]any
	callbacks map[string][]any
}

// init : This function initializes the package by calling the Reset() function.
func init() {
	Reset()
}

// GetServiceByIdentifier : This function returns the registered service associated with the given identifier. If no service is associated with this identifier, the function returns nil.
func (s *Service) GetServiceByIdentifier(identifier string) any {
	instance.Lock()
	defer instance.Unlock()

	service, ok := instance.services[identifier]

	if !ok {
		return nil
	}

	return service
}

// GetCallbacksByIdentifier : This function returns the list of registered callbacks associated with the given identifier. If no callback is associated with this identifier, the function returns an empty list
func (s *Service) GetCallbacksByIdentifier(identifier string) []any {
	instance.Lock()
	defer instance.Unlock()

	callbacks, ok := instance.callbacks[identifier]

	if !ok {
		return []any{}
	}

	return callbacks
}

// SetServiceByIdentifier : This function registers a service associated with the given identifier.
func (s *Service) SetServiceByIdentifier(identifier string, module any) {
	instance.Lock()
	defer instance.Unlock()

	instance.services[identifier] = module
}

// SetCallbacksByIdentifier : This function adds a callback to the list of registered callbacks associated with the given identifier.
func (s *Service) SetCallbacksByIdentifier(identifier string, fn any) {
	instance.Lock()
	defer instance.Unlock()

	instance.callbacks[identifier] = append(instance.callbacks[identifier], fn)
}

// resolveIdentifier returns the identifier of the registered service associated with the given interface. If the type given as a parameter is not an interface, an error is raised.
func resolveIdentifier[T any]() string {
	var t T
	r := reflect.TypeOf(&t)

	if r.Elem().Kind() != reflect.Interface && r.Elem().Kind() != reflect.Func {
		log.Panicf("Registered need to be a interface or Func not a '(%s)'", r.Elem().Kind())
	}

	return fmt.Sprintf("%s/%s", r.Elem().PkgPath(), r.Elem().Name())
}

// Register a service associated with the given interface. If a service is already registered for this interface, the function does nothing. If callbacks have been registered for this interface, they are executed with the registered service as a parameter.
func Register[T any](module T) {

	identifier := resolveIdentifier[T]()

	if instance.GetServiceByIdentifier(identifier) != nil {
		return
	}

	instance.SetServiceByIdentifier(identifier, module)

	callbacks := instance.GetCallbacksByIdentifier(identifier)
	if len(callbacks) == 0 {
		return
	}

	for _, fn := range callbacks {
		fn.(func(T))(module)
	}
}

// Get returns the registered service associated with the given interface.
func Get[T any]() (module T) {
	return require[T]()
}

// GetAfterRegister is a special case for the init function, the callback is triggered after the service is registered. If the service has already been registered for the given interface, the callback is executed with the registered service as a parameter. If the service has not yet been registered, the callback is added to the list of registered callbacks for this interface.
func GetAfterRegister[T any](fn func(module T)) {
	identifier := resolveIdentifier[T]()

	module := instance.GetServiceByIdentifier(identifier)
	if module != nil {
		fn(module.(T))
		return
	}

	instance.SetCallbacksByIdentifier(identifier, fn)
}

// Reset resets all registered services by creating new empty maps for the services and callbacks.
func Reset() {
	instance = &Service{
		services:  make(map[string]any),
		callbacks: make(map[string][]any),
	}
}

// require returns the registered service associated with the given interface. If no service is associated with this interface, an error is raised.
func require[T any]() T {
	identifier := resolveIdentifier[T]()

	module := instance.GetServiceByIdentifier(identifier)
	if module == nil {
		log.Panicf("Service with identifier '%s' not found. Make sure to call the Register() function before attempting to retrieve the service using Get()", identifier)
	}

	return module.(T)
}
