// Package classifies commands as platform-specific (routed to gh/tea)
// or native git (passed through to git).
package classify

// platformCommands is the set of commands that are routed to the
// platform-specific CLI (gh or tea) instead of git.
var platformCommands = map[string]bool{
	"pr":       true,
	"issue":    true,
	"repo":     true,
	"workflow": true,
	"run":      true,
	"release":  true,
	"gist":     true,
	"label":    true,
	"project":  true,
	"variable": true,
	"secret":   true,
}

// IsPlatform returns true if cmd should be routed to the platform CLI
// (gh/tea) rather than to git.
func IsPlatform(cmd string) bool {
	return platformCommands[cmd]
}
