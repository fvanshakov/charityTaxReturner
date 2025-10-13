package migration

import (
	"charityTax/internal/database"
	"flag"
	"log"
)

func main() {
	migrationVersion := flag.Int("version", 20250911112905, "migration version")
	flag.Parse()
	var migrator *database.Migrator
	var err error
	migrator, err = database.NewMigrator("migrations")
	err = migrator.Migrate(uint(*migrationVersion))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer migrator.Close()
}
