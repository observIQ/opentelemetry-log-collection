// Copyright 2021, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package operator

// DefaultRegistry is a global registry of operator types to operator builders.
var DefaultRegistry = NewRegistry()

// Registry is a registry for operators and plugins that is used for
// building types from IDs
type Registry struct {
	operators map[string]func() Builder
	plugins   map[string]func() Builder
}

// NewRegistry creates a new registry
func NewRegistry() *Registry {
	return &Registry{
		operators: make(map[string]func() Builder),
		plugins:   make(map[string]func() Builder),
	}
}

// Register will register a function to an operator type.
// This function will return a builder for the supplied type.
func (r *Registry) Register(operatorType string, newBuilder func() Builder) {
	r.operators[operatorType] = newBuilder
}

// RegisterPlugin will register a function to an plugin type.
// This function will return a builder for the supplied type.
func (r *Registry) RegisterPlugin(pluginName string, newBuilder func() Builder) {
	r.plugins[pluginName] = newBuilder
}

// Lookup looks up a given config type, prioritizing builtin operators
// before looking in registered plugins. Its second return value will
// be false if no builder is registered for that type
func (r *Registry) Lookup(configType string) (func() Builder, bool) {
	b, ok := r.operators[configType]
	if ok {
		return b, ok
	}

	mb, ok := r.plugins[configType]
	if ok {
		return mb, ok
	}

	return nil, false
}

// Register will register an operator in the default registry
func Register(operatorType string, newBuilder func() Builder) {
	DefaultRegistry.Register(operatorType, newBuilder)
}

// RegisterPlugin will register a plugin in the default registry
func RegisterPlugin(pluginName string, newBuilder func() Builder) {
	DefaultRegistry.RegisterPlugin(pluginName, newBuilder)
}

// Lookup looks up a given config type, prioritizing builtin operators
// before looking in registered plugins. Its second return value will
// be false if no builder is registered for that type
func Lookup(configType string) (func() Builder, bool) {
	return DefaultRegistry.Lookup(configType)
}
