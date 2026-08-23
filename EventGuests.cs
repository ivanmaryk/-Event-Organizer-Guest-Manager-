// EventGuests.cs
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.Json;
using System.Text.Json.Serialization;

class Guest
{
    [JsonPropertyName("id")] public string Id { get; set; }
    [JsonPropertyName("name")] public string Name { get; set; }
    [JsonPropertyName("email")] public string Email { get; set; }
    [JsonPropertyName("rsvp")] public string RSVP { get; set; }
    [JsonPropertyName("plus_one")] public bool PlusOne { get; set; }
    [JsonPropertyName("dietary")] public string Dietary { get; set; }
    [JsonPropertyName("notes")] public string Notes { get; set; }
    [JsonPropertyName("added")] public string Added { get; set; }

    public Guest() { }
    public Guest(string name, string email, bool plusOne, string dietary, string notes)
    {
        Id = Guid.NewGuid().ToString().Substring(0,8);
        Name = name;
        Email = email;
        RSVP = "Pending";
        PlusOne = plusOne;
        Dietary = dietary;
        Notes = notes;
        Added = DateTime.Now.ToString("o");
    }
}

class Event
{
    [JsonPropertyName("id")] public string Id { get; set; }
    [JsonPropertyName("name")] public string Name { get; set; }
    [JsonPropertyName("date")] public string Date { get; set; }
    [JsonPropertyName("venue")] public string Venue { get; set; }
    [JsonPropertyName("guests")] public List<Guest> Guests { get; set; } = new List<Guest>();
    [JsonPropertyName("created")] public string Created { get; set; }

    public Event() { }
    public Event(string name, string date, string venue)
    {
        Id = Guid.NewGuid().ToString().Substring(0,8);
        Name = name;
        Date = date;
        Venue = venue;
        Created = DateTime.Now.ToString("o");
    }
}

class Organizer
{
    private List<Event> events = new List<Event>();
    private readonly string dataFile = "events.json";
    private readonly JsonSerializerOptions options = new JsonSerializerOptions { WriteIndented = true };

    public Organizer() => Load();

    private void Load()
    {
        if (!File.Exists(dataFile)) return;
        string json = File.ReadAllText(dataFile);
        events = JsonSerializer.Deserialize<List<Event>>(json) ?? new List<Event>();
    }

    private void Save()
    {
        string json = JsonSerializer.Serialize(events, options);
        File.WriteAllText(dataFile, json);
    }

    private Event GetEvent(string id) => events.FirstOrDefault(e => e.Id == id);
    private Guest GetGuest(Event e, string id) => e.Guests.FirstOrDefault(g => g.Id == id);

    public void Create(string name, string date, string venue)
    {
        var e = new Event(name, date, venue);
        events.Add(e);
        Save();
        Console.WriteLine($"✅ Event created: {e.Name} (ID: {e.Id})");
    }

    public void AddGuest(string eventId, string name, string email, bool plusOne, string dietary, string notes)
    {
        var e = GetEvent(eventId);
        if (e == null) {
            Console.WriteLine($"Event {eventId} not found.");
            return;
        }
        var g = new Guest(name, email, plusOne, dietary, notes);
        e.Guests.Add(g);
        Save();
        Console.WriteLine($"✅ Guest added: {g.Name} (ID: {g.Id})");
    }

    public void ListGuests(string eventId, string status)
    {
        var e = GetEvent(eventId);
        if (e == null) {
            Console.WriteLine($"Event {eventId} not found.");
            return;
        }
        var guests = e.Guests;
        if (!string.IsNullOrEmpty(status)) {
            guests = guests.Where(g => g.RSVP == status).ToList();
        }
        if (!guests.Any()) {
            Console.WriteLine("No guests.");
            return;
        }
        Console.WriteLine($"\n📋 Event: {e.Name}");
        if (!string.IsNullOrEmpty(e.Date)) Console.WriteLine($"Date: {e.Date}");
        if (!string.IsNullOrEmpty(e.Venue)) Console.WriteLine($"Venue: {e.Venue}");
        Console.WriteLine($"\n👤 Guests ({guests.Count}):");
        for (int i = 0; i < guests.Count; i++) {
            var g = guests[i];
            string plus = g.PlusOne ? " (plus‑one: Yes)" : " (plus‑one: No)";
            string dietary = !string.IsNullOrEmpty(g.Dietary) ? $" 🥗 {g.Dietary}" : "";
            string notes = !string.IsNullOrEmpty(g.Notes) ? $" 📝 {g.Notes}" : "";
            Console.WriteLine($"  {i+1}. {g.Name} ({g.Email}) – {g.RSVP}{plus}{dietary}{notes}");
        }
    }

    public void Rsvp(string eventId, string guestId, string status)
    {
        var e = GetEvent(eventId);
        if (e == null) {
            Console.WriteLine($"Event {eventId} not found.");
            return;
        }
        var g = GetGuest(e, guestId);
        if (g == null) {
            Console.WriteLine($"Guest {guestId} not found.");
            return;
        }
        string[] validStatuses = {"Pending", "Attending", "Declined"};
        if (!validStatuses.Contains(status)) {
            Console.WriteLine($"Invalid status. Choose: {string.Join(", ", validStatuses)}");
            return;
        }
        g.RSVP = status;
        Save();
        Console.WriteLine($"✅ {g.Name} RSVP updated to: {status}");
    }

    public void Stats(string eventId)
    {
        var e = GetEvent(eventId);
        if (e == null) {
            Console.WriteLine($"Event {eventId} not found.");
            return;
        }
        int total = e.Guests.Count;
        int attending = e.Guests.Count(g => g.RSVP == "Attending");
        int pending = e.Guests.Count(g => g.RSVP == "Pending");
        int declined = e.Guests.Count(g => g.RSVP == "Declined");
        int plusOnes = e.Guests.Count(g => g.PlusOne);
        Console.WriteLine($"\n📊 Event: {e.Name}");
        Console.WriteLine($"  Total guests: {total}");
        Console.WriteLine($"  Attending: {attending}");
        Console.WriteLine($"  Pending: {pending}");
        Console.WriteLine($"  Declined: {declined}");
        Console.WriteLine($"  Plus‑ones: {plusOnes}");
    }

    public void Search(string eventId, string term)
    {
        var e = GetEvent(eventId);
        if (e == null) {
            Console.WriteLine($"Event {eventId} not found.");
            return;
        }
        var lower = term.ToLower();
        var results = e.Guests.Where(g =>
            g.Name.ToLower().Contains(lower) ||
            g.Email.ToLower().Contains(lower)
        ).ToList();
        if (!results.Any()) {
            Console.WriteLine("No matches.");
            return;
        }
        Console.WriteLine($"\n🔍 Found {results.Count} guest(s):");
        for (int i = 0; i < results.Count; i++) {
            Console.WriteLine($"  {i+1}. {results[i].Name} ({results[i].Email}) – {results[i].RSVP}");
        }
    }

    public void Export(string eventId, string filename)
    {
        var e = GetEvent(eventId);
        if (e == null) {
            Console.WriteLine($"Event {eventId} not found.");
            return;
        }
        using var writer = new StreamWriter(filename);
        writer.WriteLine("ID,Name,Email,RSVP,Plus-One,Dietary,Notes,Added");
        foreach (var g in e.Guests) {
            writer.WriteLine($"{g.Id},{g.Name},{g.Email},{g.RSVP},{g.PlusOne},{g.Dietary},{g.Notes},{g.Added}");
        }
        Console.WriteLine($"✅ Exported {e.Guests.Count} guests to {filename}");
    }

    static void Main(string[] args)
    {
        if (args.Length < 1) {
            Console.WriteLine("Usage: EventGuests <command> [options]");
            return;
        }
        var app = new Organizer();
        var cmd = args[0];
        var parsed = ParseArgs(args);
        switch (cmd) {
            case "create":
                if (args.Length < 2) { Console.WriteLine("create <name> [--date DATE] [--venue VENUE]"); return; }
                app.Create(args[1], parsed.GetValueOrDefault("date", ""), parsed.GetValueOrDefault("venue", ""));
                break;
            case "add-guest":
                if (args.Length < 4) { Console.WriteLine("add-guest <event_id> <name> <email> [--plus-one] [--dietary DIET] [--notes NOTES]"); return; }
                app.AddGuest(args[1], args[2], args[3],
                    parsed.ContainsKey("plus-one"),
                    parsed.GetValueOrDefault("dietary", ""),
                    parsed.GetValueOrDefault("notes", ""));
                break;
            case "list-guests":
                if (args.Length < 2) { Console.WriteLine("list-guests <event_id> [--status STATUS]"); return; }
                app.ListGuests(args[1], parsed.GetValueOrDefault("status", ""));
                break;
            case "rsvp":
                if (args.Length < 4) { Console.WriteLine("rsvp <event_id> <guest_id> <status>"); return; }
                app.Rsvp(args[1], args[2], args[3]);
                break;
            case "stats":
                if (args.Length < 2) { Console.WriteLine("stats <event_id>"); return; }
                app.Stats(args[1]);
                break;
            case "search":
                if (args.Length < 3) { Console.WriteLine("search <event_id> <term>"); return; }
                app.Search(args[1], args[2]);
                break;
            case "export":
                if (args.Length < 2) { Console.WriteLine("export <event_id> [--filename FILE]"); return; }
                app.Export(args[1], parsed.GetValueOrDefault("filename", "guests.csv"));
                break;
            default:
                Console.WriteLine("Unknown command.");
                break;
        }
    }

    static Dictionary<string, string> ParseArgs(string[] args)
    {
        var dict = new Dictionary<string, string>();
        for (int i = 1; i < args.Length; i++) {
            if (args[i].StartsWith("--") && i + 1 < args.Length) {
                dict[args[i].Substring(2)] = args[++i];
            }
        }
        return dict;
    }
}
