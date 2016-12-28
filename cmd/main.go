package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/gizak/termui"
	"github.com/lazybeaver/irishrail"
)

type ByDueTime struct{ irishrail.TrainDetailSorter }

func (a ByDueTime) Less(i, j int) bool {
	return a.TrainDetailSorter[i].DueInMinutes < a.TrainDetailSorter[j].DueInMinutes
}

func DueString(td irishrail.TrainDetail) string {
	switch {
	case td.DueInMinutes <= 0:
		return fmt.Sprintf("Arriving")
	case td.DueInMinutes == 1:
		return "1 min"
	case td.DueInMinutes > 1:
		return fmt.Sprintf("%d mins", td.DueInMinutes)
	}
	return "(unknown)"
}

func DelayString(td irishrail.TrainDetail) string {
	switch {
	case td.LateByMinutes == 0:
		return ""
	case td.LateByMinutes == 1:
		return "1 min late"
	case td.LateByMinutes == -1:
		return "1 min early"
	case td.LateByMinutes > 1:
		return fmt.Sprintf("%d mins late", td.LateByMinutes)
	case td.LateByMinutes < -1:
		return fmt.Sprintf("%d mins early", -1*td.LateByMinutes)
	}
	return "(unknown)"
}

func GetDirections(traindetails []irishrail.TrainDetail) []string {
	result := []string{}
	existence := make(map[string]bool)
	for _, td := range traindetails {
		_, ok := existence[td.Direction]
		if !ok {
			existence[td.Direction] = true
			result = append(result, td.Direction)
		}
	}
	sort.Strings(result)
	return result
}

func Filter(traindetails []irishrail.TrainDetail, pred func(td irishrail.TrainDetail) bool) []irishrail.TrainDetail {
	result := traindetails[:0]
	for _, td := range traindetails {
		if pred(td) {
			result = append(result, td)
		}
	}
	return result
}

func GetRows(client irishrail.Client, station irishrail.Station) ([][]string, error) {
	traindetails, err := client.GetTrainDetails(station)
	if err != nil {
		return nil, err
	}
	result := [][]string{
		[]string{"Destination", "Due In", "Delay", "Arrival", "Origin", "Direction", "Status"},
	}
	for _, direction := range GetDirections(traindetails) {
		details := Filter(traindetails, func(td irishrail.TrainDetail) bool {
			return td.Direction == direction
		})
		result = append(result, []string{"", "", "", "", "", "", ""})
		for _, t := range details {
			result = append(result, []string{t.DestinationName, DueString(t), DelayString(t),
				t.ExpectedArrivalTime, t.OriginName, t.Direction, t.LastLocation})
		}
	}
	return result, nil
}

func DisplayUI(client irishrail.Client, station irishrail.Station) error {
	if err := termui.Init(); err != nil {
		return err
	}
	defer termui.Close()

	title := termui.NewPar("")
	title.Height = 2
	title.Border = false
	title.TextFgColor = termui.ColorGreen
	title.BorderLabel = fmt.Sprintf("%s (Last update: 0 seconds ago)", station.Name)

	table := termui.NewTable()
	table.Height = termui.TermHeight() - 3
	table.Border = false
	table.Seperator = false

	traindetails, err := GetRows(client, station)
	if err != nil {
		return err
	}
	table.Rows = traindetails
	lastupdate := time.Now()

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, title),
		),
		termui.NewRow(
			termui.NewCol(12, 0, table),
		),
	)

	termui.Body.Align()
	termui.Render(termui.Body)

	termui.Handle("/sys/kbd/<escape>", func(e termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/timer/1s", func(e termui.Event) {
		if e.Data.(termui.EvtTimer).Count%20 != 0 {
			title.BorderLabel = fmt.Sprintf("%s (Last update: %d seconds ago)", station.Name, int(time.Since(lastupdate).Seconds()))
			termui.Render(termui.Body)
			return
		}
		traindetails, err := GetRows(client, station)
		if err != nil {
			return
		}
		table.Rows = traindetails
		lastupdate = time.Now()
		termui.Render(termui.Body)
	})

	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		table.Height = termui.TermHeight() - 3
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Clear()
		termui.Render(termui.Body)
	})

	termui.Loop()
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s StationCode", os.Args[0])
	}

	client := irishrail.NewClient(irishrail.DefaultOptions)
	name := os.Args[1]
	station, err := client.LookupStation(name)
	if err != nil {
		log.Fatalf("Error lookup up station name: %s", err)
	}

	if err := DisplayUI(client, station); err != nil {
		log.Fatal("Error displaying UI: %s", err)
	}
}
