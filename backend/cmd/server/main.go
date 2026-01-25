package main

import (
	"context"
	"os"

	"github.com/aereal/nikki/backend/entrypoint"
)

func main() { os.Exit(entrypoint.Run(entrypoint.NewDevEntrypoint(context.Background()))) }
