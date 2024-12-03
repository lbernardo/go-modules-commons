package room

import (
	"strings"
	"testing"
)

func TestGenerateAuthCode(t *testing.T) {
	length := 5
	code, err := GenerateCode(length)
	if err != nil {
		t.Fatalf("failed to get euthentication: %v", err)
	}

	expectedLength := length
	if len(code) != expectedLength {
		t.Errorf("incorrect length: expected %d, got %d", expectedLength, len(code))
	}

	for _, char := range code {
		if !(strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", char)) {
			t.Errorf("invalid: '%c'", char)
		}
	}
}

func TestGenerateRoomID(t *testing.T) {
	roomID, err := GenerateRoomID()
	if err != nil {
		t.Fatalf("failed to generate roomId: %v", err)
	}

	sections := strings.Split(roomID, "-")
	if len(sections) != 3 {
		t.Errorf("invalid format: %s", roomID)
	}

	if len(sections[0]) != 3 || len(sections[1]) != 4 || len(sections[2]) != 3 {
		t.Errorf("length incorrect: %s", roomID)
	}

	for _, section := range sections {
		for _, char := range section {
			if !(strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", char)) {
				t.Errorf("invalid char: '%c'", char)
			}
		}
	}
}
