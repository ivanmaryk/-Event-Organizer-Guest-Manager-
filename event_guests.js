// event_guests.js
#!/usr/bin/env node
const fs = require('fs');
const { program } = require('commander');
const { v4: uuidv4 } = require('uuid');

const DATA_FILE = 'events.json';

class Guest {
    constructor(name, email, plusOne = false, dietary = '', notes = '') {
        this.id = uuidv4().slice(0,8);
        this.name = name;
        this.email = email;
        this.rsvp = 'Pending';
        this.plusOne = plusOne;
        this.dietary = dietary;
        this.notes = notes;
        this.added = new Date().toISOString();
    }
}

class Event {
    constructor(name, date = '', venue = '') {
        this.id = uuidv4().slice(0,8);
        this.name = name;
        this.date = date;
        this.venue = venue;
        this.guests = [];
        this.created = new Date().toISOString();
    }
}

class Organizer {
    constructor() {
        this.events = [];
        this.load();
    }

    load() {
        if (fs.existsSync(DATA_FILE)) {
            try {
                const data = JSON.parse(fs.readFileSync(DATA_FILE));
                this.events = data.map(e => {
                    const ev = new Event(e.name, e.date, e.venue);
                    ev.id = e.id;
                    ev.guests = e.guests.map(g => {
                        const guest = new Guest(g.name, g.email, g.plusOne, g.dietary, g.notes);
                        guest.id = g.id;
                        guest.rsvp = g.rsvp || 'Pending';
                        guest.added = g.added;
                        return guest;
                    });
                    ev.created = e.created;
                    return ev;
                });
            } catch (e) {}
        }
    }

    save() {
        fs.writeFileSync(DATA_FILE, JSON.stringify(this.events, null, 2));
    }

    getEvent(id) {
        return this.events.find(e => e.id === id);
    }

    getGuest(event, id) {
        return event.guests.find(g => g.id === id);
    }

    create(name, date, venue) {
        const e = new Event(name, date, venue);
        this.events.push(e);
        this.save();
        console.log(`✅ Event created: ${e.name} (ID: ${e.id})`);
    }

    addGuest(eventId, name, email, plusOne, dietary, notes) {
        const e = this.getEvent(eventId);
        if (!e) {
            console.log(`Event ${eventId} not found.`);
            return;
        }
        const g = new Guest(name, email, plusOne, dietary, notes);
        e.guests.push(g);
        this.save();
        console.log(`✅ Guest added: ${g.name} (ID: ${g.id})`);
    }

    listGuests(eventId, status) {
        const e = this.getEvent(eventId);
        if (!e) {
            console.log(`Event ${eventId} not found.`);
            return;
        }
        let guests = e.guests;
        if (status) {
            guests = guests.filter(g => g.rsvp === status);
        }
        if (guests.length === 0) {
            console.log('No guests.');
            return;
        }
        console.log(`\n📋 Event: ${e.name}`);
        if (e.date) console.log(`Date: ${e.date}`);
        if (e.venue) console.log(`Venue: ${e.venue}`);
        console.log(`\n👤 Guests (${guests.length}):`);
        guests.forEach((g, i) => {
            const plus = g.plusOne ? ' (plus‑one: Yes)' : ' (plus‑one: No)';
            const dietary = g.dietary ? ` 🥗 ${g.dietary}` : '';
            const notes = g.notes ? ` 📝 ${g.notes}` : '';
            console.log(`  ${i+1}. ${g.name} (${g.email}) – ${g.rsvp}${plus}${dietary}${notes}`);
        });
    }

    rsvp(eventId, guestId, status) {
        const e = this.getEvent(eventId);
        if (!e) {
            console.log(`Event ${eventId} not found.`);
            return;
        }
        const g = this.getGuest(e, guestId);
        if (!g) {
            console.log(`Guest ${guestId} not found.`);
            return;
        }
        const validStatuses = ['Pending', 'Attending', 'Declined'];
        if (!validStatuses.includes(status)) {
            console.log(`Invalid status. Choose: ${validStatuses.join(', ')}`);
            return;
        }
        g.rsvp = status;
        this.save();
        console.log(`✅ ${g.name} RSVP updated to: ${status}`);
    }

    stats(eventId) {
        const e = this.getEvent(eventId);
        if (!e) {
            console.log(`Event ${eventId} not found.`);
            return;
        }
        const total = e.guests.length;
        const attending = e.guests.filter(g => g.rsvp === 'Attending').length;
        const pending = e.guests.filter(g => g.rsvp === 'Pending').length;
        const declined = e.guests.filter(g => g.rsvp === 'Declined').length;
        const plusOnes = e.guests.filter(g => g.plusOne).length;
        console.log(`\n📊 Event: ${e.name}`);
        console.log(`  Total guests: ${total}`);
        console.log(`  Attending: ${attending}`);
        console.log(`  Pending: ${pending}`);
        console.log(`  Declined: ${declined}`);
        console.log(`  Plus‑ones: ${plusOnes}`);
    }

    search(eventId, term) {
        const e = this.getEvent(eventId);
        if (!e) {
            console.log(`Event ${eventId} not found.`);
            return;
        }
        const lower = term.toLowerCase();
        const results = e.guests.filter(g =>
            g.name.toLowerCase().includes(lower) ||
            g.email.toLowerCase().includes(lower)
        );
        if (results.length === 0) {
            console.log('No matches.');
            return;
        }
        console.log(`\n🔍 Found ${results.length} guest(s):`);
        results.forEach((g, i) => {
            console.log(`  ${i+1}. ${g.name} (${g.email}) – ${g.rsvp}`);
        });
    }

    export(eventId, filename) {
        const e = this.getEvent(eventId);
        if (!e) {
            console.log(`Event ${eventId} not found.`);
            return;
        }
        let csv = 'ID,Name,Email,RSVP,Plus-One,Dietary,Notes,Added\n';
        for (const g of e.guests) {
            csv += `${g.id},${g.name},${g.email},${g.rsvp},${g.plusOne},${g.dietary},${g.notes},${g.added}\n`;
        }
        fs.writeFileSync(filename, csv);
        console.log(`✅ Exported ${e.guests.length} guests to ${filename}`);
    }
}

program
    .command('create <name>')
    .option('--date <date>', 'Event date')
    .option('--venue <venue>', 'Event venue')
    .action((name, options) => {
        const app = new Organizer();
        app.create(name, options.date || '', options.venue || '');
    });

program
    .command('add-guest <eventId> <name> <email>')
    .option('--plus-one', 'Guest can bring a plus-one')
    .option('--dietary <diet>', 'Dietary restrictions')
    .option('--notes <notes>', 'Additional notes')
    .action((eventId, name, email, options) => {
        const app = new Organizer();
        app.addGuest(eventId, name, email, options.plusOne || false, options.dietary || '', options.notes || '');
    });

program
    .command('list-guests <eventId>')
    .option('--status <status>', 'Filter by status (Pending, Attending, Declined)')
    .action((eventId, options) => {
        const app = new Organizer();
        app.listGuests(eventId, options.status || '');
    });

program
    .command('rsvp <eventId> <guestId> <status>')
    .action((eventId, guestId, status) => {
        const app = new Organizer();
        app.rsvp(eventId, guestId, status);
    });

program
    .command('stats <eventId>')
    .action((eventId) => {
        const app = new Organizer();
        app.stats(eventId);
    });

program
    .command('search <eventId> <term>')
    .action((eventId, term) => {
        const app = new Organizer();
        app.search(eventId, term);
    });

program
    .command('export <eventId>')
    .option('--filename <filename>', 'Output filename', 'guests.csv')
    .action((eventId, options) => {
        const app = new Organizer();
        app.export(eventId, options.filename);
    });

program.parse(process.argv);
