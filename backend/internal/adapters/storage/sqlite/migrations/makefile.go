package main

// import (
// 	"fmt"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/gofrs/uuid/v5"
// )

// type User struct {
// 	Id          uuid.UUID
// 	Email       string
// 	Password    string
// 	FirstName   string
// 	LastName    string
// 	DateOfBirth time.Time
// 	Avatar      string
// 	Nickname    string
// 	About       string
// 	CreatedAt   time.Time
// 	IsPublic    bool
// }

// type Follower struct {
// 	Id         uuid.UUID
// 	FollowerId uuid.UUID
// 	FolloweeId uuid.UUID
// 	// Status     string
// 	Date time.Time
// }

// type Post struct {
// 	Id      uuid.UUID
// 	UserId  uuid.UUID
// 	Title   string
// 	Content string
// 	Image   string
// 	Privacy string
// 	date    time.Time
// }

// type Group struct {
// 	Id          uuid.UUID
// 	CreatorId   uuid.UUID
// 	Title       string
// 	Description string
// 	CreatedAt   time.Time
// }

// type Event struct {
// 	Id          uuid.UUID
// 	GroupId     uuid.UUID
// 	CreatorId   uuid.UUID
// 	Title       string
// 	Description string
// 	Option      string
// 	CreatedAt   time.Time
// }

// type Notification struct {
// 	Id        uuid.UUID
// 	UserId    uuid.UUID
// 	Type      string
// 	Message   string
// 	CreatedAt time.Time
// 	DeletedAt time.Time
// }

// // GenerateMigrationFiles generates blank migration files for each struct
// func GenerateMigrationFiles(structs []string) {
// 	for _, s := range structs {
// 		upFilename := fmt.Sprintf("000001_create_table_%s.up.sql", strings.ToLower(s))
// 		downFilename := fmt.Sprintf("000001_create_table_%s.down.sql", strings.ToLower(s))

// 		// Create up migration file
// 		upFile, err := os.Create(upFilename)
// 		if err != nil {
// 			fmt.Printf("Error creating %s: %s\n", upFilename, err)
// 			continue
// 		}
// 		defer upFile.Close()
// 		downFile, err := os.Create(downFilename)
// 		if err != nil {
// 			fmt.Printf("Error creating %s: %s\n", downFilename, err)
// 			continue
// 		}
// 		defer downFile.Close()

// 		fmt.Printf("Created migration files for struct %s\n", s)
// 	}
// }

// func main() {
// 	structs := []string{"User", "Follower", "Post", "Group", "Event", "Notification", "Chat", "Comment"}
// 	GenerateMigrationFiles(structs)
// }
