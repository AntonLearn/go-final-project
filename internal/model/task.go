// Package model defines the core domain entities and data structures
// utilized throughout the task scheduler application lifecycle.
package model

// Task represents a scheduled item within the system, encapsulating its
// unique identifier, execution timeline, metadata, and optional recurrence rules.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
