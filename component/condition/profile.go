package condition

import (
	"github.com/codnect/procyoncore/component"
	"github.com/codnect/procyoncore/component/filter"
	"github.com/codnect/procyoncore/runtime"
)

type OnProfileCondition struct {
	profiles []string
}

func OnProfile(profiles ...string) *OnProfileCondition {
	return &OnProfileCondition{
		profiles: profiles,
	}
}

func (c *OnProfileCondition) MatchesCondition(ctx component.ConditionContext) bool {
	result, err := ctx.Container().GetObject(ctx, filter.ByTypeOf[runtime.Environment]())
	if err != nil {
		//
		return false
	}

	environment := result.(runtime.Environment)

	for _, profile := range c.profiles {
		if !environment.IsProfileActive(profile) {
			return false
		}
	}

	return true
}
