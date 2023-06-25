package main

import (
	"github.com/daan0220/youtube-go/my"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// main program

func main() {
	my.Migrate()
}
