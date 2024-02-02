// Package scribble provides a simple disk-based JSON storage system.
package scribble

import (
	"encoding/json"
	"github.com/D7682/scribble/pkg/errors"
	"github.com/jcelliott/lumber"
	"os"
	"path/filepath"
	"sync"
)

// Logger defines the interface for logging methods.
type Logger interface {
	Fatal(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
}

// Driver represents the main struct for interacting with the scribble database.
type Driver struct {
	mutex         sync.RWMutex
	resourceLocks sync.Map
	dir           string
	log           Logger
}

// Options represents the optional configurations for the scribble driver.
type Options struct {
	Logger
}

// New creates a new scribble database driver instance.
// It initializes a new database if it does not already exist.
func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	opts := Options{}

	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		dir:           dir,
		resourceLocks: sync.Map{},
		log:           opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating scribble database at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

// Write writes the given data to a resource within a collection in the scribble database.
func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return errors.ErrMissingCollection
	}

	if resource == "" {
		return errors.ErrResourceNotFound
	}

	mutex := d.getOrCreateLock(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	return write(dir, tmpPath, fnlPath, v)
}

// write is a helper function for writing data to a file.
func write(dir, tmpPath, dstPath string, v interface{}) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.NewFileIOError(dir, err)
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := os.WriteFile(tmpPath, b, 0644); err != nil {
		return errors.NewFileIOError(dir, err)
	}

	return os.Rename(tmpPath, dstPath)
}

// Read reads data from a resource within a collection in the scribble database.
func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" {
		return errors.ErrMissingCollection
	}

	if resource == "" {
		return errors.ErrResourceNotFound
	}

	record := filepath.Join(d.dir, collection, resource)
	return read(record, v)
}

// read is a helper function for reading data from a file.
func read(record string, v interface{}) error {
	b, err := os.ReadFile(record + ".json")
	if err != nil {
		return errors.NewFileIOError(record+".json", err)
	}

	return json.Unmarshal(b, v)
}

// ReadAll retrieves all records from a collection in the scribble database.
func (d *Driver) ReadAll(collection string) ([][]byte, error) {
	if collection == "" {
		return nil, errors.ErrMissingCollection
	}

	dir := filepath.Join(d.dir, collection)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.NewFileIOError(dir, err)
	}

	return readAll(files, dir)
}

// readAll is a helper function for reading all records from a collection.
func readAll(files []os.DirEntry, dir string) ([][]byte, error) {
	var records [][]byte

	for _, file := range files {
		b, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, errors.NewFileIOError(filepath.Join(dir, file.Name()), err)
		}
		records = append(records, b)
	}

	return records, nil
}

// Delete removes a resource within a collection from the scribble database.
func (d *Driver) Delete(collection, resource string) error {
	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateLock(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)
	fi, err := stat(dir)

	if fi == nil || err != nil {
		return errors.NewNotFoundError(path, os.ErrNotExist)
	}

	switch {
	case fi.Mode().IsDir():
		if err := os.RemoveAll(dir); err != nil {
			return errors.NewFileIOError(dir, err)
		}
	case fi.Mode().IsRegular():
		if err := os.RemoveAll(dir + ".json"); err != nil {
			return errors.NewFileIOError(dir+".json", err)
		}
	}

	return nil
}

// stat is a helper function for obtaining file information.
func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

// getOrCreateLock retrieves or creates a lock for a collection to ensure thread safety.
func (d *Driver) getOrCreateLock(collection string) *sync.Mutex {
	// Load or store a new lock for the collection
	l, loaded := d.resourceLocks.LoadOrStore(collection, &sync.Mutex{})

	// If the lock was not loaded, it means a new one was stored, and we need to unlock
	// the main mutex to avoid potential contention.
	if !loaded {
		d.mutex.Unlock()
	}

	return l.(*sync.Mutex)
}
