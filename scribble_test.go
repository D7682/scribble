package scribble

import (
	"os"
	"testing"
)

// Fish represents a fish with a type.
type Fish struct {
	Type string `json:"type"`
}

var (
	db         *Driver
	database   = "./deep/school"
	collection = "fish"
	onefish    = Fish{}
	twofish    = Fish{}
	redfish    = Fish{Type: "red"}
	bluefish   = Fish{Type: "blue"}
)

// TestMain sets up and runs tests, ensuring cleanup afterward.
func TestMain(m *testing.M) {
	setupAndRunTests(m)
}

// setupAndRunTests removes the test directory before and after running tests.
func setupAndRunTests(m *testing.M) {
	removeTestDir()
	defer removeTestDir()

	code := m.Run()

	os.Exit(code)
}

// removeTestDir removes the test directory.
func removeTestDir() {
	err := os.RemoveAll("./deep")
	if err != nil {
		return
	}
}

// createDB creates a new Scribble database.
func createDB() error {
	var err error
	if db, err = New(database, nil); err != nil {
		return err
	}
	return nil
}

// TestNew tests the creation of a new database.
func TestNew(t *testing.T) {
	assertDatabaseNotExists(t)

	err := createDB()
	if err != nil {
		return
	}

	assertDatabaseExists(t)

	err = createDB()
	if err != nil {
		return
	}

	assertDatabaseExists(t)
}

// assertDatabaseNotExists asserts that the database does not exist.
func assertDatabaseNotExists(t *testing.T) {
	if _, err := os.Stat(database); err == nil {
		t.Error("Expected nothing, got database")
	}
}

// assertDatabaseExists asserts that the database exists.
func assertDatabaseExists(t *testing.T) {
	if _, err := os.Stat(database); err != nil {
		t.Error("Expected database, got nothing")
	}
}

// TestWriteAndRead tests writing and reading fish from the database.
func TestWriteAndRead(t *testing.T) {
	err := createDB()
	if err != nil {
		return
	}

	writeFishToDatabase(t, "redfish", redfish)

	readFishFromDatabase(t, "redfish")

	assertFishType(t, "red")
	err = destroySchool()
	if err != nil {
		return
	}
}

// writeFishToDatabase writes a fish to the database.
func writeFishToDatabase(t *testing.T, key string, fish Fish) {
	if err := db.Write(collection, key, fish); err != nil {
		t.Error("Create fish failed: ", err.Error())
	}
}

// readFishFromDatabase reads a fish from the database.
func readFishFromDatabase(t *testing.T, key string) {
	if err := db.Read(collection, key, &onefish); err != nil {
		t.Error("Failed to read: ", err.Error())
	}
}

// assertFishType asserts that the type of the fish is as expected.
func assertFishType(t *testing.T, expectedType string) {
	if onefish.Type != expectedType {
		t.Error("Expected red fish, got: ", onefish.Type)
	}
}

// assertFishCountNonZero asserts that the count of fish is non-zero.
func assertFishCountNonZero(t *testing.T) {
	fish, err := db.ReadAll(collection)
	if err != nil {
		t.Error("Failed to read: ", err.Error())
	}

	if len(fish) == 0 {
		t.Error("Expected non-zero fish count, got zero")
	}
}

// TestReadall tests reading all fish from the database.
func TestReadall(t *testing.T) {
	err := createDB()
	if err != nil {
		return
	}
	err = createSchool()
	if err != nil {
		return
	}

	readAllFishFromDatabase(t)

	assertFishCountNonZero(t)
	err = destroySchool()
	if err != nil {
		return
	}
}

// readAllFishFromDatabase reads all fish from the database.
func readAllFishFromDatabase(t *testing.T) {
	fish, err := db.ReadAll(collection)
	if err != nil {
		t.Error("Failed to read: ", err.Error())
	}

	if len(fish) <= 0 {
		t.Error("Expected some fish, have none")
	}
}

// TestWriteAndReadEmpty tests writing and reading empty fish to/from the database.
func TestWriteAndReadEmpty(t *testing.T) {
	err := createDB()
	if err != nil {
		return
	}

	writeEmptyFishToDatabase(t, "redfish", redfish)
	writeEmptyFishToDatabase(t, "", redfish)

	readEmptyFishFromDatabase(t, "redfish")

	err = destroySchool()
	if err != nil {
		return
	}
}

// writeEmptyFishToDatabase attempts to write empty fish to the database, expecting an error.
func writeEmptyFishToDatabase(t *testing.T, key string, fish Fish) {
	if err := db.Write(key, "redfish", fish); err == nil {
		t.Error("Allowed write of empty resource", err.Error())
	}
}

// readEmptyFishFromDatabase attempts to read empty fish from the database, expecting an error.
func readEmptyFishFromDatabase(t *testing.T, key string) {
	if err := db.Read("", key, onefish); err == nil {
		t.Error("Allowed read of empty resource", err.Error())
	}
}

// TestDelete tests deleting a fish from the database.
func TestDelete(t *testing.T) {
	err := createDB()
	if err != nil {
		return
	}
	writeFishToDatabase(t, "redfish", redfish)

	deleteFishFromDatabase(t, "redfish")

	readDeletedFishFromDatabase(t, "redfish")

	err = destroySchool()
	if err != nil {
		return
	}
}

// deleteFishFromDatabase deletes a fish from the database.
func deleteFishFromDatabase(t *testing.T, key string) {
	if err := db.Delete(collection, key); err != nil {
		t.Error("Failed to delete: ", err.Error())
	}
}

// readDeletedFishFromDatabase attempts to read a deleted fish from the database, expecting an error.
func readDeletedFishFromDatabase(t *testing.T, key string) {
	if err := db.Read(collection, key, &onefish); err == nil {
		t.Error("Expected nothing, got fish")
	}
}

// TestDeleteall tests deleting all fish from the database.
func TestDeleteall(t *testing.T) {
	err := createDB()
	if err != nil {
		return
	}
	err = createSchool()
	if err != nil {
		return
	}

	deleteAllFishFromDatabase(t)

	assertCollectionNotExists(t)
	err = destroySchool()
	if err != nil {
		return
	}
}

// deleteAllFishFromDatabase deletes all fish from the database.
func deleteAllFishFromDatabase(t *testing.T) {
	if err := db.Delete(collection, ""); err != nil {
		t.Error("Failed to delete: ", err.Error())
	}
}

// assertCollectionNotExists asserts that the collection does not exist.
func assertCollectionNotExists(t *testing.T) {
	if _, err := os.Stat(collection); err == nil {
		t.Error("Expected nothing, have fish")
	}
}

// createFish creates a fish in the database.
func createFish(fish Fish) error {
	return db.Write(collection, fish.Type, fish)
}

// createSchool creates multiple fish in the database.
func createSchool() error {
	for _, f := range []Fish{{Type: "red"}, {Type: "blue"}} {
		if err := db.Write(collection, f.Type, f); err != nil {
			return err
		}
	}

	return nil
}

// destroyFish deletes a specific fish from the database.
func destroyFish(name string) error {
	return db.Delete(collection, name)
}

// destroySchool deletes all fish from the database.
func destroySchool() error {
	return db.Delete(collection, "")
}
