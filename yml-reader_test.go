package ymlreader

import (
	"context"
	"io"
	"os"
	"testing"
)

func TestUnmarshalFile(t *testing.T) {
	tt := []struct {
		Name      string
		Path      string
		MustError bool
	}{
		{Name: "Valid path to file", Path: "./examples/simple.yml", MustError: false},
		{Name: "Invalid path to file", Path: "./examples/no_exists.yml", MustError: true},
		{Name: "Broken file", Path: "./examples/broken.yml", MustError: true},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := UnmarshalFile(tc.Path)
			if tc.MustError && err == nil {
				t.Errorf("failed to unmarshal: %s", err)
			}
		})
	}
}

func TestNewFromFile(t *testing.T) {
	tt := []struct {
		Name      string
		Path      string
		MustError bool
	}{
		{Name: "Valid path to file", Path: "./examples/simple.yml", MustError: false},
		{Name: "Invalid path to file", Path: "./examples/no_exists.yml", MustError: true},
		// New has not run unmarshal, so broken file is NOT necessary.
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := NewFromFile(tc.Path)
			if tc.MustError && err == nil {
				t.Errorf("failed to unmarshal: %s", err)
			}
		})
	}
}

func TestYMLReader_Close(t *testing.T) {
	t.Run("YMLReader is not init", func(t *testing.T) {
		err := (*YMLReader)(nil).Close()
		if err == nil {
			t.Errorf("must be failed: %s", "yml reader: is not initialized")
		}
	})

	t.Run("YMLReader is init", func(t *testing.T) {
		y, _ := NewFromFile("./examples/simple.yml")
		err := y.Close()
		if err != nil {
			t.Errorf("failed to clode: %s", err)
		}
	})
}

func TestNewReader(t *testing.T) {
	validReader, _ := os.Open("./examples/simple.yml")
	defer validReader.Close()
	tt := []struct {
		Name      string
		reader    io.ReadCloser
		MustError bool
	}{
		{Name: "Valid reader", reader: validReader, MustError: false},
		{Name: "Invalid path to file", reader: nil, MustError: true},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := NewFromReader(tc.reader)
			if tc.MustError && err == nil {
				t.Errorf("failed to unmarshal: %s", err)
			}
		})
	}
}

func TestYMLReader_StartRead(t *testing.T) {
	validReader, _ := os.Open("./examples/simple.yml")
	defer validReader.Close()
	y, _ := NewFromReader(validReader)
	err := y.StartRead(context.TODO())
	if err != nil {
		t.Errorf("failed to read: %s", err)
	}
	if y.Len() != 1 {
		t.Errorf("length of offers must be 1")
	}
	if y.ReadCount() != 1 {
		t.Errorf("read count must be 1")
	}
}
