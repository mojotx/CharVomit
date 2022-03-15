package arg

const (
	// CharVomitVersion is the current version of the application
	CharVomitVersion = "v1.2.1-1-gdad925d-dirty"
)

// Version returns the current version of the application
func Version() string {
	return CharVomitVersion
}
