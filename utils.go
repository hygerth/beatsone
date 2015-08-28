package beatsone

import (
    "io/ioutil"
    "log"
    "net/http"
    "regexp"
    "strings"
)

const useragent = "mozilla/5.0 (iphone; cpu iphone os 7_0_2 like mac os x) applewebkit/537.51.1 (khtml, like gecko) version/7.0 mobile/11a501 safari/9537.53"

func trimSpaces(s string) string {
    s = strings.Trim(s, " ")
    re := regexp.MustCompile(`\s{2,}`)
    return re.ReplaceAllString(s, " ")
}

func getPage(url string) ([]byte, bool) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    checkerr(err)
    req.Header.Set("User-Agent", useragent)
    resp, err := client.Do(req)
    if err != nil || resp.StatusCode != 200 {
        return []byte{}, false
    }
    defer resp.Body.Close()
    b, _ := ioutil.ReadAll(resp.Body)
    return b, true
}

func checkerr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
