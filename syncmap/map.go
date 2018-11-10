package syncmap

import (
	"errors"
	"sync"

	"github.com/philippgille/gokv/util"
)

// Store is a gokv.Store implementation for a Go sync.Map.
type Store struct {
	m             *sync.Map
	marshalFormat MarshalFormat
}

// Set stores the given object for the given key.
// Values are marshalled to JSON automatically.
func (m Store) Set(k string, v interface{}) error {
	var data []byte
	var err error
	switch m.marshalFormat {
	case JSON:
		data, err = util.ToJSON(v)
	case Gob:
		data, err = util.ToGob(v)
	default:
		return errors.New("The store seems to be configured with a marshal format that's not implemented yet")
	}
	if err != nil {
		return err
	}

	m.m.Store(k, data)
	return nil
}

// Get retrieves the stored value for the given key.
// You need to pass a pointer to the value, so in case of a struct
// the automatic unmarshalling can populate the fields of the object
// that v points to with the values of the retrieved object's values.
func (m Store) Get(k string, v interface{}) (bool, error) {
	data, found := m.m.Load(k)
	if !found {
		return false, nil
	}

	switch m.marshalFormat {
	case JSON:
		return true, util.FromJSON(data.([]byte), v)
	case Gob:
		return true, util.FromGob(data.([]byte), v)
	default:
		return false, errors.New("The store seems to be configured with a marshal format that's not implemented yet")
	}
}

// Delete deletes the stored value for the given key.
// Deleting a non-existing key-value pair does NOT lead to an error.
func (m Store) Delete(k string) error {
	m.m.Delete(k)
	return nil
}

// MarshalFormat is an enum for the available (un-)marshal formats of this gokv.Store implementation.
type MarshalFormat int

const (
	// JSON is the MarshalFormat for (un-)marshalling to/from JSON
	JSON MarshalFormat = iota
	// Gob is the MarshalFormat for (un-)marshalling to/from gob
	Gob
)

// Options are the options for the Go sync.Map store.
type Options struct {
	// (Un-)marshal format.
	// Optional (JSON by default).
	MarshalFormat MarshalFormat
}

// DefaultOptions is an Options object with default values.
// MarshalFormat: JSON
var DefaultOptions = Options{
	// No need to set MarshalFormat to JSON
	// because its zero value is fine.
}

// NewStore creates a new Go sync.Map store.
func NewStore(options Options) Store {
	return Store{
		m:             &sync.Map{},
		marshalFormat: options.MarshalFormat,
	}
}
