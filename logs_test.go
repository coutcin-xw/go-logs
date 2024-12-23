package logs

import (
	"testing"
)

func TestLogger_Console(t *testing.T) {
	Log.SetFile("1.txt")
	Log.SetIsLogToFile(true)
	defer Log.Close(false)
	Log.SetLevel(Debug)
	Log.Important("test")
	Log.SetColor(true)
	Log.Importantf("%stest%s", "aaa", "sd")
	Log.Info("info")
	Log.Hint("hint")
	Log.Important("important")
	Log.Debug("debug")
	Log.Warn("warn")
	Log.Error("Error")
	// Log.Info("%stest", "1")
	AddLevel(1, "test")
	Log.Log(1, "test")
}
