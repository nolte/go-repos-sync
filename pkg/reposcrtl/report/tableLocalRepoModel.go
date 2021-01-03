package report

import (
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/store"

	"github.com/cppforlife/go-cli-ui/ui"
	uitbl "github.com/cppforlife/go-cli-ui/ui/table"
	log "github.com/sirupsen/logrus"
)

func BrainElements(models []store.LocalRepoModel) {
	confUI := ui.NewConfUI(ui.NewNoopLogger())
	defer confUI.Flush()

	table := uitbl.Table{
		Content: "repos",

		Header: []uitbl.Header{
			uitbl.NewHeader("Remote Project ID"),
			uitbl.NewHeader("Brain ID"),
			uitbl.NewHeader("local path"),
			uitbl.NewHeader("Status"),
		},

		SortBy: []uitbl.ColumnSort{
			{Column: rowNumberStatus, Asc: true},
		},
	}

	for _, model := range models {
		log.Debugf("Model %s", model.Path)
		table.Rows = append(table.Rows, []uitbl.Value{
			uitbl.NewValueInt(int(model.RemoteRefer)),
			uitbl.NewValueInt(int(model.ID)),
			uitbl.NewValueString(model.Path),
			uitbl.NewValueString(model.SyncSha),
		})
	}

	confUI.PrintTable(table)
}
