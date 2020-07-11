package util

// Artist is an artist with an ID and a name.
type Artist struct {
	ID   int
	Name string
}

// NewArtist default initializes Artist.
func NewArtist() Artist {
	a := Artist{}
	a.ID = -1
	a.Name = ""
	return a
}

// Album is an album with some metadata.
type Album struct {
	ID       int
	Title    string
	Genre    string
	Year     int
	ArtistID int
}

// NewAlbum default initializes Album.
func NewAlbum() Album {
	a := Album{}
	a.ID = -1
	a.Title = ""
	a.Genre = ""
	a.Year = -1
	a.ArtistID = -1
	return a
}

// Song is a song with related metadata.
type Song struct {
	ID       int
	Title    string
	Duration int
	Size     int
	Track    int
	Path     string
	BitRate  int
	Ext      string
	AlbumID  int
}

// NewSong default initializes Song.
func NewSong() Song {
	s := Song{}
	s.ID = -1
	s.Title = ""
	s.Duration = -1
	s.Size = -1
	s.Track = -1
	s.Path = ""
	s.BitRate = -1
	s.Ext = ""
	s.AlbumID = -1
	return s
}
