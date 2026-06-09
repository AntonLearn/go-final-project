// Package format contains immutable date layout templates used for parsing,
// data-layer serialization, and user-facing representations across the task scheduler.
package format

const (
	// YYYYMMDD defines the layout used for internal date storage,
	// cross-package sorting, and database-level temporal calculations.
	YYYYMMDD = "20060102"

	// DDMMYYYY defines the localized layout used primarily for user-facing
	// date displays and visual output formatting.
	DDMMYYYY = "02.01.2006"
)
