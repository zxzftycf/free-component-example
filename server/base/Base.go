package base

import (
	"../common"
)

type Base struct {
	Config        *common.OneComponentConfig
	ComponentName string
}

func (self *Base) LoadComponent(config *common.OneComponentConfig, componentName string) {
	self.ComponentName = componentName
	self.Config = config
	return
}
