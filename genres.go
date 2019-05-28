package metadata

// PredefinedGenre represents a predefined ID3 genre.
type PredefinedGenre int

// PredefinedGenre returns genre as a string if its value maps to a predefined
// ID3 genre or otherwise returns an empty string.
func (pg PredefinedGenre) PredefinedGenre() string {
	switch pg {
	case Blues:
		return "Blues"
	case ClassicRock:
		return "Blues"
	case Country:
		return "Country"
	case Dance:
		return "Dance"
	case Disco:
		return "Disco"
	case Funk:
		return "Funk"
	case Grunge:
		return "Grunge"
	case HipHop:
		return "Hip-Hop"
	case Jazz:
		return "Jazz"
	case Metal:
		return "Metal"
	case NewAge:
		return "New Age"
	case Oldies:
		return "Oldies"
	case Other:
		return "Other"
	case Pop:
		return "Pop"
	case RAndB:
		return "R&B"
	case Rap:
		return "Rap"
	case Reggae:
		return "Reggae"
	case Rock:
		return "Rock"
	case Techno:
		return "Techno"
	case Industrial:
		return "Industrial"
	case Alternative:
		return "Alternative"
	case Ska:
		return "Ska"
	case DeathMetal:
		return "Death Metal"
	case Pranks:
		return "Pranks"
	case Soundtrack:
		return "Soundtrack"
	case EuroTechno:
		return "Euro-Techno"
	case Ambient:
		return "Ambient"
	case TripHop:
		return "Trip-Hop"
	case Vocal:
		return "Vocal"
	case JazzFunk:
		return "Jazz+Funk"
	case Fusion:
		return "Fusion"
	case Trance:
		return "Trance"
	case Classical:
		return "Classical"
	case Instrumental:
		return "Instrumental"
	case Acid:
		return "Acid"
	case House:
		return "House"
	case Game:
		return "Game"
	case SoundClip:
		return "Sound Clip"
	case Gospel:
		return "Gospel"
	case Noise:
		return "Noise"
	case AlternRock:
		return "Alternative Rock"
	case Bass:
		return "Bass"
	case Soul:
		return "Soul"
	case Punk:
		return "Punk"
	case Space:
		return "Space"
	case Meditative:
		return "Meditative"
	case InstrumentalPop:
		return "Instrumental Pop"
	case InstrumentalRock:
		return "Instrumental Rock"
	case Ethnic:
		return "Ethnic"
	case Gothic:
		return "Gothic"
	case Darkwave:
		return "Darkwave"
	case TechnoIndustrial:
		return "Techno-Industrial"
	case Electronic:
		return "Electronic"
	case PopFolk:
		return "Pop-Folk"
	case Eurodance:
		return "Eurodance"
	case Dream:
		return "Dream"
	case SouthernRock:
		return "Southern Rock"
	case Comedy:
		return "Comedy"
	case Cult:
		return "Cult"
	case Gangsta:
		return "Gangsta"
	case Top40:
		return "Top 40"
	case ChristianRap:
		return "Christian Rap"
	case PopFunk:
		return "Pop/Funk"
	case Jungle:
		return "Jungle"
	case NativeAmerican:
		return "Native US"
	case Cabaret:
		return "Cabaret"
	case NewWave:
		return "New Wave"
	case Psychadelic:
		return "Psychadelic"
	case Rave:
		return "Rave"
	case Showtunes:
		return "Showtunes"
	case Trailer:
		return "Trailer"
	case LoFi:
		return "Lo-Fi"
	case Tribal:
		return "Tribal"
	case AcidPunk:
		return "Acid Punk"
	case AcidJazz:
		return "Acid Jazz"
	case Polka:
		return "Polka"
	case Retro:
		return "Retro"
	case Musical:
		return "Musical"
	case RockAndRoll:
		return "Rock & Roll"
	case HardRock:
		return "Hard Rock"
	default:
		return ""
	}
}

// Constants represents the predefined ID3 genres.
const (
	Blues PredefinedGenre = iota
	ClassicRock
	Country
	Dance
	Disco
	Funk
	Grunge
	HipHop
	Jazz
	Metal
	NewAge
	Oldies
	Other
	Pop
	RAndB
	Rap
	Reggae
	Rock
	Techno
	Industrial
	Alternative
	Ska
	DeathMetal
	Pranks
	Soundtrack
	EuroTechno
	Ambient
	TripHop
	Vocal
	JazzFunk
	Fusion
	Trance
	Classical
	Instrumental
	Acid
	House
	Game
	SoundClip
	Gospel
	Noise
	AlternRock
	Bass
	Soul
	Punk
	Space
	Meditative
	InstrumentalPop
	InstrumentalRock
	Ethnic
	Gothic
	Darkwave
	TechnoIndustrial
	Electronic
	PopFolk
	Eurodance
	Dream
	SouthernRock
	Comedy
	Cult
	Gangsta
	Top40
	ChristianRap
	PopFunk
	Jungle
	NativeAmerican
	Cabaret
	NewWave
	Psychadelic
	Rave
	Showtunes
	Trailer
	LoFi
	Tribal
	AcidPunk
	AcidJazz
	Polka
	Retro
	Musical
	RockAndRoll
	HardRock
)
