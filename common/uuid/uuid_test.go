package uuid

import (
	"fmt"
	"testing"
)

func TestNewUUID(t *testing.T) {
	uuid := New()

	if uuid == nil {
		t.Error("uuid is nil")
	}

	if len(uuid.Bytes()) != 16 {
		t.Error("uuid length is not 16")
	}

	t.Log(uuid)

}

func TestUUIDString(t *testing.T) {
	uuid := New()
	str := uuid.String()
	if str == "" {
		t.Error("uuid string is empty")
	}
	if len(str) != 36 {
		t.Error("uuid string length is not 36")
	}

	if str[8] != '-' || str[13] != '-' || str[18] != '-' || str[23] != '-' {
		t.Error("uuid string format is not correct")
	}

	t.Log(str)

}

func TestUUIDOperations(t *testing.T) {

	t.Run("TestNewUUID", func(t *testing.T) {
		uuid := New()
		if len(uuid.Bytes()) != 16 {
			t.Errorf("Expected UUID length to be 16 bytes, got %d", len(uuid.Bytes()))
		}
		fmt.Println("New UUID:", uuid.String())
	})

	t.Run("TestParseBytes", func(t *testing.T) {
		originalUUID := New()
		bytes := originalUUID.Bytes()
		uuid, err := ParseBytes(bytes)
		if err != nil {
			t.Fatalf("Failed to parse UUID from bytes: %v", err)
		}
		if !uuid.Equals(originalUUID) {
			t.Errorf("Parsed UUID does not match original UUID")
		}
		fmt.Println("Parsed UUID from bytes:", uuid.String())
	})

	t.Run("TestParseString", func(t *testing.T) {
		originalUUID := New()
		str := originalUUID.String()
		uuid, err := ParseString(str)
		if err != nil {
			t.Fatalf("Failed to parse UUID from string: %v", err)
		}
		if !uuid.Equals(originalUUID) {
			t.Errorf("Parsed UUID does not match original UUID")
		}
		fmt.Println("Parsed UUID from string:", uuid.String())
	})

	t.Run("TestParseStringInvalid", func(t *testing.T) {
		_, err := ParseString("invalid-uuid-string")
		if err == nil {
			t.Errorf("Expected error parsing invalid UUID string, got nil")
		}
		fmt.Println("Parsing invalid UUID string resulted in error:", err)
	})
}

func TestUUIDEquals(t *testing.T) {
	uuid1 := New()
	uuid2 := New()
	uuid3 := *uuid1

	if !uuid1.Equals(&uuid3) {
		t.Errorf("Expected UUIDs to be equal, but they are not.")
	}

	if uuid1.Equals(uuid2) {
		t.Errorf("Expected UUIDs to be different, but they are equal.")
	}

	var nilUUID *UUID
	if uuid1.Equals(nilUUID) {
		t.Errorf("Expected UUID and nil to be different, but they are equal.")
	}

	fmt.Printf("UUID1: %s\n", uuid1.String())
	fmt.Printf("UUID2: %s\n", uuid2.String())
}

func TestUUIDVersionAndVariant(t *testing.T) {
	uuid := New()
	str := uuid.String()

	if str[14] != '4' {
		t.Errorf("Expected UUID version to be 4, but got version %c", str[14])
	}

	if str[19] != '8' && str[19] != '9' && str[19] != 'a' && str[19] != 'b' {
		t.Errorf("Expected UUID variant to be 8, 9, a, or b, but got variant %c", str[19])
	}

	fmt.Printf("UUID: %s\n", str)
}

func TestUUIDStringConversion(t *testing.T) {
	originalUUID := New()
	str := originalUUID.String()
	uuid, err := ParseString(str)
	if err != nil {
		t.Fatalf("Failed to parse UUID from string: %v", err)
	}
	if !uuid.Equals(originalUUID) {
		t.Errorf("Parsed UUID does not match original UUID")
	}
	fmt.Printf("Original UUID: %s\n", originalUUID.String())
	fmt.Printf("Parsed UUID: %s\n", uuid.String())
}

func TestParseBytes(t *testing.T) {

	bytes := [16]byte{0x24, 0x18, 0xd0, 0x87, 0x64, 0x8d, 0x49, 0x90, 0x86, 0xe8, 0x19, 0xdc, 0xa1, 0xd0, 0x06, 0xd3}

	uuid, err := ParseBytes(bytes[:])

	if err != nil {
		t.Fatalf("Failed to parse UUID from bytes: %v", err)
	}

	if !uuid.Equals(&UUID{0x24, 0x18, 0xd0, 0x87, 0x64, 0x8d, 0x49, 0x90, 0x86, 0xe8, 0x19, 0xdc, 0xa1, 0xd0, 0x06, 0xd3}) {
		t.Errorf("Parsed UUID does not match original UUID")
	}

	_, err = ParseBytes([]byte{1, 3, 2, 4})
	if err == nil {
		t.Fatal("Expect error but nil")
	}
}

func TestRandom(t *testing.T) {
	uuid := New()
	uuid2 := New()

	if uuid.String() == uuid2.String() {
		t.Error("duplicated uuid")
	}
}

func TestEquals(t *testing.T) {
	var uuid *UUID
	var uuid2 *UUID

	if !uuid.Equals(uuid2) {
		t.Error("empty uuid should equal")
	}

	uuid3 := New()
	if uuid.Equals(uuid3) {
		t.Error("nil uuid equals non-nil uuid")
	}
}
