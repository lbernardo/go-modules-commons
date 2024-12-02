package sealer

// github.com/lbernardo/go-modules-commons/configs & github.com/lbernardo/go-modules-commons/logger are required

import "go.uber.org/fx"

var Module = fx.Module("gmc.sealer-module", fx.Provide(
	NewSealerProvider,
))
