// event_guests.go
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
	"github.com/google/uuid"
)

type Guest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	RSVP     string `json:"rsvp"`
	PlusOne  bool   `json:"plus_one"`
	Dietary  string `json:"dietary"`
	Notes    string `json:"notes"`
	Added    string `json:"added"`
}

type Event struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Date    string  `json:"date"`
	Venue   string  `json:"venue"`
	Guests  []Guest `json:"guests"`
	Created string  `json:"created"`
}

type Organizer struct {
	Events []Event `json:"events"`
}

var dataFile = "events.json"

func (o *Organizer) load() {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return
	}
	json.Unmarshal(data, o)
}

func (o *Organizer) save() {
	data, _ := json.MarshalIndent(o, "", "  ")
	os.WriteFile(dataFile, data, 0644)
}

func (o *Organizer) getEvent(id string) *Event {
	for i := range o.Events {
		if o.Events[i].ID == id {
			return &o.Events[i]
		}
	}
	return nil
}

func (o *Organizer) getGuest(e *Event, id string) *Guest {
	for i := range e.Guests {
		if e.Guests[i].ID == id {
			return &e.Guests[i]
		}
	}
	return nil
}

func (o *Organizer) create(name, date, venue string) {
	e := Event{
		ID:      uuid.New().String()[:8],
		Name:    name,
		Date:    date,
		Venue:   venue,
		Created: time.Now().Format(time.RFC3339),
	}
	o.Events = append(o.Events, e)
	o.save()
	fmt.Printf("✅ Event created: %s (ID: %s)\n", e.Name, e.ID)
}

func (o *Organizer) addGuest(eventID, name, email string, plusOne bool, dietary, notes string) {
	e := o.getEvent(eventID)
	if e == nil {
		fmt.Printf("Event %s not found.\n", eventID)
		return
	}
	g := Guest{
		ID:      uuid.New().String()[:8],
		Name:    name,
		Email:   email,
		RSVP:    "Pending",
		PlusOne: plusOne,
		Dietary: dietary,
		Notes:   notes,
		Added:   time.Now().Format(time.RFC3339),
	}
	e.Guests = append(e.Guests, g)
	o.save()
	fmt.Printf("✅ Guest added: %s (ID: %s)\n", g.Name, g.ID)
}

func (o *Organizer) listGuests(eventID, status string) {
	e := o.getEvent(eventID)
	if e == nil {
		fmt.Printf("Event %s not found.\n", eventID)
		return
	}
	guests := e.Guests
	if status != "" {
		var filtered []Guest
		for _, g := range guests {
			if g.RSVP == status {
				filtered = append(filtered, g)
			}
		}
		guests = filtered
	}
	if len(guests) == 0 {
		fmt.Println("No guests.")
		return
	}
	fmt.Printf("\n📋 Event: %s\n", e.Name)
	if e.Date != "" {
		fmt.Printf("Date: %s\n", e.Date)
	}
	if e.Venue != "" {
		fmt.Printf("Venue: %s\n", e.Venue)
	}
	fmt.Printf("\n👤 Guests (%d):\n", len(guests))
	for i, g := range guests {
		plus := " (plus‑one: Yes)"
		if !g.PlusOne {
			plus = " (plus‑one: No)"
		}
		dietary := ""
		if g.Dietary != "" {
			dietary = fmt.Sprintf(" 🥗 %s", g.Dietary)
		}
		notes := ""
		if g.Notes != "" {
			notes = fmt.Sprintf(" 📝 %s", g.Notes)
		}
		fmt.Printf("  %d. %s (%s) – %s%s%s%s\n", i+1, g.Name, g.Email, g.RSVP, plus, dietary, notes)
	}
}

func (o *Organizer) rsvp(eventID, guestID, status string) {
	e := o.getEvent(eventID)
	if e == nil {
		fmt.Printf("Event %s not found.\n", eventID)
		return
	}
	g := o.getGuest(e, guestID)
	if g == nil {
		fmt.Printf("Guest %s not found.\n", guestID)
		return
	}
	validStatuses := map[string]bool{"Pending": true, "Attending": true, "Declined": true}
	if !validStatuses[status] {
		fmt.Println("Invalid status. Choose: Pending, Attending, Declined")
		return
	}
	g.RSVP = status
	o.save()
	fmt.Printf("✅ %s RSVP updated to: %s\n", g.Name, status)
}

func (o *Organizer) stats(eventID string) {
	e := o.getEvent(eventID)
	if e == nil {
		fmt.Printf("Event %s not found.\n", eventID)
		return
	}
	total := len(e.Guests)
	attending, pending, declined, plusOnes := 0, 0, 0, 0
	for _, g := range e.Guests {
		switch g.RSVP {
		case "Attending":
			attending++
		case "Pending":
			pending++
		case "Declined":
			declined++
		}
		if g.PlusOne {
			plusOnes++
		}
	}
	fmt.Printf("\n📊 Event: %s\n", e.Name)
	fmt.Printf("  Total guests: %d\n", total)
	fmt.Printf("  Attending: %d\n", attending)
	fmt.Printf("  Pending: %d\n", pending)
	fmt.Printf("  Declined: %d\n", declined)
	fmt.Printf("  Plus‑ones: %d\n", plusOnes)
}

func (o *Organizer) search(eventID, term string) {
	e := o.getEvent(eventID)
	if e == nil {
		fmt.Printf("Event %s not found.\n", eventID)
		return
	}
	var results []Guest
	for _, g := range e.Guests {
		if contains(g.Name, term) || contains(g.Email, term) {
			results = append(results, g)
		}
	}
	if len(results) == 0 {
		fmt.Println("No matches.")
		return
	}
	fmt.Printf("\n🔍 Found %d guest(s):\n", len(results))
	for i, g := range results {
		fmt.Printf("  %d. %s (%s) – %s\n", i+1, g.Name, g.Email, g.RSVP)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(s[0:len(substr)] == substr) || (len(s) > len(substr) && (s[len(s)-len(substr):] == substr || 
		indexOfSubstr(s, substr) != -1)))
}

func indexOfSubstr(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func (o *Organizer) export(eventID, filename string) {
	e := o.getEvent(eventID)
	if e == nil {
		fmt.Printf("Event %s not found.\n", eventID)
		return
	}
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write([]string{"ID", "Name", "Email", "RSVP", "Plus-One", "Dietary", "Notes", "Added"})
	for _, g := range e.Guests {
		w.Write([]string{g.ID, g.Name, g.Email, g.RSVP, fmt.Sprintf("%v", g.PlusOne), g.Dietary, g.Notes, g.Added})
	}
	fmt.Printf("✅ Exported %d guests to %s\n", len(e.Guests), filename)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: event_guests <command> [options]")
		return
	}
	o := &Organizer{}
	o.load()
	cmd := os.Args[1]

	switch cmd {
	case "create":
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		name := createCmd.String("name", "", "")
		date := createCmd.String("date", "", "")
		venue := createCmd.String("venue", "", "")
		createCmd.Parse(os.Args[2:])
		if *name == "" && len(createCmd.Args()) > 0 {
			*name = createCmd.Args()[0]
		}
		if *name == "" {
			fmt.Println("create requires a name")
			return
		}
		o.create(*name, *date, *venue)

	case "add-guest":
		addCmd := flag.NewFlagSet("add-guest", flag.ExitOnError)
		eventID := addCmd.String("event-id", "", "")
		name := addCmd.String("name", "", "")
		email := addCmd.String("email", "", "")
		plusOne := addCmd.Bool("plus-one", false, "")
		dietary := addCmd.String("dietary", "", "")
		notes := addCmd.String("notes", "", "")
		addCmd.Parse(os.Args[2:])
		args := addCmd.Args()
		if *eventID == "" && len(args) > 0 {
			*eventID = args[0]
		}
		if *name == "" && len(args) > 1 {
			*name = args[1]
		}
		if *email == "" && len(args) > 2 {
			*email = args[2]
		}
		if *eventID == "" || *name == "" || *email == "" {
			fmt.Println("add-guest requires event-id, name, email")
			return
		}
		o.addGuest(*eventID, *name, *email, *plusOne, *dietary, *notes)

	case "list-guests":
		listCmd := flag.NewFlagSet("list-guests", flag.ExitOnError)
		eventID := listCmd.String("event-id", "", "")
		status := listCmd.String("status", "", "")
		listCmd.Parse(os.Args[2:])
		if *eventID == "" && len(listCmd.Args()) > 0 {
			*eventID = listCmd.Args()[0]
		}
		if *eventID == "" {
			fmt.Println("list-guests requires event-id")
			return
		}
		o.listGuests(*eventID, *status)

	case "rsvp":
		rsvpCmd := flag.NewFlagSet("rsvp", flag.ExitOnError)
		eventID := rsvpCmd.String("event-id", "", "")
		guestID := rsvpCmd.String("guest-id", "", "")
		status := rsvpCmd.String("status", "", "")
		rsvpCmd.Parse(os.Args[2:])
		if *eventID == "" && len(rsvpCmd.Args()) > 0 {
			*eventID = rsvpCmd.Args()[0]
		}
		if *guestID == "" && len(rsvpCmd.Args()) > 1 {
			*guestID = rsvpCmd.Args()[1]
		}
		if *status == "" && len(rsvpCmd.Args()) > 2 {
			*status = rsvpCmd.Args()[2]
		}
		if *eventID == "" || *guestID == "" || *status == "" {
			fmt.Println("rsvp requires event-id, guest-id, status")
			return
		}
		o.rsvp(*eventID, *guestID, *status)

	case "stats":
		statsCmd := flag.NewFlagSet("stats", flag.ExitOnError)
		eventID := statsCmd.String("event-id", "", "")
		statsCmd.Parse(os.Args[2:])
		if *eventID == "" && len(statsCmd.Args()) > 0 {
			*eventID = statsCmd.Args()[0]
		}
		if *eventID == "" {
			fmt.Println("stats requires event-id")
			return
		}
		o.stats(*eventID)

	case "search":
		searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
		eventID := searchCmd.String("event-id", "", "")
		term := searchCmd.String("term", "", "")
		searchCmd.Parse(os.Args[2:])
		if *eventID == "" && len(searchCmd.Args()) > 0 {
			*eventID = searchCmd.Args()[0]
		}
		if *term == "" && len(searchCmd.Args()) > 1 {
			*term = searchCmd.Args()[1]
		}
		if *eventID == "" || *term == "" {
			fmt.Println("search requires event-id and term")
			return
		}
		o.search(*eventID, *term)

	case "export":
		exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
		eventID := exportCmd.String("event-id", "", "")
		filename := exportCmd.String("filename", "guests.csv", "")
		exportCmd.Parse(os.Args[2:])
		if *eventID == "" && len(exportCmd.Args()) > 0 {
			*eventID = exportCmd.Args()[0]
		}
		if *eventID == "" {
			fmt.Println("export requires event-id")
			return
		}
		o.export(*eventID, *filename)

	default:
		fmt.Println("Unknown command. Use create, add-guest, list-guests, rsvp, stats, search, export.")
	}
}
