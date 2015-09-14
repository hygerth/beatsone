package beatsone

import (
    "bytes"
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
    s := "Title: " + np.Title + "\n"
    if len(np.Artist) > 0 {
        s += "Artist: " + np.Artist + "\n"
    }
    if len(np.Album) > 0 {
        s += "Album: " + np.Album + "\n"
    }
    s += "Artwork: " + np.Artwork
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

    // Restructure the AAC file for simpler information extration
    aacfile = RestructureAACFile(aacfile)
    // Retrieve only the information for the first song in aac file
    aacfile = splitFileIfMultipleSongs(aacfile)

    aacstring := string(aacfile)
    np.Artwork = getArtwork(aacstring)
    np.Album = getAlbum(aacstring)
    np.Artist = getArtist(aacstring)
    np.Title = getTitle(aacstring)
    return np
}

func getArtwork(s string) string {
    r := regexp.MustCompile(`artworkURL_640x\t((http|ftp|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?)`)
    m := r.FindStringSubmatch(s)
    if len(m) > 0 {
        return cleanImageURL(m[1])
    }
    return ""
}

func getAlbum(s string) string {
    r := regexp.MustCompile(`TALB(?:[^\v]+\v)+([^\t]+)\tTPE1`)
    m := r.FindStringSubmatch(s)
    if len(m) > 0 {
        return trimSpaces(m[1])
    }
    return ""
}

func getArtist(s string) string {
    r := regexp.MustCompile(`TPE1(?:[^\v]+\v)([^\t]+)\tTIT2`)
    m := r.FindStringSubmatch(s)
    if len(m) > 0 {
        return trimSpaces(m[1])
    }
    return ""
}

func getTitle(s string) string {
    r := regexp.MustCompile(`TIT2(?:[^\v]+\v)([^\t]+)\t`)
    m := r.FindStringSubmatch(s)
    if len(m) > 0 {
        return trimSpaces(m[1])
    }
    return ""
}

func splitFileIfMultipleSongs(s []byte) []byte {
    r := regexp.MustCompile(`artworkURL_640x`)
    m := r.FindAllIndex(s, 2)
    if len(m) > 1 {
        s = s[:m[len(m)-1][0]]
    }
    return s
}

func RestructureAACFile(file []byte) []byte {
    file = bytes.Replace(file, []byte{9}, []byte{32}, -1)
    file = bytes.Replace(file, []byte{0}, []byte{9}, -1)
    file = bytes.Replace(file, []byte{11}, []byte{32}, -1)
    file = bytes.Replace(file, []byte{3}, []byte{11}, -1)
    return file
}

func cleanImageURL(s string) string {
    if len(s) > 0 {
        lastjpg := strings.LastIndex(s, ".jpg")
        s = s[:lastjpg + 4]
    }
    return s
}
