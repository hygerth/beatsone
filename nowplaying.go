package beatsone

import (
    "encoding/json"
    "regexp"
    "strings"
)

const playlistURL = "http://itsliveradiobackup.apple.com/streams/hub02/session02/64k/"
const layout = "2006-01-02 15:04"

// NowPlaying describes the structure for the metadata of a song
type NowPlaying struct {
    Artwork string
    Album   string
    Artist  string
    Title   string
}

// JSONString returns the song playing in a JSON structure as a string
func (np *NowPlaying) JSONString() string {
    jsonobject, err := json.Marshal(np)
    checkerr(err)
    return string(jsonobject)
}

// String returns the song playing as a string separated by newlines
func (np *NowPlaying) String() string {
    s := "Artwork: " + np.Artwork + "\n"
    if len(np.Album) > 0 {
        s += "Album: " + np.Album + "\n"
    }
    if len(np.Artist) > 0 {
        s += "Artist: " + np.Artist + "\n"
    }
    s += "Title: " + np.Title
    return s
}

func getNowPlaying() NowPlaying {
    var np NowPlaying
    progplaylist, result := getPage(playlistURL + "prog.m3u8")
    if !result {
        return np
    }
    m3u8 := string(progplaylist)
    lines := strings.Split(m3u8, "\n")
    if len(lines) < 2 {
        return np
    }
    lastFile := lines[len(lines)-2]
    aacfile, success := getPage(playlistURL + lastFile)
    if !success {
        return np
    }
    file := string(aacfile)
    cleanFile := removeNonASCIICharsAndClean(file)
    // Retrieve only the information for the first song in file
    cleanFile = splitFile(cleanFile)

    np.Artwork = getArtwork(cleanFile)
    np.Album = getAlbum(cleanFile)
    np.Artist = getArtist(cleanFile)
    np.Title = getTitle(cleanFile)
    return np
}

func getArtwork(s string) string {
    r := regexp.MustCompile(`artworkURL_640x ((http|ftp|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?)`)
    m := r.FindStringSubmatch(s)
    if len(m) > 0 {
        return cleanImageURL(m[1])
    }
    return ""
}

func getAlbum(s string) string {
    r := regexp.MustCompile(`TALB(?:[^\t]+\t)+([\S ]+) TPE1`)
    m := r.FindStringSubmatch(s)
    if len(m) > 0 {
        return trimSpaces(m[1])
    }
    return ""
}

func getArtist(s string) string {
    r := regexp.MustCompile(`TPE1(?:[^\t]+\t)([\S ]+) TIT2`)
    m := r.FindStringSubmatch(s)
    if len(m) > 0 {
        return trimSpaces(m[1])
    }
    return ""
}

func getTitle(s string) string {
    r := regexp.MustCompile(`TIT2(?:[^\t]+\t)([\w \(\)\&\.\,\-\'\"\_]+) \\`)
    m := r.FindStringSubmatch(s)
    if len(m) > 0 {
        return trimSpaces(m[1])
    }
    return ""
}

func splitFile(s string) string {
    r := regexp.MustCompile(`artworkURL_640x`)
    m := r.FindAllStringIndex(s, 2)
    if len(m) > 1 {
        s = s[:m[len(m)-1][0]]
    }
    return s
}

func removeNonASCIICharsAndClean(s string) string {
    regex := regexp.MustCompile(`[^\x00-\x7F]+`)
    s = regex.ReplaceAllStringFunc(s, func(w string) string {
        return ""
    })
    var res string
    for i := 0; i < len(s); i++ {
        if (s[i] >= 32 && s[i] < 128) || s[i] == 10 {
            res += string(s[i])
        } else if s[i] == 0 {
            res += " "
        } else if s[i] == 3 {
            res += "\t"
        } else if s[i] == 12 {
            res += "\n"
        }
    }
    return res
}

func cleanImageURL(s string) string {
    if len(s) > 0 {
        s = strings.SplitAfter(s, ".jpg")[0]
    }
    return s
}
