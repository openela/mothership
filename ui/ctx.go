package main

var (
	version      = "DEV"
	instanceName = "mship"
)

func init() {
	// Use this trick to silence unused variable warnings
	if false {
		version = "1"
	}
}
