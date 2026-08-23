# event_guests.py
import json
import os
import sys
import argparse
import uuid
from datetime import datetime

DATA_FILE = "events.json"

class Guest:
    def __init__(self, name, email, plus_one=False, dietary="", notes="", guest_id=None):
        self.id = guest_id or str(uuid.uuid4())[:8]
        self.name = name
        self.email = email
        self.rsvp = "Pending"  # Pending, Attending, Declined
        self.plus_one = plus_one
        self.dietary = dietary
        self.notes = notes
        self.added = datetime.now().isoformat()

    def to_dict(self):
        return {
            "id": self.id,
            "name": self.name,
            "email": self.email,
            "rsvp": self.rsvp,
            "plus_one": self.plus_one,
            "dietary": self.dietary,
            "notes": self.notes,
            "added": self.added
        }

    @classmethod
    def from_dict(cls, data):
        g = cls(data["name"], data["email"], data.get("plus_one", False),
                data.get("dietary", ""), data.get("notes", ""), data.get("id"))
        g.rsvp = data.get("rsvp", "Pending")
        g.added = data.get("added", datetime.now().isoformat())
        return g

class Event:
    def __init__(self, name, date="", venue="", event_id=None):
        self.id = event_id or str(uuid.uuid4())[:8]
        self.name = name
        self.date = date
        self.venue = venue
        self.guests = []
        self.created = datetime.now().isoformat()

    def to_dict(self):
        return {
            "id": self.id,
            "name": self.name,
            "date": self.date,
            "venue": self.venue,
            "guests": [g.to_dict() for g in self.guests],
            "created": self.created
        }

    @classmethod
    def from_dict(cls, data):
        e = cls(data["name"], data.get("date", ""), data.get("venue", ""), data.get("id"))
        e.guests = [Guest.from_dict(g) for g in data.get("guests", [])]
        e.created = data.get("created", datetime.now().isoformat())
        return e

class Organizer:
    def __init__(self):
        self.events = []
        self.load()

    def load(self):
        if os.path.exists(DATA_FILE):
            with open(DATA_FILE, "r") as f:
                data = json.load(f)
                self.events = [Event.from_dict(e) for e in data]

    def save(self):
        with open(DATA_FILE, "w") as f:
            json.dump([e.to_dict() for e in self.events], f, indent=2)

    def get_event(self, event_id):
        for e in self.events:
            if e.id == event_id:
                return e
        return None

    def get_guest(self, event, guest_id):
        for g in event.guests:
            if g.id == guest_id:
                return g
        return None

    def create(self, name, date="", venue=""):
        e = Event(name, date, venue)
        self.events.append(e)
        self.save()
        print(f"✅ Event created: {e.name} (ID: {e.id})")

    def add_guest(self, event_id, name, email, plus_one=False, dietary="", notes=""):
        e = self.get_event(event_id)
        if not e:
            print(f"Event {event_id} not found.")
            return
        g = Guest(name, email, plus_one, dietary, notes)
        e.guests.append(g)
        self.save()
        print(f"✅ Guest added: {g.name} (ID: {g.id})")

    def list_guests(self, event_id, status=None):
        e = self.get_event(event_id)
        if not e:
            print(f"Event {event_id} not found.")
            return
        guests = e.guests
        if status:
            guests = [g for g in guests if g.rsvp.lower() == status.lower()]
        if not guests:
            print("No guests.")
            return
        print(f"\n📋 Event: {e.name}")
        if e.date:
            print(f"Date: {e.date}")
        if e.venue:
            print(f"Venue: {e.venue}")
        print(f"\n👤 Guests ({len(guests)}):")
        for i, g in enumerate(guests, 1):
            plus = " (plus‑one: Yes)" if g.plus_one else " (plus‑one: No)"
            dietary = f" 🥗 {g.dietary}" if g.dietary else ""
            notes = f" 📝 {g.notes}" if g.notes else ""
            print(f"  {i}. {g.name} ({g.email}) – {g.rsvp}{plus}{dietary}{notes}")

    def rsvp(self, event_id, guest_id, status):
        e = self.get_event(event_id)
        if not e:
            print(f"Event {event_id} not found.")
            return
        g = self.get_guest(e, guest_id)
        if not g:
            print(f"Guest {guest_id} not found.")
            return
        valid_statuses = ["Pending", "Attending", "Declined"]
        if status not in valid_statuses:
            print(f"Invalid status. Choose: {', '.join(valid_statuses)}")
            return
        g.rsvp = status
        self.save()
        print(f"✅ {g.name} RSVP updated to: {status}")

    def stats(self, event_id):
        e = self.get_event(event_id)
        if not e:
            print(f"Event {event_id} not found.")
            return
        total = len(e.guests)
        attending = sum(1 for g in e.guests if g.rsvp == "Attending")
        pending = sum(1 for g in e.guests if g.rsvp == "Pending")
        declined = sum(1 for g in e.guests if g.rsvp == "Declined")
        plus_ones = sum(1 for g in e.guests if g.plus_one)
        print(f"\n📊 Event: {e.name}")
        print(f"  Total guests: {total}")
        print(f"  Attending: {attending}")
        print(f"  Pending: {pending}")
        print(f"  Declined: {declined}")
        print(f"  Plus‑ones: {plus_ones}")

    def search(self, event_id, term):
        e = self.get_event(event_id)
        if not e:
            print(f"Event {event_id} not found.")
            return
        term_lower = term.lower()
        results = [g for g in e.guests if term_lower in g.name.lower() or term_lower in g.email.lower()]
        if not results:
            print("No matches.")
            return
        print(f"\n🔍 Found {len(results)} guest(s):")
        for i, g in enumerate(results, 1):
            print(f"  {i}. {g.name} ({g.email}) – {g.rsvp}")

    def export(self, event_id, filename):
        import csv
        e = self.get_event(event_id)
        if not e:
            print(f"Event {event_id} not found.")
            return
        with open(filename, 'w', newline='') as f:
            writer = csv.writer(f)
            writer.writerow(["ID", "Name", "Email", "RSVP", "Plus-One", "Dietary", "Notes", "Added"])
            for g in e.guests:
                writer.writerow([g.id, g.name, g.email, g.rsvp, g.plus_one, g.dietary, g.notes, g.added])
        print(f"✅ Exported {len(e.guests)} guests to {filename}")

def main():
    parser = argparse.ArgumentParser(description="Event Organizer - Guest Manager")
    subparsers = parser.add_subparsers(dest="cmd", required=True)

    create_parser = subparsers.add_parser("create")
    create_parser.add_argument("name")
    create_parser.add_argument("--date", default="")
    create_parser.add_argument("--venue", default="")

    add_parser = subparsers.add_parser("add-guest")
    add_parser.add_argument("event_id")
    add_parser.add_argument("name")
    add_parser.add_argument("email")
    add_parser.add_argument("--plus-one", action="store_true")
    add_parser.add_argument("--dietary", default="")
    add_parser.add_argument("--notes", default="")

    list_parser = subparsers.add_parser("list-guests")
    list_parser.add_argument("event_id")
    list_parser.add_argument("--status", choices=["Pending", "Attending", "Declined"])

    rsvp_parser = subparsers.add_parser("rsvp")
    rsvp_parser.add_argument("event_id")
    rsvp_parser.add_argument("guest_id")
    rsvp_parser.add_argument("status", choices=["Pending", "Attending", "Declined"])

    stats_parser = subparsers.add_parser("stats")
    stats_parser.add_argument("event_id")

    search_parser = subparsers.add_parser("search")
    search_parser.add_argument("event_id")
    search_parser.add_argument("term")

    export_parser = subparsers.add_parser("export")
    export_parser.add_argument("event_id")
    export_parser.add_argument("--filename", default="guests.csv")

    args = parser.parse_args()
    app = Organizer()

    if args.cmd == "create":
        app.create(args.name, args.date, args.venue)
    elif args.cmd == "add-guest":
        app.add_guest(args.event_id, args.name, args.email, args.plus_one, args.dietary, args.notes)
    elif args.cmd == "list-guests":
        app.list_guests(args.event_id, args.status)
    elif args.cmd == "rsvp":
        app.rsvp(args.event_id, args.guest_id, args.status)
    elif args.cmd == "stats":
        app.stats(args.event_id)
    elif args.cmd == "search":
        app.search(args.event_id, args.term)
    elif args.cmd == "export":
        app.export(args.event_id, args.filename)

if __name__ == "__main__":
    main()
