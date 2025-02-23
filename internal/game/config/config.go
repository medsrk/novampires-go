// internal/game/config/config.go
package config

// DisplayConfig contains display-related settings
type DisplayConfig struct {
	Width         int
	Height        int
	Fullscreen    bool
	VSync         bool
	TargetFPS     int
	ShowFPS       bool
	ShowDebugInfo bool
}

// DefaultDisplay returns sensible display defaults
func DefaultDisplay() DisplayConfig {
	return DisplayConfig{
		Width:         1600,
		Height:        900,
		Fullscreen:    false,
		VSync:         true,
		TargetFPS:     60,
		ShowFPS:       true,
		ShowDebugInfo: true,
	}
}

// AudioConfig contains audio-related settings
type AudioConfig struct {
	MasterVolume      float64
	MusicVolume       float64
	SFXVolume         float64
	VoiceVolume       float64
	MuteWhenFocusLost bool
}

// DefaultAudio returns sensible audio defaults
func DefaultAudio() AudioConfig {
	return AudioConfig{
		MasterVolume:      1.0,
		MusicVolume:       0.7,
		SFXVolume:         0.8,
		VoiceVolume:       1.0,
		MuteWhenFocusLost: true,
	}
}

// GameplayConfig contains gameplay-related settings
type GameplayConfig struct {
	Difficulty        int
	AutoAimStrength   float64
	CameraShakeAmount float64
	ScreenShake       bool
	HitMarkers        bool
	DamageNumbers     bool
}

// DefaultGameplay returns sensible gameplay defaults
func DefaultGameplay() GameplayConfig {
	return GameplayConfig{
		Difficulty:        1, // 0=Easy, 1=Normal, 2=Hard
		AutoAimStrength:   0.6,
		CameraShakeAmount: 0.7,
		ScreenShake:       true,
		HitMarkers:        true,
		DamageNumbers:     true,
	}
}

// Config aggregates all configuration categories
type Config struct {
	Display  DisplayConfig
	Audio    AudioConfig
	Gameplay GameplayConfig
}

// Default returns a default configuration
func Default() Config {
	return Config{
		Display:  DefaultDisplay(),
		Audio:    DefaultAudio(),
		Gameplay: DefaultGameplay(),
	}
}
