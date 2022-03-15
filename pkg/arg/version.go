package arg

const (
	// CharVomitVersion is the current version of the application
	CharVomitVersion = "v1.3.6"
)

// Version returns the current version of the application
func Version() string {
	return CharVomitVersion
}
