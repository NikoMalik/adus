package bytespool

import (
	"bytes"
	"testing"
)

func TestAllocateAndFree(t *testing.T) {
	// Initialize the pool
	Init()

	// Test Allocate function with different sizes
	sizes := []int{
		1 << 11, // 2048 bytes
		1 << 12, // 4096 bytes
		1 << 13, // 8192 bytes
	}

	for _, size := range sizes {
		// Allocate a buffer
		buf := Allocate(size)
		if len(buf) != size {
			t.Errorf("Expected buffer size %d, got %d", size, len(buf))
		}

		// Ensure the buffer is zeroed out (if necessary)
		if !bytes.Equal(buf, make([]byte, size)) {
			t.Errorf("Buffer not zeroed after allocation for size %d", size)
		}

		// Free the buffer
		Free(buf)

		// Allocate again and check if it's reused
		buf2 := Allocate(size)
		if len(buf2) != size {
			t.Errorf("Expected buffer size %d after reuse, got %d", size, len(buf2))
		}

		// Free the second buffer
		Free(buf2)
	}
}

func TestAllocateLargeSize(t *testing.T) {
	// Test that allocating a buffer larger than the largest pool size
	// correctly allocates using `make` instead of the pool.
	largeSize := 1 << 14 // 16384 bytes, larger than any pool size
	buf := Allocate(largeSize)
	if len(buf) != largeSize {
		t.Errorf("Expected buffer size %d, got %d", largeSize, len(buf))
	}

	// Free should not panic even though the buffer is larger than any pool
	Free(buf)
}

func TestBytePoolUsage(t *testing.T) {

	Init()

	size := 1 << 12 // 4096 bytes // 4 kb
	buf := Allocate(size)
	if len(buf) != size {
		t.Fatalf("Expected buffer of size %d, got %d", size, len(buf))
	}

	// write data to buffer
	for i := range buf {
		buf[i] = byte(i % 256)
	}

	// return buffer to pull
	Free(buf)

	// again get buffer from pool
	buf2 := Allocate(size)
	if len(buf2) != size {
		t.Fatalf("Expected buffer of size %d, got %d", size, len(buf2))
	}

	// check if data was written correctly
	for i := range buf2 {
		if buf2[i] != byte(i%256) {
			t.Errorf("Buffer contents don't match at index %d: expected %d, got %d", i, byte(i%256), buf2[i])
		}
	}

	// check another buffer size
	buf3 := Allocate(100)
	if len(buf3) != 100 {
		t.Fatalf("Expected buffer of size 100, got %d", len(buf3))
	}
}

func TestPoolUsage(t *testing.T) {
	size := 100
	buf := Allocate(size)
	if len(buf) < size {
		t.Errorf("Expected buffer of at least size %d, got %d", size, len(buf))
	}
}

func TestAllocate(t *testing.T) {
	Init()

	tests := [5]struct {
		name     string
		size     int
		expected int
	}{
		{"small size", 100, 2048},
		{"medium size", 2048, 2048},
		{"large size", 8192, 8192},
		{"extra large size", 16384, 16384},
		{"too large", 20000, 20000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := Allocate(tt.size)
			if len(buf) != tt.size {
				t.Errorf("expected buffer of size %d, got %d", tt.size, len(buf))
			}
			if cap(buf) < tt.expected {
				t.Errorf("expected buffer capacity to be at least %d, got %d", tt.expected, cap(buf))
			}
		})
	}
}

func TestFree(t *testing.T) {
	Init()

	size := 2048
	buf := Allocate(size)
	if len(buf) != size {
		t.Errorf("expected buffer of size %d, got %d", size, len(buf))
	}

	// Check the buffer's capacity before freeing it
	if cap(buf) != 2048 {
		t.Errorf("expected buffer capacity to be %d, got %d", 2048, cap(buf))
	}

	Free(buf)

	// Allocate again and ensure the buffer was returned to the pool
	newBuf := Allocate(size)
	if cap(newBuf) < 2048 {
		t.Errorf("expected buffer capacity to be at least %d, got %d", 2048, cap(newBuf))
	}
}

func TestAllocateAndFreeCycle(t *testing.T) {
	Init()

	sizes := [3]int{100, 2048, 8192}
	for _, size := range sizes {
		buf := Allocate(size)
		if len(buf) != size {
			t.Errorf("expected buffer of size %d, got %d", size, len(buf))
		}

		Free(buf)

		// Re-allocate and ensure buffer is reused
		newBuf := Allocate(size)
		if cap(newBuf) < size {
			t.Errorf("expected buffer capacity to be at least %d, got %d", size, cap(newBuf))
		}
	}
}
