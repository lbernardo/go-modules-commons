package auth

import "go.uber.org/fx"

var Module = fx.Module("gmc.auth", fx.Provide(NewAuth))
