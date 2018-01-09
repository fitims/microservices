package filetype

import "testing"

type testMedia struct {
	media     []byte
	mediaType FileType
}

var media = []testMedia{
	{[]byte{0x49, 0x44, 0x33}, MP3},
	{[]byte{0xFF, 0xfb, 0x40}, MP3},
	{[]byte{0x52, 0x49, 0x46, 0x46, 0xFF, 0xFF, 0xFF, 0xFF, 0x57, 0x41, 0x56, 0x45}, WAV},
	{[]byte{0x46, 0x4F, 0x52, 0x4D, 0xFF, 0xFF, 0xFF, 0xFF, 0x41, 0x49, 0x46, 0x46}, AIFF},
	{[]byte{0x42, 0x44, 0x5f, 0x4f, 0x4d, 0x45, 0x43, 0x42}, Unknown},
	{[]byte{0x42}, Unknown},
}

func TestIsMp3(t *testing.T) {
	for _, item := range media {
		actual := IsMp3(item.media)
		expected := item.mediaType == MP3
		if actual != expected {
			t.Error("For", item.media, " - ", item.mediaType.String(), "expected", expected, "got", actual)
		}
	}
}

func TestIsWAV(t *testing.T) {
	for _, item := range media {
		actual := IsWav(item.media)
		expected := item.mediaType == WAV
		if actual != expected {
			t.Error("For", item.media, " - ", item.mediaType.String(), "expected", expected, "got", actual)
		}
	}
}

func TestIsAIFF(t *testing.T) {
	for _, item := range media {
		actual := IsAiff(item.media)
		expected := item.mediaType == AIFF
		if actual != expected {
			t.Error("For", item.mediaType.String(), "expected", expected, "got", actual)
		}
	}
}

func TestIsAudio(t *testing.T) {
	for _, item := range media {
		actual := IsAudio(item.media)
		expected := item.mediaType != Unknown
		if actual != expected {
			t.Error("For", item.mediaType.String(), "expected", expected, "got", actual)
		}
	}
}
