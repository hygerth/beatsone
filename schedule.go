package beatsone

import (
    "encoding/json"
    "strings"
    "time"
)

const scheduleURL = "http://fuse-music.herokuapp.com/api/programs"

type program struct {
    Programs []show
}

type show struct {
    Image string
    Start int64
    End   int64
    Title string
    URL   string
}

// Entries is a list of "Entry" objects
type Entries []Entry

// Entry is basically the same as the struct "show", with the unix time values
// converted to golang time objects
type Entry struct {
    Image string
    Start time.Time
    End   time.Time
    Title string
    URL   string
}

// JSONString returns the entries in a JSON structure as a string
func (e *Entries) JSONString() string {
    jsonobject, err := json.Marshal(e)
    checkerr(err)
    return string(jsonobject)
}

// String returns all the entries as a string separated by newlines
func (e *Entries) String() string {
    s := "Starts:\t\t\tEnds:\t\t\tTitle:\t\n"
    s += strings.Repeat("-", 80) + "\n"
    for _, entry := range *e {
        s += entry.Start.Format(layout) + "\t"
        s += entry.End.Format(layout) + "\t"
        s += entry.Title + "\n"
    }
    return s
}

func getSchedule() Entries {
    var entries Entries
    source, result := getPage(scheduleURL)
    if !result {
        return entries
    }
    var shows program
    json.Unmarshal(source, &shows)
    for i, show := range shows.Programs {
        start := time.Unix(show.Start/1000, 0)
        end := time.Unix(show.End/1000, 0)
        if i < len(shows.Programs)-1 {
            end = time.Unix(shows.Programs[i+1].Start/1000, 0).Add(-1 * time.Minute)
        }
        entry := Entry{Image: show.Image, Start: start, End: end, Title: show.Title, URL: show.URL}
        entries = append(entries, entry)
    }
    return entries
}
