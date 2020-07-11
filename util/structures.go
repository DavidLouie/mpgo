package util

type Artist struct {
	Id   int
	Name string
}

func NewArtist() Artist {
	a := Artist{}
	a.Id = -1
	a.Name = ""
	return a
}

type Album struct {
	Id       int
	Title    string
	Genre    string
	Year     int
	ArtistId int
}

func NewAlbum() Album {
	a := Album{}
	a.Id = -1
	a.Title = ""
	a.Genre = ""
	a.Year = -1
	a.ArtistId = -1
	return a
}

type Song struct {
	Id       int
	Title    string
	Duration int
	Size     int
	Track    int
	Path     string
	BitRate  int
	Ext      string
	AlbumId  int
}

func NewSong() Song {
	s := Song{}
	s.Id = -1
	s.Title = ""
	s.Duration = -1
	s.Size = -1
	s.Track = -1
	s.Path = ""
	s.BitRate = -1
	s.Ext = ""
	s.AlbumId = -1
	return s
}
