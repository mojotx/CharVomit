package arg

// CharVomitVersion is the version reported by the application. Release builds
// set it with Go linker flags.
var CharVomitVersion = "dev"

// Version returns the current version of the application
func Version() string {
	return CharVomitVersion
}
