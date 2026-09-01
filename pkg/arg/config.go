package arg

// ConfigType holds the character-pool and output options shared by the CLI and CharVomit package.
type ConfigType struct {
	PasswordLen int
	Digits      bool
	ShowHelp    bool
	LowerCase   bool
	Symbols     bool
	UpperCase   bool
	WeakChars   bool
	Stdout      bool
	OutputFile  string
	Version     bool
	Excluded    string
}
