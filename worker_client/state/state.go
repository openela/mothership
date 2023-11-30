package state

type PackageState struct {
	// Packages is a map of RPM path to SHA256 hash of the package.
	// The RPM path is relative to the data directory (usually /opt/srcs)
	Packages map[string]string `json:"packages"`
}

type State interface {
	Close() error
	UpdatePackageState() error
}
