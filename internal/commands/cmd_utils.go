package commands

// Contains checks if a command is in a list of commands
func Contains(commands []Command, command Command) bool {
	for _, cmd := range commands {
		if cmd == command {
			return true
		}
	}
	return false
}

