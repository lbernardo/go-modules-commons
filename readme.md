# commons

```shell
# install service
go get github.com/lbernardo/go-modules-commons
go get github.com/lbernardo/go-modules-commons/services/<service>
```

## Services

- apiserver
- database
- metrics
- auth
- rabbitmq
- sealer


## Requirements

```go
package main

import (
	"github.com/lbernardo/go-modules-commons/configs"
	"github.com/lbernardo/go-modules-commons/logger"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		configs.Module, // configs is required
		logger.Module, // logger is required
	)
	app.Run()
}
```

### Configs

configs is a required module used to read application configurations. We use spf13/viper to read configuration files.
When the module is instantiated it will try to find the following files
`config.dev.yaml` or `/etc/secrets/config.yaml`, or you can set the configuration file as the environment variable: `APP_CONFIG_FILE`

the configuration file follows the following parameters

```yaml
app:
  debug: false # when enable the logger is enabled as DEBUG MODE
  auth: # auth configurations
    secret: XXXX_DDDD_FFF # auth secret to JWT
  database: # database configurations
    mongodb: # mongodb configurations
      uri: mongodb://localhost:27017 # mongo URI
      name: mydb # database name
  queue: # rabbitmq configurations
    host: <HOST_WITH_PASSWORD> # rabbitmq host
```

> You can set any configurations and read with `*viper.Viper`