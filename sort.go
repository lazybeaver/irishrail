package irishrail

// All the type scaffolding required for clients to sort irishrail structs
// by any arbitrary element. To fully implement the sort interface, you need
// to do define the Less method on the field by which you want to sort:
//
// type ByDueTime struct {
//   TrainDetailSorter
// }
//
// func (a ByDueTime) Less(i, j int) bool {
//   return a.TrainDetailSorter[i].DueInMinutes < a.TrainDetailSorter[j].DueInMinutes
// }
//
// sort.Sort(ByDueTime(traindetails))

// Station

type StationSorter []Station

func (a StationSorter) Len() int      { return len(a) }
func (a StationSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Train

type TrainSorter []Train

func (a TrainSorter) Len() int      { return len(a) }
func (a TrainSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// StationDetail

type StationDetailSorter []StationDetail

func (a StationDetailSorter) Len() int      { return len(a) }
func (a StationDetailSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// TrainDetail

type TrainDetailSorter []TrainDetail

func (a TrainDetailSorter) Len() int      { return len(a) }
func (a TrainDetailSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
