package parser

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jeremybouzigard/metadata"
)

// An Atom is the basic unit of an MPEG-4 file. MPEG-4 is the container used by
// .m4a files.
type Atom struct {
	Start    int64  // byte index of the start of the atom in the file
	End      int64  // byte index of the end of the atom in the file
	Size     int    // number of bytes in the atom
	Type     string // four-character code that identifies the atom type
	Contents []byte // contents of the atom from start to end
}

// Constants are used to parse audio files and report messages.
const (
	SizeLen                  = 4 // length of size field that specifies the number of bytes in atom
	TypeLen                  = 4 // length of type field that specifies format of the data in atom
	ReadAtomAtError          = "failed to read %s for atom %s at %x: %v"
	ReadSiblingWithTypeError = "failed to read %s for sibling of %s atom with target type %s: %v"
	ReadSiblingEror          = "failed to read %s for sibling of %s atom: %v"
)

// M4A parses an .m4a file and returns metadata.
func M4A(path string) (metadata.Metadata, error) {
	var m metadata.Metadata

	fi, err := os.Stat(path)
	if err != nil {
		return m, err
	}

	f, err := os.OpenFile(path, os.O_RDONLY, fi.Mode())
	if err != nil {
		return m, err
	}
	defer f.Close()

	if isM4AContainer(f) != true {
		return m, fmt.Errorf("file does not use M4A container")
	}

	ftyp, err := readAtomAt(f, 0)
	if err != nil {
		return m, err
	}

	moov, err := readSiblingWithType(f, ftyp, "moov")
	if err != nil {
		return m, err
	}

	udta := findChildWithType(moov, "udta")
	if udta.Size == 0 {
		return m, fmt.Errorf("failed to parse udta atom")
	}

	meta := findChildWithType(udta, "meta")
	if meta.Size == 0 {
		return m, fmt.Errorf("failed to parse meta atom")
	}

	ilst := findChildWithType(meta, "ilst")
	if ilst.Size == 0 {
		return m, fmt.Errorf("failed to parse ilst atom")
	}

	m = parse(ilst, path)

	return m, nil
}

// isM4AContainer reads the file at the expected location of the file type atom
// to confirm that the file uses the M4A container.
func isM4AContainer(file *os.File) bool {
	b := make([]byte, 7)
	file.Seek(4, 0)
	file.Read(b)

	return strings.EqualFold(string(b), "ftypM4A")
}

// readAtomAt reads the file at the specified offset and returns an Atom if one
// starts at that location.
func readAtomAt(f *os.File, off int64) (Atom, error) {
	var a Atom
	b := make([]byte, SizeLen)

	// The first four bytes of an atom are always the size field.
	_, err := f.ReadAt(b, off)
	if err != nil {
		return a, fmt.Errorf(ReadAtomAtError, f.Name(), "size", off, err)
	}
	size, err := parseSize(b)
	if err != nil {
		return a, fmt.Errorf(ReadAtomAtError, f.Name(), "size", off, err)
	}

	// The second four bytes of an atom are always the type field.
	_, err = f.ReadAt(b, off+SizeLen)
	if err != nil {
		return a, fmt.Errorf(ReadAtomAtError, f.Name(), "type", off, err)
	}
	typ := parseType(b)

	// An atom, including size and type fields, starts at offset and continues
	// for the number of contiguous bytes indicated by the size field.
	b = make([]byte, size)
	_, err = f.ReadAt(b, off)
	if err != nil {
		return a, fmt.Errorf(ReadAtomAtError, f.Name(), "content", off, err)
	}

	a = Atom{
		Start:    off,
		End:      off + int64(size),
		Size:     size,
		Type:     typ,
		Contents: b}
	return a, nil
}

// parseSize  returns the byte buffer containing atom size as an int.
func parseSize(buffer []byte) (int, error) {
	var boxSizeInBytes int
	encodedStr := hex.EncodeToString(buffer)
	s, err := strconv.ParseInt(encodedStr, 16, 0)
	if err != nil {
		return boxSizeInBytes, err
	}

	boxSizeInBytes = int(s)
	return boxSizeInBytes, nil
}

// parseType returns the byte buffer containing atom type as a string.
func parseType(buffer []byte) string {
	return string(buffer)
}

// readSiblingWithType reads the file for siblings of the given atom and
// returns the one with the target type.
func readSiblingWithType(f *os.File, a Atom, typ string) (Atom, error) {
	for {
		sib, err := readSibling(f, a)
		if err != nil {
			return sib, fmt.Errorf(ReadSiblingWithTypeError, f.Name(), a.Type, typ, err)
		}
		if sib.Type == typ {
			return sib, nil
		}
		a = sib
	}
}

// readSibling reads the file and returns the sibling immediately adjacent to
// the given atom.
func readSibling(f *os.File, a Atom) (Atom, error) {
	sib, err := readAtomAt(f, a.End)
	if err != nil {
		return sib, fmt.Errorf(ReadSiblingEror, f.Name(), a.Type, err)
	}
	return sib, nil
}

// findChildWithType searches the parent atom for a child of the given type and
// if successful returns the child atom.
func findChildWithType(parent Atom, typ string) Atom {
	for i := SizeLen + TypeLen; i < len(parent.Contents); {
		off := parent.Start + int64(i)
		child := parseAtom(parent.Contents[i:], off)
		if child.Size != 0 {
			if strings.Contains(child.Type, typ) {
				return child
			}
			i += child.Size
		} else {
			i += 4
		}
	}

	return Atom{}
}

// parseAtom parses the given byte buffer into an Atom and returns the struct.
func parseAtom(b []byte, off int64) Atom {
	var a Atom
	if len(b) >= SizeLen+TypeLen {
		size, _ := parseSize(b[0:4])
		if size <= len(b) {
			typ := parseType(b[4:8])
			contents := b[0:size]
			a = Atom{
				Start:    off,
				End:      off + int64(size),
				Size:     size,
				Type:     typ,
				Contents: contents}
		}
	}

	return a
}

// parse parses the metdata item list atom and returns metadata.
func parse(ilstAtom Atom, path string) metadata.Metadata {
	title := getTitleMetadata(ilstAtom)
	artist := getArtistMetadata(ilstAtom)
	artistSort := getArtistSortMetadata(ilstAtom)
	album := getAlbumMetadata(ilstAtom)
	albumSort := getAlbumSortMetadata(ilstAtom)
	year := getYearMetadata(ilstAtom)
	track := getTrackMetadata(ilstAtom)
	genre := getGenreMetadata(ilstAtom)
	comment := getCommentMetadata(ilstAtom)
	lyrics := getLyricsMetadata(ilstAtom)

	m := metadata.Metadata{
		Title:      title,
		Artist:     artist,
		ArtistSort: artistSort,
		Album:      album,
		AlbumSort:  albumSort,
		Year:       year,
		Track:      track,
		Genre:      genre,
		Comment:    comment,
		Lyrics:     lyrics,
		FileBase:   filepath.Base(path),
		FileDir:    filepath.Dir(path)}

	return m
}

func getTitleMetadata(ilstAtom Atom) string {
	titleAtom := findChildWithType(ilstAtom, "nam")
	if titleAtom.Size > 0 {
		contents := getMetadataValueFromDataAtom(titleAtom)
		title := string(contents)
		return title
	}

	return ""
}

func getArtistMetadata(ilstAtom Atom) string {
	artistAtom := findChildWithType(ilstAtom, "ART")
	if artistAtom.Size > 0 {
		contents := getMetadataValueFromDataAtom(artistAtom)
		artist := string(contents)
		return artist
	}

	return ""
}

func getArtistSortMetadata(ilstAtom Atom) string {
	artistSortAtom := findChildWithType(ilstAtom, "soar")
	if artistSortAtom.Size > 0 {
		contents := getMetadataValueFromDataAtom(artistSortAtom)
		artistSort := string(contents)
		return artistSort
	}

	return ""
}

func getAlbumMetadata(ilstAtom Atom) string {
	albumAtom := findChildWithType(ilstAtom, "alb")
	if albumAtom.Size > 0 {
		contents := getMetadataValueFromDataAtom(albumAtom)
		album := string(contents)
		return album
	}
	return ""
}

func getAlbumSortMetadata(ilstAtom Atom) string {
	albumSortAtom := findChildWithType(ilstAtom, "soal")
	if albumSortAtom.Size > 0 {
		contents := getMetadataValueFromDataAtom(albumSortAtom)
		albumSort := string(contents)
		return albumSort
	}

	return ""
}

func getYearMetadata(ilstAtom Atom) string {
	yearAtom := findChildWithType(ilstAtom, "day")
	if yearAtom.Size > 0 {
		contents := getMetadataValueFromDataAtom(yearAtom)
		year := string(contents)
		return year
	}

	return ""
}

func getTrackMetadata(ilstAtom Atom) string {
	trackAtom := findChildWithType(ilstAtom, "trkn")
	if trackAtom.Size > 0 {
		contents := getMetadataValueFromDataAtom(trackAtom)
		trackNumber := binary.BigEndian.Uint32(contents[0:4])
		track := strconv.FormatUint(uint64(trackNumber), 10)
		return track
	}

	return ""
}

func getGenreMetadata(ilstAtom Atom) string {
	genreAtom := findChildWithType(ilstAtom, "gnre")

	if genreAtom.Size < 1 {
		genreAtom = findChildWithType(ilstAtom, "gen")
	}

	if genreAtom.Size > 0 {
		contents := getMetadataValueFromDataAtom(genreAtom)
		genreField := binary.BigEndian.Uint16(contents)
		genreNumber := metadata.PredefinedGenre(genreField - 1)
		genre := genreNumber.PredefinedGenre()
		if genre == "" {
			genre = string(contents)
		}

		return genre
	}

	return ""
}

func getCommentMetadata(ilstAtom Atom) string {
	commentAtom := findChildWithType(ilstAtom, "cmt")
	if commentAtom.Size > 0 {
		contents := getMetadataValueFromDataAtom(commentAtom)
		comment := string(contents)
		return comment
	}

	return ""
}

func getLyricsMetadata(ilstAtom Atom) string {
	lyricsAtom := findChildWithType(ilstAtom, "lyr")

	if lyricsAtom.Size > 0 {
		contents := getMetadataValueFromDataAtom(lyricsAtom)
		lyrics := string(contents)
		return lyrics
	}

	return ""
}

// getMetadataValueFromDataAtom parses a data atom for relevant metadata value
// and returns the metadata value contents.
//
// The form of a data atom is as follows:
// - size: [0][1][2][3]
// - name: [4][5][6][7]
// - typeset: [8]
// - wellKnownType: [9][10][11]
// - countryIndicator: [12][13]
// - languageIndicator: [14][15]
func getMetadataValueFromDataAtom(dataAtom Atom) []byte {
	valueAtom := dataAtom.Contents[8:]
	metadataValue := valueAtom[16:]
	return metadataValue
}
