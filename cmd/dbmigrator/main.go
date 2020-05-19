package main

import (
	"context"
	"github.com/Kalinin-Andrey/dbmigrator/internal/app/cmd"
)

func main() {
	cmd.Execute(context.Background())
}
