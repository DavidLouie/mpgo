package subsonic

import (
    "encoding/xml"
)

type subResp struct {
    XMLName         xml.Name        `xml:"subsonic-response"`
    Status          string          `xml:"status,attr"`
    Version         string          `xml:"version,attr"`
    ErrorCode       *errorCode      `xml:"error,omitempty"`
    MusicFolders    *musicFolders   `xml:"musicFolders,omitempty"i`
    Genres          *genres         `xml:"genres,omitempty"i`
}

type errorCode struct {
    Code      string `xml:"code,attr"`
    Message   string `xml:"message,attr"`
}

type musicFolder struct {
    Id      string `xml:"id,attr"`
    Name    string `xml:"name,attr"`
}

type musicFolders struct {
    Folders []musicFolder `xml:"musicFolder"`
}

type genre struct  {
    SongCount  string `xml:"songCount,attr"`
    AlbumCount string `xml:"albumCount,attr"`
    Genre      string `xml:",chardata"`
}

type genres struct {
    Genres []genre `xml:"genre"`
}
