package subsonic

import (
    "encoding/xml"
)

type SubResp struct {
    XMLName         xml.Name        `xml:"subsonic-response"`
    Status          string          `xml:"status,attr"`
    Version         string          `xml:"version,attr"`
    MusicFolders    *MusicFolders   `xml:"musicFolders,omitempty"i`
}

type MusicFolder struct {
    Id      string `xml:"id,attr"`
    Name    string `xml:"name,attr"`
}

type MusicFolders struct {
    Folders []MusicFolder `xml:"musicFolder"`
}
