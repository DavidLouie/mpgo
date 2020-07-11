package subsonic

import (
	"encoding/xml"
)

type subResp struct {
	XMLName      xml.Name      `xml:"subsonic-response"`
	Status       string        `xml:"status,attr"`
	Version      string        `xml:"version,attr"`
	ErrorCode    *errorCode    `xml:"error,omitempty"`
	MusicFolders *musicFolders `xml:"musicFolders,omitempty"i`
	Genres       *genres       `xml:"genres,omitempty"i`
}

type errorCode struct {
	Code    string `xml:"code,attr"`
	Message string `xml:"message,attr"`
}

type musicFolder struct {
	Id   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type musicFolders struct {
	Folders []musicFolder `xml:"musicFolder"`
}

type directory struct {
	Children []dirChild `xml:"child"`
}

type dirChild struct {
	Id         int    `xml:"id,attr"`
	Parent     int    `xml:"parent,attr"`
	Title      string `xml:"title,attr"`
	IsDir      bool   `xml:"isDir,attr"`
	Album      string `xml:"album,attr,omitempty"`
	Artist     string `xml:"artist,attr"`
	Track      int    `xml:"track,attr,omitempty"`
	Year       int    `xml:"year,attr,omitempty"`
	Genre      string `xml:"genre,attr,omitempty"`
	CoverArt   string `xml:"coverArt,attr"`
	Size       int    `xml:"size,attr,omitempty"`
	ContType   string `xml:"contentType,attr,omitempty"`
	Suffix     string `xml:"suffix,attr,omitempty"`
	TrContType string `xml:"transcodedContentType,attr,omitempty"`
	TrSuffix   string `xml:"transcodedSuffix,attr,omitempty"`
	Duration   int    `xml:"duration,attr,omitempty"`
	BitRate    int    `xml:"bitRate,attr,omitempty"`
	Path       string `xml:"path,attr,omitempty"`
}

type genre struct {
	SongCount  string `xml:"songCount,attr"`
	AlbumCount string `xml:"albumCount,attr"`
	Genre      string `xml:",chardata"`
}

type genres struct {
	Genres []genre `xml:"genre"`
}
