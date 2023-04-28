package depkit

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

var instance *dependency

type dependency struct {
	sync.Mutex
	dependencies map[string]any
	callbacks    map[string][]any
}

// init : This function initializes the package by calling the Reset() function.
func init() {
	Reset()
}

// GetDependencyByIdentifier : This function returns the registered dependency associated with the given identifier. If no dependency is associated with this identifier, the function returns nil.
func (s *dependency) GetDependencyByIdentifier(identifier string) any {
	instance.Lock()
	defer instance.Unlock()

	dependency, ok := instance.dependencies[identifier]

	if !ok {
		return nil
	}

	return dependency
}

// GetCallbacksByIdentifier : This function returns the list of registered callbacks associated with the given identifier. If no callback is associated with this identifier, the function returns an empty list
func (s *dependency) GetCallbacksByIdentifier(identifier string) []any {
	instance.Lock()
	defer instance.Unlock()

	callbacks, ok := instance.callbacks[identifier]

	if !ok {
		return []any{}
	}

	return callbacks
}

// SetDependencyByIdentifier : This function registers a dependency associated with the given identifier.
func (s *dependency) SetDependencyByIdentifier(identifier string, module any) {
	instance.Lock()
	defer instance.Unlock()

	instance.dependencies[identifier] = module
}

// SetCallbacksByIdentifier : This function adds a callback to the list of registered callbacks associated with the given identifier.
func (s *dependency) SetCallbacksByIdentifier(identifier string, fn any) {
	instance.Lock()
	defer instance.Unlock()

	instance.callbacks[identifier] = append(instance.callbacks[identifier], fn)
}

// resolveIdentifier returns the identifier of the registered dependency associated with the given interface. If the type given as a parameter is not an interface, an error is raised.
func resolveIdentifier[T any]() string {
	var t T
	r := reflect.TypeOf(&t)

	if r.Elem().Kind() != reflect.Interface && r.Elem().Kind() != reflect.Func {
		log.Panicf("register method need to be a interface or Func not a '(%s)'", r.Elem().Kind())
	}

	return fmt.Sprintf("%s/%s", r.Elem().PkgPath(), r.Elem().Name())
}

// Register a dependency associated with the given interface. If a dependency is already registered for this interface, the function does nothing. If callbacks have been registered for this interface, they are executed with the registered dependency as a parameter.
func Register[T any](module T) {

	identifier := resolveIdentifier[T]()

	if instance.GetDependencyByIdentifier(identifier) != nil {
		return
	}

	instance.SetDependencyByIdentifier(identifier, module)

	callbacks := instance.GetCallbacksByIdentifier(identifier)
	if len(callbacks) == 0 {
		return
	}

	for _, fn := range callbacks {
		fn.(func(T))(module)
	}
}

// UnRegister removes the dependency associated with the given interface. If no dependency is registered for this interface, the function does nothing.
func UnRegister[T any]() {
	identifier := resolveIdentifier[T]()

	instance.Lock()
	defer instance.Unlock()

	delete(instance.dependencies, identifier)
	delete(instance.callbacks, identifier)
}

// Get returns the registered dependency associated with the given interface.
func Get[T any]() (module T) {
	return require[T]()
}

// Dependencies returns the list of registered dependencies.
func Dependencies() []string {
	var dependencies []string
	for key := range instance.dependencies {
		dependencies = append(dependencies, key)
	}
	return dependencies
}

// GetAfterRegister is a special case for the init function, the callback is triggered after the dependency is registered. If the dependency has already been registered for the given interface, the callback is executed with the registered dependency as a parameter. If the dependency has not yet been registered, the callback is added to the list of registered callbacks for this interface.
func GetAfterRegister[T any](fn func(module T)) {
	identifier := resolveIdentifier[T]()

	module := instance.GetDependencyByIdentifier(identifier)
	if module != nil {
		fn(module.(T))
		return
	}

	instance.SetCallbacksByIdentifier(identifier, fn)
}

// Reset resets all registered dependencies by creating new empty maps for the dependencies and callbacks.
func Reset() {
	instance = &dependency{
		dependencies: make(map[string]any),
		callbacks:    make(map[string][]any),
	}
}

// require returns the registered dependency associated with the given interface. If no dependency is associated with this interface, an error is raised.
func require[T any]() T {
	identifier := resolveIdentifier[T]()

	module := instance.GetDependencyByIdentifier(identifier)
	if module == nil {
		log.Panicf("dependency with identifier '%s' not found. Make sure to call the Register() function before attempting to retrieve the dependency using Get()", identifier)
	}

	return module.(T)
}
