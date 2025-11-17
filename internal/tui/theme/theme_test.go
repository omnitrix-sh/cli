package theme

import (
	"testing"
)

func TestThemeRegistration(t *testing.T) {
	// Get list of available themes
	availableThemes := AvailableThemes()
	
	// Check if only "dracula" theme is registered
	if len(availableThemes) != 1 {
		t.Errorf("Expected exactly 1 theme, got %d", len(availableThemes))
	}
	
	if availableThemes[0] != "dracula" {
		t.Errorf("Expected dracula theme, got %s", availableThemes[0])
	}
	
	// Try to get the dracula theme and make sure it's not nil
	dracula := GetTheme("dracula")
	if dracula == nil {
		t.Errorf("Dracula theme is nil")
	}
	
	// Test current theme should be dracula
	if CurrentThemeName() != "dracula" {
		t.Errorf("Current theme should be dracula, got %s", CurrentThemeName())
	}
	
	// Test setting theme to dracula should work
	err := SetTheme("dracula")
	if err != nil {
		t.Errorf("Failed to set theme to dracula: %v", err)
	}
	
	// Test setting theme to non-existent theme should fail
	err = SetTheme("nonexistent")
	if err == nil {
		t.Errorf("Setting non-existent theme should have failed")
	}
}
