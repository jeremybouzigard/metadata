package metadata

// Metadata represents the metadata of an audio file.
type Metadata struct {
	Title      string
	Artist     string
	ArtistSort string
	Album      string
	AlbumSort  string
	Year       string
	Track      string
	Genre      string
	Comment    string
	Lyrics     string
	FileBase   string
	FileDir    string
}

// Service is an interface for fetching the metadata of an audio file.
type Service interface {
	Metadata(path string) (Metadata, error)
}

// Constants are used to report messages.
const (
	ParseError                = "failed to parse metadata for filepath: %s"
	UnsupportedContainerError = "unable to parse metadata for extension: %s"
)
