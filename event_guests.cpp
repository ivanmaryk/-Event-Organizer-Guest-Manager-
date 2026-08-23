// event_guests.cpp
#include <iostream>
#include <fstream>
#include <string>
#include <vector>
#include <map>
#include <algorithm>
#include <cctype>
#include <random>
#include <nlohmann/json.hpp>

using namespace std;
using json = nlohmann::json;

string generateId() {
    const char* hex = "0123456789abcdef";
    string id;
    random_device rd;
    mt19937 gen(rd());
    uniform_int_distribution<> dis(0, 15);
    for (int i=0; i<8; i++) id += hex[dis(gen)];
    return id;
}

string currentTime() {
    time_t t = time(nullptr);
    char buf[30];
    strftime(buf, sizeof(buf), "%Y-%m-%dT%H:%M:%S%z", localtime(&t));
    return string(buf);
}

struct Guest {
    string id, name, email, rsvp, dietary, notes, added;
    bool plus_one;
};

struct Event {
    string id, name, date, venue, created;
    vector<Guest> guests;
};

class Organizer {
private:
    vector<Event> events;
    string dataFile = "events.json";

    void load() {
        ifstream f(dataFile);
        if (!f.is_open()) return;
        json j;
        f >> j;
        for (auto& item : j) {
            Event e;
            e.id = item["id"];
            e.name = item["name"];
            e.date = item["date"];
            e.venue = item["venue"];
            e.created = item["created"];
            for (auto& g : item["guests"]) {
                Guest guest;
                guest.id = g["id"];
                guest.name = g["name"];
                guest.email = g["email"];
                guest.rsvp = g["rsvp"];
                guest.plus_one = g["plus_one"];
                guest.dietary = g["dietary"];
                guest.notes = g["notes"];
                guest.added = g["added"];
                e.guests.push_back(guest);
            }
            events.push_back(e);
        }
    }

    void save() {
        json j = json::array();
        for (auto& e : events) {
            json guests = json::array();
            for (auto& g : e.guests) {
                guests.push_back({
                    {"id", g.id}, {"name", g.name}, {"email", g.email},
                    {"rsvp", g.rsvp}, {"plus_one", g.plus_one},
                    {"dietary", g.dietary}, {"notes", g.notes}, {"added", g.added}
                });
            }
            j.push_back({
                {"id", e.id}, {"name", e.name}, {"date", e.date},
                {"venue", e.venue}, {"guests", guests}, {"created", e.created}
            });
        }
        ofstream f(dataFile);
        f << setw(2) << j << endl;
    }

    Event* getEvent(const string& id) {
        for (auto& e : events) {
            if (e.id == id) return &e;
        }
        return nullptr;
    }

    Guest* getGuest(Event* e, const string& id) {
        for (auto& g : e->guests) {
            if (g.id == id) return &g;
        }
        return nullptr;
    }

public:
    Organizer() { load(); }

    void create(const string& name, const string& date, const string& venue) {
        Event e;
        e.id = generateId();
        e.name = name;
        e.date = date;
        e.venue = venue;
        e.created = currentTime();
        events.push_back(e);
        save();
        cout << "✅ Event created: " << e.name << " (ID: " << e.id << ")\n";
    }

    void addGuest(const string& eventId, const string& name, const string& email,
                  bool plusOne, const string& dietary, const string& notes) {
        Event* e = getEvent(eventId);
        if (!e) {
            cout << "Event " << eventId << " not found.\n";
            return;
        }
        Guest g;
        g.id = generateId();
        g.name = name;
        g.email = email;
        g.rsvp = "Pending";
        g.plus_one = plusOne;
        g.dietary = dietary;
        g.notes = notes;
        g.added = currentTime();
        e->guests.push_back(g);
        save();
        cout << "✅ Guest added: " << g.name << " (ID: " << g.id << ")\n";
    }

    void listGuests(const string& eventId, const string& status) {
        Event* e = getEvent(eventId);
        if (!e) {
            cout << "Event " << eventId << " not found.\n";
            return;
        }
        vector<Guest> guests = e->guests;
        if (!status.empty()) {
            vector<Guest> filtered;
            for (auto& g : guests) {
                if (g.rsvp == status) filtered.push_back(g);
            }
            guests = filtered;
        }
        if (guests.empty()) {
            cout << "No guests.\n";
            return;
        }
        cout << "\n📋 Event: " << e->name << "\n";
        if (!e->date.empty()) cout << "Date: " << e->date << "\n";
        if (!e->venue.empty()) cout << "Venue: " << e->venue << "\n";
        cout << "\n👤 Guests (" << guests.size() << "):\n";
        for (size_t i=0; i<guests.size(); i++) {
            auto& g = guests[i];
            string plus = g.plus_one ? " (plus‑one: Yes)" : " (plus‑one: No)";
            string dietary = g.dietary.empty() ? "" : " 🥗 " + g.dietary;
            string notes = g.notes.empty() ? "" : " 📝 " + g.notes;
            cout << "  " << i+1 << ". " << g.name << " (" << g.email << ") – " << g.rsvp << plus << dietary << notes << "\n";
        }
    }

    void rsvp(const string& eventId, const string& guestId, const string& status) {
        Event* e = getEvent(eventId);
        if (!e) {
            cout << "Event " << eventId << " not found.\n";
            return;
        }
        Guest* g = getGuest(e, guestId);
        if (!g) {
            cout << "Guest " << guestId << " not found.\n";
            return;
        }
        vector<string> valid = {"Pending", "Attending", "Declined"};
        if (find(valid.begin(), valid.end(), status) == valid.end()) {
            cout << "Invalid status. Choose: Pending, Attending, Declined\n";
            return;
        }
        g->rsvp = status;
        save();
        cout << "✅ " << g->name << " RSVP updated to: " << status << "\n";
    }

    void stats(const string& eventId) {
        Event* e = getEvent(eventId);
        if (!e) {
            cout << "Event " << eventId << " not found.\n";
            return;
        }
        int total = e->guests.size();
        int attending=0, pending=0, declined=0, plusOnes=0;
        for (auto& g : e->guests) {
            if (g.rsvp == "Attending") attending++;
            else if (g.rsvp == "Pending") pending++;
            else if (g.rsvp == "Declined") declined++;
            if (g.plus_one) plusOnes++;
        }
        cout << "\n📊 Event: " << e->name << "\n";
        cout << "  Total guests: " << total << "\n";
        cout << "  Attending: " << attending << "\n";
        cout << "  Pending: " << pending << "\n";
        cout << "  Declined: " << declined << "\n";
        cout << "  Plus‑ones: " << plusOnes << "\n";
    }

    void search(const string& eventId, const string& term) {
        Event* e = getEvent(eventId);
        if (!e) {
            cout << "Event " << eventId << " not found.\n";
            return;
        }
        string t = term;
        transform(t.begin(), t.end(), t.begin(), ::tolower);
        vector<Guest> results;
        for (auto& g : e->guests) {
            string n = g.name, em = g.email;
            transform(n.begin(), n.end(), n.begin(), ::tolower);
            transform(em.begin(), em.end(), em.begin(), ::tolower);
            if (n.find(t) != string::npos || em.find(t) != string::npos) {
                results.push_back(g);
            }
        }
        if (results.empty()) {
            cout << "No matches.\n";
            return;
        }
        cout << "\n🔍 Found " << results.size() << " guest(s):\n";
        for (size_t i=0; i<results.size(); i++) {
            cout << "  " << i+1 << ". " << results[i].name << " (" << results[i].email << ") – " << results[i].rsvp << "\n";
        }
    }

    void exportCSV(const string& eventId, const string& filename) {
        Event* e = getEvent(eventId);
        if (!e) {
            cout << "Event " << eventId << " not found.\n";
            return;
        }
        ofstream f(filename);
        f << "ID,Name,Email,RSVP,Plus-One,Dietary,Notes,Added\n";
        for (auto& g : e->guests) {
            f << g.id << "," << g.name << "," << g.email << "," << g.rsvp << ","
              << (g.plus_one ? "true" : "false") << "," << g.dietary << "," << g.notes << "," << g.added << "\n";
        }
        f.close();
        cout << "✅ Exported " << e->guests.size() << " guests to " << filename << "\n";
    }
};

int main(int argc, char* argv[]) {
    if (argc < 2) {
        cerr << "Usage: event_guests <command> [options]\n";
        return 1;
    }
    Organizer app;
    string cmd = argv[1];

    if (cmd == "create") {
        if (argc < 3) { cerr << "create <name> [--date DATE] [--venue VENUE]\n"; return 1; }
        string name = argv[2];
        string date, venue;
        for (int i=3; i<argc; i++) {
            if (string(argv[i]) == "--date" && i+1 < argc) date = argv[++i];
            if (string(argv[i]) == "--venue" && i+1 < argc) venue = argv[++i];
        }
        app.create(name, date, venue);
    } else if (cmd == "add-guest") {
        if (argc < 5) { cerr << "add-guest <event_id> <name> <email> [--plus-one] [--dietary DIET] [--notes NOTES]\n"; return 1; }
        string eventId = argv[2];
        string name = argv[3];
        string email = argv[4];
        bool plusOne = false;
        string dietary, notes;
        for (int i=5; i<argc; i++) {
            if (string(argv[i]) == "--plus-one") plusOne = true;
            if (string(argv[i]) == "--dietary" && i+1 < argc) dietary = argv[++i];
            if (string(argv[i]) == "--notes" && i+1 < argc) notes = argv[++i];
        }
        app.addGuest(eventId, name, email, plusOne, dietary, notes);
    } else if (cmd == "list-guests") {
        if (argc < 3) { cerr << "list-guests <event_id> [--status STATUS]\n"; return 1; }
        string eventId = argv[2];
        string status;
        for (int i=3; i<argc; i++) {
            if (string(argv[i]) == "--status" && i+1 < argc) status = argv[++i];
        }
        app.listGuests(eventId, status);
    } else if (cmd == "rsvp") {
        if (argc < 5) { cerr << "rsvp <event_id> <guest_id> <status>\n"; return 1; }
        app.rsvp(argv[2], argv[3], argv[4]);
    } else if (cmd == "stats") {
        if (argc < 3) { cerr << "stats <event_id>\n"; return 1; }
        app.stats(argv[2]);
    } else if (cmd == "search") {
        if (argc < 4) { cerr << "search <event_id> <term>\n"; return 1; }
        app.search(argv[2], argv[3]);
    } else if (cmd == "export") {
        if (argc < 3) { cerr << "export <event_id> [--filename FILE]\n"; return 1; }
        string eventId = argv[2];
        string filename = "guests.csv";
        for (int i=3; i<argc; i++) {
            if (string(argv[i]) == "--filename" && i+1 < argc) filename = argv[++i];
        }
        app.exportCSV(eventId, filename);
    } else {
        cerr << "Unknown command. Use create, add-guest, list-guests, rsvp, stats, search, export.\n";
        return 1;
    }
    return 0;
}
