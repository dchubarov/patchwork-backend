package main

import "testing"

func TestPluginInfo(t *testing.T) {
	if info, err := PluginInfo(); err != nil {
		t.Fatalf("PluginInfo() failed: %v", err)
	} else if info.DefaultExtension() == nil {
		t.Fatal("DefaultExtension() is nil")
	}
}
