package main

import (
	"cicdgf/internal/boot"
	_ "cicdgf/internal/boot"
	_ "cicdgf/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"cicdgf/internal/cmd"
)

func main() {
	boot.InitView()
	cmd.Main.Run(gctx.GetInitCtx())
}
