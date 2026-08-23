🎉 Event Organizer (Guest Manager) — Multi‑Language Event Guest Management
8 languages, one complete guest management system – create events, manage guest lists, track RSVPs, send reminders, and export guest data – right from your terminal.

✨ Features
👤 Add guests – name, email, RSVP status, plus‑one, dietary restrictions

📋 List guests – view all guests with RSVP status (filter by status)

📊 Event statistics – see total guests, attending, declined, pending

✅ RSVP tracking – update guest status (Pending, Attending, Declined)

📝 Add notes – attach custom notes to each guest

📤 Export to CSV – for printing or spreadsheet analysis

🔍 Search – find guests by name or email

💾 Persistent storage – all data saved in events.json

🧰 Supported Languages & Files
Language	File	Dependencies
Python	event_guests.py	none (stdlib)
Go	event_guests.go	none (stdlib)
JavaScript (Node)	event_guests.js	commander (optional)
Ruby	event_guests.rb	json, date
PHP	event_guests.php	none (extensions)
Java	EventGuests.java	Java 8+
C#	EventGuests.cs	.NET Core 3.1+
C++	event_guests.cpp	nlohmann/json
🚀 Quick Start
All implementations follow the same CLI pattern:

bash
# Create a new event
<command> create "Birthday Party" --date "2026-12-25" --venue "My House"

# Add a guest to an event (use event ID from list)
<command> add-guest <event_id> "Alice" "alice@email.com"

# Add guest with plus‑one and dietary restrictions
<command> add-guest <event_id> "Bob" "bob@email.com" --plus-one --dietary "Vegetarian"

# List all guests for an event
<command> list-guests <event_id>

# List guests by RSVP status
<command> list-guests <event_id> --status Attending

# Update RSVP status
<command> rsvp <event_id> <guest_id> Attending

# Show event statistics
<command> stats <event_id>

# Export guest list to CSV
<command> export <event_id> --filename guests.csv
Commands:

create <name> [--date DATE] [--venue VENUE] – create an event

add-guest <event_id> <name> <email> [--plus-one] [--dietary DIET] [--notes TEXT] – add guest

list-guests <event_id> [--status STATUS] – list guests (filter by status)

rsvp <event_id> <guest_id> <status> – update RSVP status

stats <event_id> – show event statistics

search <event_id> <term> – search guests

export <event_id> [--filename FILE] – export to CSV

📸 Example Output
text
📋 Event: Birthday Party
Date: 2026-12-25
Venue: My House

👤 Guests (5):
  1. Alice (alice@email.com) – Pending
  2. Bob (bob@email.com) – Attending (plus‑one: Yes) 🥗 Vegetarian
  3. Charlie (charlie@email.com) – Declined
  4. Diana (diana@email.com) – Pending (plus‑one: No)
  5. Eve (eve@email.com) – Attending

📊 Statistics:
  Total: 5
  Attending: 2
  Pending: 2
  Declined: 1
  Plus‑ones: 1
📁 Repository Structure
text
.
├── README.md
├── python/
│   └── event_guests.py
├── go/
│   └── event_guests.go
├── javascript/
│   └── event_guests.js
├── ruby/
│   └── event_guests.rb
├── php/
│   └── event_guests.php
├── java/
│   └── EventGuests.java
├── csharp/
│   └── EventGuests.cs
└── cpp/
    └── event_guests.cpp
