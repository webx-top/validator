package validator

import (
	"sync"

	"github.com/webx-top/echo"
)

var (
	cachedValidators = make(map[string]*Validate)
	validatorsMu     sync.RWMutex
)

func Middleware(skipper ...echo.Skipper) echo.MiddlewareFunc {
	var skip echo.Skipper
	if len(skipper) > 0 {
		skip = skipper[0]
	} else {
		skip = echo.DefaultSkipper
	}
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if skip(c) {
				return h.Handle(c)
			}
			locale := c.Lang().Format(true, `_`)
			v := GetCachedValidator(c, locale)
			c.Internal().Set(`validator`, v)
			c.SetValidator(v)
			return h.Handle(c)
		})
	}
}

func GetCachedValidator(c echo.Context, locale string) *Validate {
	var transtation TranslationRegister
	transtation, locale = getTranslationRegister(locale)

	validatorsMu.RLock()
	v, ok := cachedValidators[locale]
	validatorsMu.RUnlock()
	if ok {
		return v
	}

	validatorsMu.Lock()
	defer validatorsMu.Unlock()
	// Double-check
	if v, ok = cachedValidators[locale]; ok {
		return v
	}

	v = newWithTranslationRegister(c, transtation, locale)
	cachedValidators[locale] = v
	return v
}
