package arg

import (
	"github.com/rs/zerolog/log"
)

func Debug() {

	log.Debug().Msgf("Value of Config.WeakChars = %+v", Config.WeakChars)
	log.Debug().Msgf("Value of Config.Args = %+v", Config.Args)
	log.Debug().Msgf("Value of Config.UpperCase = %+v", Config.UpperCase)
	log.Debug().Msgf("Value of Config.LowerCase = %+v", Config.LowerCase)
	log.Debug().Msgf("Value of Config.Symbols = %+v", Config.Symbols)
	log.Debug().Msgf("Value of Config.Digits = %+v", Config.Digits)

}
