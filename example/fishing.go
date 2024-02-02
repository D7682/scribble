package example

import (
	"encoding/json"

	"github.com/D7682/scribble"
)

// Fish represents a fish in the example
type Fish struct {
	Name string
}

// FishingExample is an example demonstrating fishing functionality
type FishingExample struct {
	db *scribble.Driver
}

// NewFishingExample creates a new FishingExample instance
func NewFishingExample(db *scribble.Driver) *FishingExample {
	return &FishingExample{db: db}
}

// WriteFishToDatabase writes a fish to the database
func (fe *FishingExample) WriteFishToDatabase(name string) error {
	return fe.db.Write("fish", name, Fish{Name: name})
}

// ReadFishFromDatabase reads a fish from the database
func (fe *FishingExample) ReadFishFromDatabase(name string) (*Fish, error) {
	fish := &Fish{}
	err := fe.db.Read("fish", name, fish)
	return fish, err
}

// ReadAllFishFromDatabase reads all fish from the database
func (fe *FishingExample) ReadAllFishFromDatabase() ([]Fish, error) {
	records, err := fe.db.ReadAll("fish")
	if err != nil {
		return nil, err
	}

	var fishies []Fish
	for _, f := range records {
		fishFound := Fish{}
		if err := json.Unmarshal(f, &fishFound); err != nil {
			return nil, err
		}
		fishies = append(fishies, fishFound)
	}

	return fishies, nil
}

// DeleteFishFromDatabase deletes a fish from the database
func (fe *FishingExample) DeleteFishFromDatabase(name string) error {
	return fe.db.Delete("fish", name)
}

// DeleteAllFishFromDatabase deletes all fish from the database
func (fe *FishingExample) DeleteAllFishFromDatabase() error {
	return fe.db.Delete("fish", "")
}

// func main() {
// 	// Example usage
// 	dir := "./"

// 	db, err := scribble.New(dir, nil)
// 	if err != nil {
// 		fmt.Println("Error", err)
// 	}

// 	fishingExample := NewFishingExample(db)

// 	// Write a fish to the database
// 	for _, name := range []string{"onefish", "twofish", "redfish", "bluefish"} {
// 		err := fishingExample.WriteFishToDatabase(name)
// 		if err != nil {
// 			fmt.Println("Error writing fish:", err)
// 		}
// 	}

// 	// Read a fish from the database
// 	onefish, err := fishingExample.ReadFishFromDatabase("onefish")
// 	if err != nil {
// 		fmt.Println("Error reading fish:", err)
// 	} else {
// 		fmt.Println("Read fish:", onefish)
// 	}

// 	// Read all fish from the database
// 	allFish, err := fishingExample.ReadAllFishFromDatabase()
// 	if err != nil {
// 		fmt.Println("Error reading all fish:", err)
// 	} else {
// 		fmt.Println("Read all fish:", allFish)
// 	}

// 	// Delete a fish from the database
// 	err = fishingExample.DeleteFishFromDatabase("onefish")
// 	if err != nil {
// 		fmt.Println("Error deleting fish:", err)
// 	}

// 	// Delete all fish from the database
// 	err = fishingExample.DeleteAllFishFromDatabase()
// 	if err != nil {
// 		fmt.Println("Error deleting all fish:", err)
// 	}
// }
