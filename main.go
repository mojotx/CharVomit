package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mojotx/CharVomit/pkg/CharVomit"
	"github.com/mojotx/CharVomit/pkg/arg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TO-DO:
// - Add support for duplicate character checking
// - Add support for specifying valid characters, e.g., all upper-case, etc.
func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	arg.Parse()

	var cv CharVomit.CharVomit

	if err := cv.SetAcceptableChars(arg.Config); err != nil {
		log.Fatal().Err(err).Msg("could not set acceptable chars")
		os.Exit(1)
	}

	pwLen := 32

	log.Debug().Msgf("debug: flag.NArg() = %d, Config.Args = %d", flag.NArg(), arg.Config.Args)

	if flag.NArg() == 1 {

		var err error
		pwLen, err = strconv.Atoi(flag.Arg(0))
		if err != nil {
			log.Fatal().Err(err).Msgf("cannot parse argument '%s'", flag.Arg(0))
			os.Exit(1)
		}

		// Get absolute value
		if pwLen < 0 {
			pwLen = pwLen * -1
		}
	}

	log.Debug().Msgf("pwLen = %d", pwLen)
	log.Debug().Msgf("flag.NArg() = %d", flag.NArg())
	log.Debug().Msgf("Config.Args = %d", arg.Config.Args)
	arg.Debug()

	pw, err := cv.Puke(pwLen)
	if err != nil {
		log.Fatal().Err(err).Msgf("error calling Puke(%d)", pwLen)
		os.Exit(1)
	}

	log.Debug().Msgf("pw = '%s'", pw)
	fmt.Println(pw)
}
