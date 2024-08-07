package component

import (
	"codnect.io/reflector"
	"fmt"
	"github.com/codnect/procyoncore/component/filter"
	"sync"
)

var (
	components   = make(map[string]*Component)
	muComponents = sync.RWMutex{}
)

type Component struct {
	definition *Definition
	conditions []Condition
}

func createComponent(constructor Constructor, options ...Option) *Component {
	definition, err := MakeDefinition(constructor, options...)

	if err != nil {
		panic(err)
	}

	return &Component{
		definition: definition,
		conditions: make([]Condition, 0),
	}
}

func (c *Component) Definition() *Definition {
	return c.definition
}

func (c *Component) Conditions() []Condition {
	copyOfConditions := make([]Condition, 0)

	for _, condition := range copyOfConditions {
		copyOfConditions = append(copyOfConditions, condition)
	}

	return copyOfConditions
}

type Registration struct {
	component *Component
}

func (r Registration) ConditionalOn(condition Condition) Registration {
	if condition != nil {
		r.component.conditions = append(r.component.conditions, condition)
	}

	return r
}

func Register(constructor Constructor, options ...Option) Registration {
	defer muComponents.Unlock()
	muComponents.Lock()

	component := createComponent(constructor, options...)
	componentName := component.Definition().Name()

	if _, exists := components[componentName]; exists {
		panic(fmt.Sprintf("component with name '%s' already exists", componentName))
	}

	components[componentName] = component

	return Registration{
		component: component,
	}
}

func List(filters ...filter.Filter) []*Component {
	defer muComponents.Unlock()
	muComponents.Lock()

	filterOpts := filter.Of(filters...)
	componentList := make([]*Component, 0)

	for _, component := range components {
		definition := component.Definition()

		if filterOpts.Name != "" && filterOpts.Name != component.Definition().Name() {
			continue
		}

		if filterOpts.Type == nil {
			componentList = append(componentList, component)
			continue
		}

		if definition.Type().CanConvert(filterOpts.Type) {
			componentList = append(componentList, component)
		} else if reflector.IsPointer(definition.Type()) && !reflector.IsPointer(filterOpts.Type) && !reflector.IsInterface(filterOpts.Type) {
			pointerType := reflector.ToPointer(definition.Type())

			if pointerType.Elem().CanConvert(filterOpts.Type) {
				componentList = append(componentList, component)
			}
		}
	}

	return componentList
}
