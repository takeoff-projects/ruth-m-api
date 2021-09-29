package main
 
import (
    "encoding/json"
    "fmt"
    "log"
	"net/http"
	"io/ioutil"
	"os"
	"context"
	
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"cloud.google.com/go/datastore"
)

// Event model
type Event struct {
    Name     string `json:"id"`
    Title  string `json:"title" datastore: "title"`
    Location string `json:"location" datastore: "location"`
    When   string `json:"when" datastore: "when"`
}
 
// Events array
var Events []Event
 
func home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the Events API!")
    fmt.Println("Endpoint Hit: homeP")
}
 
func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", home)
	myRouter.HandleFunc("/events", getEvents).Methods("GET")
	myRouter.HandleFunc("/events/{id}",getEventbyID).Methods("GET")
	myRouter.HandleFunc("/events", createEvent).Methods("POST")
	myRouter.HandleFunc("/events/{id}", updateEvent).Methods("PUT")
	myRouter.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on Port: %s\n", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), myRouter))
}
 
func getEvents(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: getEvents")

    ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    dsClient, err := datastore.NewClient(ctx, projectID)
    if err != nil  {
        http.Error(w, err.Error(), 500)
        return
    }

    var events[]Event
    _, err = dsClient.GetAll(ctx, datastore.NewQuery("Event"), &events)
    if err != nil  {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(events)
}

func getEventbyID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getEventbyID")
    vars := mux.Vars(r)
	key := vars["id"]
	
	fmt.Printf("Key: %s\n", key)

	ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    dsClient, err := datastore.NewClient(ctx, projectID)
    if err != nil  {
        http.Error(w, err.Error(), 500)
        return
    }

    event := Event{}
    nameKey := datastore.NameKey("Event", key, nil)
    err = dsClient.Get(ctx, nameKey, &event)
    if err != nil  {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(event)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createEvent")
	newID := uuid.New().String()
	fmt.Println(newID)
	
    reqBody, _ := ioutil.ReadAll(r.Body)
    var event Event 
	json.Unmarshal(reqBody, &event)
	event.Name = newID

    ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    dsClient, err := datastore.NewClient(ctx, projectID)
    if err != nil  {
        http.Error(w, err.Error(), 500)
        return
    }

	key := datastore.NameKey("Event", newID, nil)
	event.Name = newID
	_, err = dsClient.Put(ctx, key, &event)
    if err != nil  {
        http.Error(w, err.Error(), 500)
        return
    }
 
    json.NewEncoder(w).Encode(event)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var updatedEvent Event
	json.Unmarshal(reqBody, &updatedEvent)
	vars := mux.Vars(r)
	id := vars["id"]

    ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    dsClient, err := datastore.NewClient(ctx, projectID)
    if err != nil  {
        http.Error(w, err.Error(), 500)
        return
    }

    key := datastore.NameKey("Event", id, nil)
    _, err = dsClient.Put(ctx, key, &updatedEvent)
    if err != nil  {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(updatedEvent)
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    fmt.Println("Endpoint hit: Delete, " + id)
 
    ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    dsClient, err := datastore.NewClient(ctx, projectID)
    if err != nil  {
        http.Error(w, err.Error(), 500)
        return
    }

    key := datastore.NameKey("Event", id, nil)
    err = dsClient.Delete(ctx, key)
    if err != nil  {
        http.Error(w, err.Error(), 500)
        return
    }
}

func main() {
    ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    dsClient, err := datastore.NewClient(ctx, projectID)
    if err != nil  {
        log.Println(err.Error())
        return
    }

    events := []*Event{
        {Name: "2944a9cb-ef2d-4632-ac1d-af2b2629d0f2",
         Title: "Dinner",
         Location: "My House",
         When: "Tonight"},
         {Name: "f88f1860-9a5d-423e-820f-9acb4db3030e",
          Title: "Go Programming Lesson",
          Location: "At School",
          When: "Tomorrow"},
         {Name: "4cb393fb-dd19-469e-a52c-22a12c0a98df",
          Title: "Company Picnic",
          Location: "At the Park",
          When: "Saturday"},
    }

    keys := []*datastore.Key{
        datastore.NameKey("Event", events[0].Name, nil),
        datastore.NameKey("Event", events[1].Name, nil),
        datastore.NameKey("Event", events[2].Name, nil),
    }

    _, err = dsClient.PutMulti(ctx, keys, events)
    if err != nil  {
        log.Println(err.Error())
    }
	
    handleRequests()
}