# sealer

Create encryption for texts using sealed

> github.com/lbernardo/go-modules-commons/configs  github.com/lbernardo/go-modules-commons/logger **are required**

### Import module

```go
package main

import (
	"github.com/lbernardo/go-modules-commons/configs"
	"github.com/lbernardo/go-modules-commons/logger"
	"github.com/lbernardo/go-modules-commons/services/sealer"
)

func main() {
	app := fx.New(
		configs.Module, // is required
		logger.Module, // is required
	sealer.Module,
)
	app.Run()
}
```

## Usage

```go
package service

import (
	"context"
	"fmt"
	"github.com/lbernardo/go-modules-commons/services/sealer"
)

type MyService struct {
	sealerProvider *sealer.SealerProvider
}

func NewMyService(s *sealer.SealerProvider) *MyService {
	return &MyService{
		sealerProvider: s,
	}
}

func (c *MyService) Run() {

	content, err := c.sealerProvider.Seal(context.Background(), "example test")
	if err != nil {
		fmt.Println("error",err)
		return
	}
	fmt.Println("text sealed", content)
	result, err := c.sealerProvider.Unseal(content)
	if err != nil {
		fmt.Println("err", err)
		return
    }
	fmt.Println("text unsealed", result)

}
```