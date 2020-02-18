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
