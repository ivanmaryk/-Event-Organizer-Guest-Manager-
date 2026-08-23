# event_guests.php
#!/usr/bin/env php
<?php

define('DATA_FILE', 'events.json');

class Guest {
    public $id;
    public $name;
    public $email;
    public $rsvp;
    public $plus_one;
    public $dietary;
    public $notes;
    public $added;

    function __construct($name, $email, $plus_one = false, $dietary = '', $notes = '') {
        $this->id = substr(bin2hex(random_bytes(4)), 0, 8);
        $this->name = $name;
        $this->email = $email;
        $this->rsvp = 'Pending';
        $this->plus_one = $plus_one;
        $this->dietary = $dietary;
        $this->notes = $notes;
        $this->added = date('c');
    }

    function toArray() {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'email' => $this->email,
            'rsvp' => $this->rsvp,
            'plus_one' => $this->plus_one,
            'dietary' => $this->dietary,
            'notes' => $this->notes,
            'added' => $this->added
        ];
    }

    static function fromArray($data) {
        $g = new self($data['name'], $data['email'], $data['plus_one'], $data['dietary'], $data['notes']);
        $g->id = $data['id'];
        $g->rsvp = $data['rsvp'] ?? 'Pending';
        $g->added = $data['added'] ?? date('c');
        return $g;
    }
}

class Event {
    public $id;
    public $name;
    public $date;
    public $venue;
    public $guests = [];
    public $created;

    function __construct($name, $date = '', $venue = '') {
        $this->id = substr(bin2hex(random_bytes(4)), 0, 8);
        $this->name = $name;
        $this->date = $date;
        $this->venue = $venue;
        $this->created = date('c');
    }

    function toArray() {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'date' => $this->date,
            'venue' => $this->venue,
            'guests' => array_map(function($g) { return $g->toArray(); }, $this->guests),
            'created' => $this->created
        ];
    }

    static function fromArray($data) {
        $e = new self($data['name'], $data['date'], $data['venue']);
        $e->id = $data['id'];
        $e->guests = array_map(function($g) { return Guest::fromArray($g); }, $data['guests']);
        $e->created = $data['created'] ?? date('c');
        return $e;
    }
}

class Organizer {
    private $events = [];

    function __construct() {
        $this->load();
    }

    function load() {
        if (file_exists(DATA_FILE)) {
            $data = json_decode(file_get_contents(DATA_FILE), true);
            $this->events = array_map(function($e) { return Event::fromArray($e); }, $data);
        }
    }

    function save() {
        $data = array_map(function($e) { return $e->toArray(); }, $this->events);
        file_put_contents(DATA_FILE, json_encode($data, JSON_PRETTY_PRINT));
    }

    function getEvent($id) {
        foreach ($this->events as $e) {
            if ($e->id == $id) return $e;
        }
        return null;
    }

    function getGuest($event, $id) {
        foreach ($event->guests as $g) {
            if ($g->id == $id) return $g;
        }
        return null;
    }

    function create($name, $date = '', $venue = '') {
        $e = new Event($name, $date, $venue);
        $this->events[] = $e;
        $this->save();
        echo "✅ Event created: {$e->name} (ID: {$e->id})\n";
    }

    function addGuest($eventId, $name, $email, $plusOne = false, $dietary = '', $notes = '') {
        $e = $this->getEvent($eventId);
        if (!$e) {
            echo "Event $eventId not found.\n";
            return;
        }
        $g = new Guest($name, $email, $plusOne, $dietary, $notes);
        $e->guests[] = $g;
        $this->save();
        echo "✅ Guest added: {$g->name} (ID: {$g->id})\n";
    }

    function listGuests($eventId, $status = null) {
        $e = $this->getEvent($eventId);
        if (!$e) {
            echo "Event $eventId not found.\n";
            return;
        }
        $guests = $e->guests;
        if ($status) {
            $guests = array_filter($guests, function($g) use ($status) {
                return $g->rsvp == $status;
            });
        }
        if (empty($guests)) {
            echo "No guests.\n";
            return;
        }
        echo "\n📋 Event: {$e->name}\n";
        if ($e->date) echo "Date: {$e->date}\n";
        if ($e->venue) echo "Venue: {$e->venue}\n";
        echo "\n👤 Guests (" . count($guests) . "):\n";
        $i = 1;
        foreach ($guests as $g) {
            $plus = $g->plus_one ? ' (plus‑one: Yes)' : ' (plus‑one: No)';
            $dietary = $g->dietary ? " 🥗 {$g->dietary}" : '';
            $notes = $g->notes ? " 📝 {$g->notes}" : '';
            echo "  $i. {$g->name} ({$g->email}) – {$g->rsvp}{$plus}{$dietary}{$notes}\n";
            $i++;
        }
    }

    function rsvp($eventId, $guestId, $status) {
        $e = $this->getEvent($eventId);
        if (!$e) {
            echo "Event $eventId not found.\n";
            return;
        }
        $g = $this->getGuest($e, $guestId);
        if (!$g) {
            echo "Guest $guestId not found.\n";
            return;
        }
        $validStatuses = ['Pending', 'Attending', 'Declined'];
        if (!in_array($status, $validStatuses)) {
            echo "Invalid status. Choose: " . implode(', ', $validStatuses) . "\n";
            return;
        }
        $g->rsvp = $status;
        $this->save();
        echo "✅ {$g->name} RSVP updated to: $status\n";
    }

    function stats($eventId) {
        $e = $this->getEvent($eventId);
        if (!$e) {
            echo "Event $eventId not found.\n";
            return;
        }
        $total = count($e->guests);
        $attending = 0; $pending = 0; $declined = 0; $plusOnes = 0;
        foreach ($e->guests as $g) {
            if ($g->rsvp == 'Attending') $attending++;
            if ($g->rsvp == 'Pending') $pending++;
            if ($g->rsvp == 'Declined') $declined++;
            if ($g->plus_one) $plusOnes++;
        }
        echo "\n📊 Event: {$e->name}\n";
        echo "  Total guests: $total\n";
        echo "  Attending: $attending\n";
        echo "  Pending: $pending\n";
        echo "  Declined: $declined\n";
        echo "  Plus‑ones: $plusOnes\n";
    }

    function search($eventId, $term) {
        $e = $this->getEvent($eventId);
        if (!$e) {
            echo "Event $eventId not found.\n";
            return;
        }
        $lower = strtolower($term);
        $results = array_filter($e->guests, function($g) use ($lower) {
            return strpos(strtolower($g->name), $lower) !== false ||
                   strpos(strtolower($g->email), $lower) !== false;
        });
        if (empty($results)) {
            echo "No matches.\n";
            return;
        }
        echo "\n🔍 Found " . count($results) . " guest(s):\n";
        $i = 1;
        foreach ($results as $g) {
            echo "  $i. {$g->name} ({$g->email}) – {$g->rsvp}\n";
            $i++;
        }
    }

    function export($eventId, $filename) {
        $e = $this->getEvent($eventId);
        if (!$e) {
            echo "Event $eventId not found.\n";
            return;
        }
        $fp = fopen($filename, 'w');
        fputcsv($fp, ['ID', 'Name', 'Email', 'RSVP', 'Plus-One', 'Dietary', 'Notes', 'Added']);
        foreach ($e->guests as $g) {
            fputcsv($fp, [$g->id, $g->name, $g->email, $g->rsvp, $g->plus_one, $g->dietary, $g->notes, $g->added]);
        }
        fclose($fp);
        echo "✅ Exported " . count($e->guests) . " guests to $filename\n";
    }
}

if ($argc < 2) {
    die("Usage: php event_guests.php <command> [options]\n");
}
$app = new Organizer();
$cmd = $argv[1];

switch ($cmd) {
    case 'create':
        if ($argc < 3) die("create <name> [--date DATE] [--venue VENUE]\n");
        $name = $argv[2];
        $date = ''; $venue = '';
        for ($i=3; $i<$argc; $i++) {
            if ($argv[$i] == '--date' && isset($argv[$i+1])) { $date = $argv[++$i]; }
            if ($argv[$i] == '--venue' && isset($argv[$i+1])) { $venue = $argv[++$i]; }
        }
        $app->create($name, $date, $venue);
        break;

    case 'add-guest':
        if ($argc < 5) die("add-guest <event_id> <name> <email> [--plus-one] [--dietary DIET] [--notes NOTES]\n");
        $eventId = $argv[2];
        $name = $argv[3];
        $email = $argv[4];
        $plusOne = false; $dietary = ''; $notes = '';
        for ($i=5; $i<$argc; $i++) {
            if ($argv[$i] == '--plus-one') { $plusOne = true; }
            if ($argv[$i] == '--dietary' && isset($argv[$i+1])) { $dietary = $argv[++$i]; }
            if ($argv[$i] == '--notes' && isset($argv[$i+1])) { $notes = $argv[++$i]; }
        }
        $app->addGuest($eventId, $name, $email, $plusOne, $dietary, $notes);
        break;

    case 'list-guests':
        if ($argc < 3) die("list-guests <event_id> [--status STATUS]\n");
        $eventId = $argv[2];
        $status = null;
        for ($i=3; $i<$argc; $i++) {
            if ($argv[$i] == '--status' && isset($argv[$i+1])) { $status = $argv[++$i]; }
        }
        $app->listGuests($eventId, $status);
        break;

    case 'rsvp':
        if ($argc < 5) die("rsvp <event_id> <guest_id> <status>\n");
        $app->rsvp($argv[2], $argv[3], $argv[4]);
        break;

    case 'stats':
        if ($argc < 3) die("stats <event_id>\n");
        $app->stats($argv[2]);
        break;

    case 'search':
        if ($argc < 4) die("search <event_id> <term>\n");
        $app->search($argv[2], $argv[3]);
        break;

    case 'export':
        if ($argc < 3) die("export <event_id> [--filename FILE]\n");
        $eventId = $argv[2];
        $filename = 'guests.csv';
        for ($i=3; $i<$argc; $i++) {
            if ($argv[$i] == '--filename' && isset($argv[$i+1])) { $filename = $argv[++$i]; }
        }
        $app->export($eventId, $filename);
        break;

    default:
        echo "Unknown command. Use create, add-guest, list-guests, rsvp, stats, search, export.\n";
}
?>
