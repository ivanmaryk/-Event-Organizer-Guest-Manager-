// EventGuests.java
import java.io.*;
import java.nio.file.*;
import java.time.*;
import java.util.*;
import com.google.gson.*;

class Guest {
    String id, name, email, rsvp, dietary, notes, added;
    boolean plus_one;

    Guest() {}
    Guest(String name, String email, boolean plusOne, String dietary, String notes) {
        this.id = UUID.randomUUID().toString().substring(0,8);
        this.name = name;
        this.email = email;
        this.rsvp = "Pending";
        this.plus_one = plusOne;
        this.dietary = dietary;
        this.notes = notes;
        this.added = Instant.now().toString();
    }
}

class Event {
    String id, name, date, venue, created;
    List<Guest> guests = new ArrayList<>();

    Event() {}
    Event(String name, String date, String venue) {
        this.id = UUID.randomUUID().toString().substring(0,8);
        this.name = name;
        this.date = date;
        this.venue = venue;
        this.created = Instant.now().toString();
    }
}

public class EventGuests {
    private List<Event> events = new ArrayList<>();
    private final String dataFile = "events.json";
    private final Gson gson = new GsonBuilder().setPrettyPrinting().create();

    public EventGuests() { load(); }

    private void load() {
        try {
            Path path = Paths.get(dataFile);
            if (Files.exists(path)) {
                String json = new String(Files.readAllBytes(path));
                Event[] arr = gson.fromJson(json, Event[].class);
                events = Arrays.asList(arr);
            }
        } catch (Exception e) {}
    }

    private void save() {
        try {
            Files.write(Paths.get(dataFile), gson.toJson(events).getBytes());
        } catch (Exception e) {}
    }

    private Event getEvent(String id) {
        for (Event e : events) if (e.id.equals(id)) return e;
        return null;
    }

    private Guest getGuest(Event e, String id) {
        for (Guest g : e.guests) if (g.id.equals(id)) return g;
        return null;
    }

    public void create(String name, String date, String venue) {
        Event e = new Event(name, date, venue);
        events.add(e);
        save();
        System.out.printf("✅ Event created: %s (ID: %s)%n", e.name, e.id);
    }

    public void addGuest(String eventId, String name, String email, boolean plusOne, String dietary, String notes) {
        Event e = getEvent(eventId);
        if (e == null) {
            System.out.printf("Event %s not found.%n", eventId);
            return;
        }
        Guest g = new Guest(name, email, plusOne, dietary, notes);
        e.guests.add(g);
        save();
        System.out.printf("✅ Guest added: %s (ID: %s)%n", g.name, g.id);
    }

    public void listGuests(String eventId, String status) {
        Event e = getEvent(eventId);
        if (e == null) {
            System.out.printf("Event %s not found.%n", eventId);
            return;
        }
        List<Guest> guests = e.guests;
        if (status != null && !status.isEmpty()) {
            guests = new ArrayList<>();
            for (Guest g : e.guests) {
                if (g.rsvp.equals(status)) guests.add(g);
            }
        }
        if (guests.isEmpty()) {
            System.out.println("No guests.");
            return;
        }
        System.out.printf("\n📋 Event: %s%n", e.name);
        if (e.date != null && !e.date.isEmpty()) System.out.printf("Date: %s%n", e.date);
        if (e.venue != null && !e.venue.isEmpty()) System.out.printf("Venue: %s%n", e.venue);
        System.out.printf("\n👤 Guests (%d):%n", guests.size());
        for (int i = 0; i < guests.size(); i++) {
            Guest g = guests.get(i);
            String plus = g.plus_one ? " (plus‑one: Yes)" : " (plus‑one: No)";
            String dietary = g.dietary != null && !g.dietary.isEmpty() ? " 🥗 " + g.dietary : "";
            String notes = g.notes != null && !g.notes.isEmpty() ? " 📝 " + g.notes : "";
            System.out.printf("  %d. %s (%s) – %s%s%s%s%n", i+1, g.name, g.email, g.rsvp, plus, dietary, notes);
        }
    }

    public void rsvp(String eventId, String guestId, String status) {
        Event e = getEvent(eventId);
        if (e == null) {
            System.out.printf("Event %s not found.%n", eventId);
            return;
        }
        Guest g = getGuest(e, guestId);
        if (g == null) {
            System.out.printf("Guest %s not found.%n", guestId);
            return;
        }
        List<String> validStatuses = Arrays.asList("Pending", "Attending", "Declined");
        if (!validStatuses.contains(status)) {
            System.out.println("Invalid status. Choose: Pending, Attending, Declined");
            return;
        }
        g.rsvp = status;
        save();
        System.out.printf("✅ %s RSVP updated to: %s%n", g.name, status);
    }

    public void stats(String eventId) {
        Event e = getEvent(eventId);
        if (e == null) {
            System.out.printf("Event %s not found.%n", eventId);
            return;
        }
        int total = e.guests.size();
        int attending = 0, pending = 0, declined = 0, plusOnes = 0;
        for (Guest g : e.guests) {
            switch (g.rsvp) {
                case "Attending": attending++; break;
                case "Pending": pending++; break;
                case "Declined": declined++; break;
            }
            if (g.plus_one) plusOnes++;
        }
        System.out.printf("\n📊 Event: %s%n", e.name);
        System.out.printf("  Total guests: %d%n", total);
        System.out.printf("  Attending: %d%n", attending);
        System.out.printf("  Pending: %d%n", pending);
        System.out.printf("  Declined: %d%n", declined);
        System.out.printf("  Plus‑ones: %d%n", plusOnes);
    }

    public void search(String eventId, String term) {
        Event e = getEvent(eventId);
        if (e == null) {
            System.out.printf("Event %s not found.%n", eventId);
            return;
        }
        String lower = term.toLowerCase();
        List<Guest> results = new ArrayList<>();
        for (Guest g : e.guests) {
            if (g.name.toLowerCase().contains(lower) || g.email.toLowerCase().contains(lower)) {
                results.add(g);
            }
        }
        if (results.isEmpty()) {
            System.out.println("No matches.");
            return;
        }
        System.out.printf("\n🔍 Found %d guest(s):%n", results.size());
        for (int i = 0; i < results.size(); i++) {
            Guest g = results.get(i);
            System.out.printf("  %d. %s (%s) – %s%n", i+1, g.name, g.email, g.rsvp);
        }
    }

    public void export(String eventId, String filename) throws IOException {
        Event e = getEvent(eventId);
        if (e == null) {
            System.out.printf("Event %s not found.%n", eventId);
            return;
        }
        try (BufferedWriter writer = Files.newBufferedWriter(Paths.get(filename))) {
            writer.write("ID,Name,Email,RSVP,Plus-One,Dietary,Notes,Added\n");
            for (Guest g : e.guests) {
                writer.write(String.format("%s,%s,%s,%s,%b,%s,%s,%s%n",
                    g.id, g.name, g.email, g.rsvp, g.plus_one, g.dietary, g.notes, g.added));
            }
        }
        System.out.printf("✅ Exported %d guests to %s%n", e.guests.size(), filename);
    }

    public static void main(String[] args) throws Exception {
        if (args.length < 1) {
            System.out.println("Usage: EventGuests <command> [options]");
            return;
        }
        EventGuests app = new EventGuests();
        String cmd = args[0];
        Map<String, String> params = new HashMap<>();
        for (int i=1; i<args.length; i++) {
            if (args[i].startsWith("--") && i+1 < args.length) {
                params.put(args[i].substring(2), args[++i]);
            }
        }
        switch (cmd) {
            case "create":
                if (args.length < 2) { System.out.println("create <name> [--date DATE] [--venue VENUE]"); return; }
                app.create(args[1], params.getOrDefault("date", ""), params.getOrDefault("venue", ""));
                break;
            case "add-guest":
                if (args.length < 4) { System.out.println("add-guest <event_id> <name> <email> [--plus-one] [--dietary DIET] [--notes NOTES]"); return; }
                app.addGuest(args[1], args[2], args[3],
                    params.containsKey("plus-one"),
                    params.getOrDefault("dietary", ""),
                    params.getOrDefault("notes", ""));
                break;
            case "list-guests":
                if (args.length < 2) { System.out.println("list-guests <event_id> [--status STATUS]"); return; }
                app.listGuests(args[1], params.get("status"));
                break;
            case "rsvp":
                if (args.length < 4) { System.out.println("rsvp <event_id> <guest_id> <status>"); return; }
                app.rsvp(args[1], args[2], args[3]);
                break;
            case "stats":
                if (args.length < 2) { System.out.println("stats <event_id>"); return; }
                app.stats(args[1]);
                break;
            case "search":
                if (args.length < 3) { System.out.println("search <event_id> <term>"); return; }
                app.search(args[1], args[2]);
                break;
            case "export":
                if (args.length < 2) { System.out.println("export <event_id> [--filename FILE]"); return; }
                app.export(args[1], params.getOrDefault("filename", "guests.csv"));
                break;
            default:
                System.out.println("Unknown command.");
        }
    }
}
