package beatsone

// GetNowPlaying collects the information about what song is currently playing on Beats1
func GetNowPlaying() NowPlaying {
    return getNowPlaying()
}

// GetSchedule collects the schedule for Beats1
func GetSchedule() Entries {
    return getSchedule()
}
