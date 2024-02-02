package example

import (
	"fmt"
	"log"
	"os"

	"github.com/D7682/scribble"
	"github.com/brianvoe/gofakeit/v6"
)

// Person defines a simple person struct
type Person struct {
	Name string
	Age  int
}

// NewPerson constructor for creating a new person
func NewPerson(name string, age int) *Person {
	return &Person{
		Name: name,
		Age:  age,
	}
}

// PeopleExample is an example demonstrating people-related functionality
type PeopleExample struct {
	db *scribble.Driver
}

// NewPeopleExample creates a new PeopleExample instance
func NewPeopleExample(db *scribble.Driver) *PeopleExample {
	return &PeopleExample{
		db: db,
	}
}

// SetupDatabase sets up the database and returns a scribble driver instance
func (pe *PeopleExample) SetupDatabase() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	pe.db, err = scribble.New(path, nil)
	return err
}

// WritePeopleToDatabase writes a slice of Person instances to the database
func (pe *PeopleExample) WritePeopleToDatabase(people []*Person) {
	for _, person := range people {
		err := pe.db.Write("people", person.Name, person)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Wrote: %v\n", person.Name)
	}
}

// DeletePersonFromDatabase deletes a person from the database by name
func (pe *PeopleExample) DeletePersonFromDatabase(personName string) {
	err := pe.db.Delete("people", personName)
	if os.IsNotExist(err) {
		log.Fatal(err)
	}
}

// GenerateFakeUsers generates a specified number of fake user instances
func (pe *PeopleExample) GenerateFakeUsers(count int) []*Person {
	var fakeUsers []*Person

	for i := 0; i < count; i++ {
		name := gofakeit.FirstName()
		age := gofakeit.Number(0, 120)
		fakeUser := NewPerson(name, age)
		fakeUsers = append(fakeUsers, fakeUser)
	}

	return fakeUsers
}

// func main() {
// 	start := time.Now()

// 	// Create a new PeopleExample instance
// 	peopleExample := NewPeopleExample(nil) // Pass nil, as it will be initialized in SetupDatabase

// 	// Setup database
// 	err := peopleExample.SetupDatabase()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Generate fake users
// 	fakeUsers := peopleExample.GenerateFakeUsers(1000)

// 	// Write fake users to the database
// 	peopleExample.WritePeopleToDatabase(fakeUsers)

// 	// Delete a fake user from the database (just an example)
// 	peopleExample.DeletePersonFromDatabase(fakeUsers[0].Name)

// 	fmt.Println("Done.")
// 	fmt.Printf("Time Elapsed: %v\n", time.Since(start))
// }
