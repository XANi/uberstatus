package uber

// Plugin interface
type Plugin interface {
	// Init() is called before main loop is run.
	// Error here will stop the main program from running so it is a good place to run any validation
	// of config or environment you need
	// Plugin should run any validation checks it needs in that function and return error if it fails, it is also good place to do stuff like precompiling templates (especially if they can be wrong from user input)
	// It is also good place to run any background goroutines plugin needs
	Init() error
	// UpdatePeriodic will run on each tick of timer in interval specified, and also will be run once to pre-populate  data
	// note that any switch-state handling (like displaying different kind of message when you click it) is still on plugin to do
	// This have to return state, even if it was previous one
	// It *can* delay for longer than interval,  engine will run at most one UpdatePeriodic(), but any long-term/heavy data gathering should be run out of separate goroutine started in Init()
	//
	UpdatePeriodic() Update
	// UpdateFromEvent is ran on each user-initiated event. If for some reason update can't be generated or doesn't make sense, return empty update
	UpdateFromEvent(Event) Update
	// GetUpdateInterval returns interval in milliseconds that will be used to run UpdatePeriodic() function by main gor
	GetUpdateInterval() int
}
