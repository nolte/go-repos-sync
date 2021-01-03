package report

import (
	"github.com/cppforlife/go-cli-ui/ui"
	uitbl "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/nolte/go-repos-sync/pkg/reposcrtl/repository"
	log "github.com/sirupsen/logrus"
)

const rowNumberStatus = 2

func RepositoryElement(elements []repository.Element) {
	log.Debugf("Generate Table Report")

	confUI := ui.NewConfUI(ui.NewNoopLogger())
	defer confUI.Flush()

	table := uitbl.Table{
		Content: "repos",

		Header: []uitbl.Header{
			uitbl.NewHeader("Remote"),
			uitbl.NewHeader("Local"),
			uitbl.NewHeader("Status"),
		},

		SortBy: []uitbl.ColumnSort{
			{Column: rowNumberStatus, Asc: true},
		},

		//Notes: []string{"(*) Currently deployed"},
	}
	for _, model := range elements {
		table.Rows = append(table.Rows, []uitbl.Value{
			uitbl.NewValueString(model.Remote.Remote.BrowserURL()),
			uitbl.NewValueString(model.Path()),
			uitbl.NewValueString(model.Status().String()),
		})
	}

	confUI.PrintTable(table)
}
