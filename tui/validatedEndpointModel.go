package tui

import (
	"fmt"
	"github.com/buonotti/apisense/v2/log"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/buonotti/apisense/v2/errors"
	"github.com/buonotti/apisense/v2/validation/pipeline"
)

var (
	selectedValidatedEndpoint pipeline.ValidatedEndpoint
	allowEndpointSelection    bool
	resultRows                []table.Row
	updateResultRows          = false
)

type validationEndpointModel struct {
	keymap      keymap
	table       table.Model
	resultModel tea.Model
}

func ValidationEndpointModel() tea.Model {
	t := table.New(
		table.WithColumns(getValidatedEndpointColumns()),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#F38BA8")).
		Background(lipgloss.Color("#1e1e2e")).
		Bold(false)
	t.SetStyles(s)

	return validationEndpointModel{
		keymap:      DefaultKeyMap,
		table:       t,
		resultModel: ResultModel(),
	}
}

func (v validationEndpointModel) Init() tea.Cmd {
	return nil
}

func (v validationEndpointModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmdModel tea.Cmd
	if choiceReportModel != "reportModel" {
		t := table.New(
			table.WithColumns(getValidatedEndpointColumns()),
			table.WithRows(validatedEndpointRows),
			table.WithFocused(true),
			table.WithHeight(7),
		)
		s := table.DefaultStyles()
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("#F38BA8")).
			Background(lipgloss.Color("#1e1e2e")).
			Bold(false)
		t.SetStyles(s)
		v.table = t
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, v.keymap.back):
				if choiceReportModel == "validatedEndpointModel" {
					choiceReportModel = "reportModel"
				}
			case key.Matches(msg, v.keymap.quit):
				return v, tea.Quit
			case key.Matches(msg, v.keymap.choose):
				if allowEndpointSelection {
					i, err := strconv.Atoi(v.table.SelectedRow()[0])
					if err != nil {
						log.TuiLogger().Fatal(err)
					}
					val, err := getSelectedValidatedEndpoint(selectedReport, i)
					if err != nil {
						log.TuiLogger().Fatal(err)
					}
					selectedValidatedEndpoint = val
					resultRows = getResultRows(selectedValidatedEndpoint.TestCaseResults)
					if choiceReportModel != "validatedEndpointModel" {
						v.resultModel, cmdModel = v.resultModel.Update(msg)
						v.table, cmd = v.table.Update(msg)
						return v, tea.Batch(cmd, cmdModel)
					}
					choiceReportModel = "resultModel"
					updateResultRows = true
					v.table, cmd = v.table.Update(msg)
				}
				return v, cmd
			}
		}
		v.resultModel, cmdModel = v.resultModel.Update(msg)
		v.table, cmd = v.table.Update(msg)
		return v, tea.Batch(cmd, cmdModel)
	}
	return v, nil
}

func (v validationEndpointModel) View() string {
	if choiceReportModel != "validatedEndpointModel" {
		return lipgloss.NewStyle().Render(v.resultModel.View())
	}
	return lipgloss.NewStyle().Render(v.table.View())
}

func getValidatedEndpointRows(validatedEndpoint pipeline.Report) []table.Row {
	rows := make([]table.Row, 0)
	for i, point := range validatedEndpoint.Endpoints {
		rows = append(rows, table.Row{fmt.Sprintf("%v", i), point.EndpointName})
	}
	allowEndpointSelection = true
	if len(rows) < 1 {
		rows = append(rows, table.Row{"", "No endpoints found"})
		allowEndpointSelection = false
	}

	return rows
}

func getValidatedEndpointColumns() []table.Column {
	return []table.Column{
		{Title: "", Width: 3},
		{Title: "Endpoint name", Width: 73},
	}
}

func getSelectedValidatedEndpoint(report pipeline.Report, index int) (pipeline.ValidatedEndpoint, error) {
	if index > len(report.Endpoints) || index < 0 {
		return pipeline.ValidatedEndpoint{}, errors.ModelError.New("Index out of range")
	}
	return report.Endpoints[index], nil
}
