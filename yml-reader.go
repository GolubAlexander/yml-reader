package ymlreader

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

const offersBufferSize = 10

// UnmarshalFile unmarshals an yandex market language file.
// It's ok for small file, but memory expensive for the big one.
// Loads file into the memory.
func UnmarshalFile(path string) (*YMLCatalog, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var catalog YMLCatalog
	if err := xml.Unmarshal(b, &catalog); err != nil {
		return nil, err
	}
	return &catalog, err
}

// YMLReader represents yandex market language reader.
// Contains reader and buffered channel of offers and done channel.
// Reader implements ReadCloser interface.
type YMLReader struct {
	reader    io.ReadCloser
	offers    chan Offer
	done      chan struct{}
	mu        sync.Mutex
	readCount int
}

// NewFromFile initializes a reader from a file. Param is the path to file.
// If file can be open then YMLReader inits reader and prepares channel.
func NewFromFile(filePath string) (*YMLReader, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &YMLReader{
		reader: f,
		offers: make(chan Offer, offersBufferSize),
		done:   make(chan struct{}, 1),
	}, nil
}

// Close closes a reader. Can return an error if reader is not initialized.
func (y *YMLReader) Close() error {
	if y == nil {
		return fmt.Errorf("yml reader: is not initialized")
	}
	return y.reader.Close()
}

// NewFromReader creates an yandex market language reader from any ReadCloser interface.
// Can return error if passed reader is nil.
func NewFromReader(r io.ReadCloser) (*YMLReader, error) {
	if r == nil {
		return nil, fmt.Errorf("reader is nil")
	}
	return &YMLReader{
		reader: r,
		offers: make(chan Offer, offersBufferSize),
		done:   make(chan struct{}, 1),
	}, nil
}

func (y *YMLReader) StartRead(_ context.Context) error {
	d := xml.NewDecoder(y.reader)
	var elName string
	var offer Offer
	for {
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				y.done <- struct{}{}
				return nil
			}
			return err
		}
		switch tag := t.(type) {
		case xml.StartElement:
			elName = tag.Name.Local
			if elName == "offer" {
				if err := d.DecodeElement(&offer, &tag); err != nil {
					return err
				}
				y.mu.Lock()
				y.readCount++
				y.mu.Unlock()
				y.offers <- offer
			}
		default:
			continue
		}
	}
}

// Len returns current queue of offers up to offersBufferSize.
func (y *YMLReader) Len() int {
	return len(y.offers)
}

// ReadCount returns amount of read offers.
func (y *YMLReader) ReadCount() int {
	y.mu.Lock()
	defer y.mu.Unlock()
	return y.readCount
}

func (y *YMLReader) Offer() Offer {
	offer := <-y.offers
	return offer
}
