package main

import (
    "flag"
    "fmt"
    "github.com/hygerth/beatsone"
    "os"
)

var (
    exit        = os.Exit
    usage       = "Usage: beatsone [OPTIONS]"
    options     = "Options:\n-h, -help \tPrint this help text and exit \n-v, -version \tPrint program version and exit\n" + json + schedule
    version     = "2015.08.28"
    help        = fmt.Sprintf("%s\nVersion: %s\n%s", usage, version, options)
    json        = "-j, -json\tPrint the result in JSON format\n"
    schedule    = "-s, -schedule\tPrint the schedule\n"
    cliVersion  = flag.Bool("version", false, version)
    cliHelp     = flag.Bool("help", false, help)
    cliJSON     = flag.Bool("json", false, json)
    cliSchedule = flag.Bool("schedule", false, schedule)
)

func init() {
    flag.BoolVar(cliVersion, "v", false, version)
    flag.BoolVar(cliHelp, "h", false, help)
    flag.BoolVar(cliJSON, "j", false, json)
    flag.BoolVar(cliSchedule, "s", false, schedule)
}

func main() {
    flag.Parse()

    if *cliVersion {
        fmt.Println(flag.Lookup("version").Usage)
        exit(0)
        return
    }
    if *cliHelp {
        fmt.Println(flag.Lookup("help").Usage)
        exit(0)
        return
    }
    if *cliJSON {
        if *cliSchedule {
            schedule := beatsone.GetSchedule()
            fmt.Println(schedule.JSONString())
            exit(0)
            return
        }
        np := beatsone.GetNowPlaying()
        fmt.Println(np.JSONString())
        exit(0)
        return
    }
    if *cliSchedule {
        schedule := beatsone.GetSchedule()
        fmt.Println(schedule.String())
        exit(0)
        return
    }
    np := beatsone.GetNowPlaying()
    fmt.Println(np.String())
    exit(0)
    return
}
