package mongodb

import "go.uber.org/fx"

var Module = fx.Module("gmc.database.mongodb",
	fx.Provide(
		NewConnection,
		NewDatabase,
	))
