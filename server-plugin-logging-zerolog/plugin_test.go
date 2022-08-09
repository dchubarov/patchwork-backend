package main

import (
	"testing"
	"twowls.org/patchwork/commons/logging"
)

func TestPluginInfo_ValidExtensionInterface(t *testing.T) {
	info, err := PluginInfo()
	if err != nil {
		t.Errorf("Invalid plugin info: %v", err)
	}

	ext := info.DefaultExtension()
	if ext == nil {
		t.Error("Default extension is nil")
	}

	if _, ok := ext.Interface().(logging.Facade); !ok {
		t.Error("Default extension's interface is not logging.Facade")
	}
}
