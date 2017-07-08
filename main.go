package main

import (
	"flag"

	"github.com/saimoon/Holmes-Totem-Dynamic/check"
	"github.com/saimoon/Holmes-Totem-Dynamic/feed"
	"github.com/saimoon/Holmes-Totem-Dynamic/lib"
	"github.com/saimoon/Holmes-Totem-Dynamic/submit"
)

func main() {
	cPath := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	ctx := &lib.Ctx{}

	err := ctx.Init(*cPath)
	if err != nil {
		panic(err.Error())
	}

	err = feed.Run(ctx, false)
	if err != nil {
		panic(err.Error())
	}

	err = check.Run(ctx, false)
	if err != nil {
		panic(err.Error())
	}

	err = submit.Run(ctx, true)
	if err != nil {
		panic(err.Error())
	}
}
