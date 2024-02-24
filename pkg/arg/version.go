package arg

const (
	// CharVomitVersion is the current version of the application
	CharVomitVersion = "v1.5.0"
)

// Version returns the current version of the application
func Version() string {
	return CharVomitVersion
}
