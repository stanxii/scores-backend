package migrations

import "github.com/raphi011/scores/migrate"

var (
	MigrationSet = []migrate.Migration{V1, V2, V3, V4, V5}
	ResetSet     = []migrate.TableNames{ResetV1, ResetV2, ResetV3, ResetV4, ResetV5}
)
