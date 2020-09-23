package dashboard

import (
	//"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	Api "github.com/miguelrcDEV/cmpMonit/api"
	Types "github.com/miguelrcDEV/cmpMonit/types"
	"log"
	//"os"
	"strconv"
	"strings"
	"time"
)

type Header struct {
	Title *widgets.List
}

type Dashboard struct {
	Header         *ui.Grid
	InstancesTable *widgets.Table
}

var DashB Dashboard

var mediaServers = []Types.MediaServer{
	{
		Instance: "DEV",
		Url:      "https://media-server.dev.councilbox.com/api/sessions",
		User:     "OPENVIDUAPP",
		Secret:   "uincBgf9ysUCIo4MNbrfMg5hsX6FYYak",
	},
	{
		Instance: "0",
		Url:      "https://media-server.councilbox.com/api/sessions",
		User:     "OPENVIDUAPP",
		Secret:   "uincBgf9ysUCIo4MNbrfMg5hsX6FYYak",
	},
	{
		Instance: "1",
		Url:      "https://media-server1.councilbox.com/api/sessions",
		User:     "OPENVIDUAPP",
		Secret:   "uincBgf9ysUCIo4MNbrfMg5hsX6FYYak",
	},
	{
		Instance: "2",
		Url:      "https://media-server2.councilbox.com/api/sessions",
		User:     "OPENVIDUAPP",
		Secret:   "uincBgf9ysUCIo4MNbrfMg5hsX6FYYak",
	},
	{
		Instance: "3",
		Url:      "https://media-server3.councilbox.com/api/sessions",
		User:     "OPENVIDUAPP",
		Secret:   "uincBgf9ysUCIo4MNbrfMg5hsX6FYYak",
	},
	{
		Instance: "4",
		Url:      "https://media-server4.councilbox.com/api/sessions",
		User:     "OPENVIDUAPP",
		Secret:   "uincBgf9ysUCIo4MNbrfMg5hsX6FYYak",
	},
	{
		Instance: "5",
		Url:      "https://media-server5.councilbox.com/api/sessions",
		User:     "OPENVIDUAPP",
		Secret:   "uincBgf9ysUCIo4MNbrfMg5hsX6FYYak",
	},
	{
		Instance: "10",
		Url:      "https://media-server10.councilbox.com/api/sessions",
		User:     "OPENVIDUAPP",
		Secret:   "uincBgf9ysUCIo4MNbrfMg5hsX6FYYak",
	},
}

func newDashboard() Dashboard {
	var header = Header{}
	var dashboard = Dashboard{}

	termWidth, termHeight := ui.TerminalDimensions()

	header.Title = widgets.NewList()
	header.Title.Rows = []string{
		"CMP MONITOR",
	}
	header.Title.Border = false

	headerGrid := ui.NewGrid()
	headerGrid.SetRect(0, 0, termWidth, termHeight)

	headerGrid.Set(
		ui.NewRow(1.0/1,
			ui.NewCol(1.0/3, nil),
			ui.NewCol(1.0/3,
				ui.NewRow(1.0/1,
					ui.NewCol(1.0/3, nil),
					ui.NewCol(1.0/3, header.Title),
					ui.NewCol(1.0/3, nil),
				),
			),
			ui.NewCol(1.0/3, nil),
		),
	)

	dashboard.Header = headerGrid

	instancesTable := widgets.NewTable()

	instancesTable.Rows = [][]string{}

	/* HEADER TABLE ROW */
	instanceRow := []string{"Instance", "Active", "Sessions", "Total Connections"}
	instancesTable.Rows = append(instancesTable.Rows, instanceRow)

	/* TABLE BODY */
	for i := 0; i < len(Api.ActiveSessions); i++ {
		var totalConnections int = 0
		for j := 0; j < len(Api.ActiveSessions[i].Content); j++ {
			totalConnections = totalConnections + Api.ActiveSessions[i].Content[j].Connections.NumberOfElements
		}

		var sessionsString []string
		for k := 0; k < len(Api.ActiveSessions[i].Content); k++ {
			sessionsString = append(sessionsString, Api.ActiveSessions[i].Content[k].SessionId)
		}

		instanceRow := []string{Api.ActiveSessions[i].MediaServer.Instance, strconv.Itoa(Api.ActiveSessions[i].NumberOfElements), strings.Join(sessionsString, ","), strconv.Itoa(totalConnections)}
		instancesTable.Rows = append(instancesTable.Rows, instanceRow)

		if Api.ActiveSessions[i].Health == true {
			if Api.ActiveSessions[i].NumberOfElements > 0 {
				instancesTable.RowStyles[i+1] = ui.NewStyle(ui.ColorMagenta)
			} else {
				instancesTable.RowStyles[i+1] = ui.NewStyle(ui.ColorWhite)
			}
		} else {
			instancesTable.RowStyles[i+1] = ui.NewStyle(ui.ColorRed)
		}
	}

	instancesTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	instancesTable.RowSeparator = true
	instancesTable.BorderStyle = ui.NewStyle(ui.ColorCyan)
	instancesTable.FillRow = true

	dashboard.InstancesTable = instancesTable

	return dashboard
}

func newDashboardGrid(dashboard Dashboard) *ui.Grid {
	termWidth, termHeight := ui.TerminalDimensions()

	grid := ui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(0.4/3,
			ui.NewCol(1.0/1, dashboard.Header),
		),
		ui.NewRow(2.6/3,
			ui.NewCol(1.0/1, dashboard.InstancesTable),
		),
	)

	return grid
}

func Render() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	uiEvents := ui.PollEvents()

	DashB = newDashboard()
	DashBGrid := newDashboardGrid(DashB)
	//fmt.Printf("DASHBOARD %v", DashB.InstancesTable.Rows[0])
	ui.Render(DashBGrid)

	ticker := time.NewTicker(15 * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			Api.GetActiveSessions(mediaServers)
			refresh()
			//refreshDashBoardInstances()
		case <-quit:
			ticker.Stop()
			return
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
				//os.Exit(1)
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				DashBGrid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(DashBGrid)
			}
		}
	}
}

func refreshDashBoardInstances() {
	instancesTable := widgets.NewTable()
	instancesTable.Rows = [][]string{}

	/* HEADER TABLE ROW */
	instanceRow := []string{"Instance", "Active", "Sessions", "Total Connections"}
	instancesTable.Rows = append(instancesTable.Rows, instanceRow)

	/* TABLE BODY */
	for i := 0; i < len(Api.ActiveSessions); i++ {
		var totalConnections int = 0
		for j := 0; j < len(Api.ActiveSessions[i].Content); j++ {
			totalConnections = totalConnections + Api.ActiveSessions[i].Content[j].Connections.NumberOfElements
		}

		var sessionsString []string
		for k := 0; k < len(Api.ActiveSessions[i].Content); k++ {
			sessionsString = append(sessionsString, Api.ActiveSessions[i].Content[k].SessionId)
		}

		instanceRow := []string{Api.ActiveSessions[i].MediaServer.Instance, strconv.Itoa(Api.ActiveSessions[i].NumberOfElements), strings.Join(sessionsString, ","), strconv.Itoa(totalConnections)}
		instancesTable.Rows = append(instancesTable.Rows, instanceRow)

		if Api.ActiveSessions[i].NumberOfElements > 0 {
			instancesTable.RowStyles[i+1] = ui.NewStyle(ui.ColorMagenta)
		}
	}

	instancesTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	instancesTable.RowSeparator = true
	instancesTable.BorderStyle = ui.NewStyle(ui.ColorCyan)
	instancesTable.FillRow = true

	DashB.InstancesTable = instancesTable

	DashBGrid := newDashboardGrid(DashB)

	ui.Clear()
	ui.Render(DashBGrid)
}

func refresh() {
	ui.Clear()
	Render()
}
