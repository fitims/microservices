package filetype

// valid file types
const (
	Unknown = iota
	MP3     = iota
	AIFF    = iota
	WAV     = iota
)

// FileType represent the valid file types
type FileType int

// Matcher is a signature for matching file types
type Matcher func([]byte) bool

// Matchers containst all valid matchers
var Matchers = make(map[FileType]Matcher)

// IsMp3 checks if the buffer is MP3 file
func IsMp3(buf []byte) bool {
	return len(buf) > 2 && ((buf[0] == 0x49 && buf[1] == 0x44 && buf[2] == 0x33) || (buf[0] == 0xFF && buf[1] == 0xfb))
}

// IsWav checks if the buffer is WAV file
func IsWav(buf []byte) bool {
	return len(buf) > 11 &&
		buf[0] == 0x52 && buf[1] == 0x49 &&
		buf[2] == 0x46 && buf[3] == 0x46 &&
		buf[8] == 0x57 && buf[9] == 0x41 &&
		buf[10] == 0x56 && buf[11] == 0x45
}

// IsAiff checks if the buffer is AIFF file
func IsAiff(buf []byte) bool {
	return len(buf) > 11 &&
		buf[0] == 0x46 && buf[1] == 0x4F &&
		buf[2] == 0x52 && buf[3] == 0x4D &&
		buf[8] == 0x41 && buf[9] == 0x49 &&
		buf[10] == 0x46 && buf[11] == 0x46
}

// IsAudio checks if the file is a valid audio file
func IsAudio(buf []byte) bool {
	for _, m := range Matchers {
		if m(buf) == true {
			return true
		}
	}
	return false
}

// Match tries to match the buffer with the valid file types and returns the valid type
func Match(buf []byte) FileType {
	for k, m := range Matchers {
		if m(buf) {
			return k
		}
	}
	return Unknown
}

func (f FileType) String() string {
	switch f {
	case MP3:
		return "MP3"
	case WAV:
		return "WAV"
	case AIFF:
		return "AIFF"
	default:
		return "Unknown"
	}
}

func init() {
	Matchers[MP3] = IsMp3
	Matchers[AIFF] = IsAiff
	Matchers[WAV] = IsWav
}
