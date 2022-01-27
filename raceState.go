package main

type RaceState string

const (
	// Freshly created before both racers have confirmed a scheduled time
	RaceStateInitial RaceState = "initial"
	// Confirmed the time but before it starts
	RaceStateScheduled RaceState = "scheduled"
	// Triggered 5 minutes before starting
	RaceStateVetoCharacters RaceState = "vetoCharacters"
	// Triggered 5 minutes before starting
	RaceStateBanningCharacters RaceState = "banningCharacters"
	RaceStatePickingCharacters RaceState = "pickingCharacters"
	RaceStateBanningBuilds     RaceState = "banningBuilds"
	RaceStatePickingBuilds     RaceState = "pickingBuilds"
	RaceStateVetoBuilds        RaceState = "vetoBuilds"
	RaceStateInProgress        RaceState = "inProgress"
	// After a score is reported
	RaceStateCompleted RaceState = "completed"
)
