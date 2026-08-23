# event_guests.rb
#!/usr/bin/env ruby
require 'json'
require 'securerandom'
require 'date'
require 'csv'

DATA_FILE = 'events.json'

class Guest
  attr_accessor :id, :name, :email, :rsvp, :plus_one, :dietary, :notes, :added

  def initialize(name, email, plus_one = false, dietary = '', notes = '')
    @id = SecureRandom.hex(4)
    @name = name
    @email = email
    @rsvp = 'Pending'
    @plus_one = plus_one
    @dietary = dietary
    @notes = notes
    @added = Time.now.iso8601
  end

  def to_hash
    { id: @id, name: @name, email: @email, rsvp: @rsvp,
      plus_one: @plus_one, dietary: @dietary, notes: @notes, added: @added }
  end

  def self.from_hash(h)
    g = new(h['name'], h['email'], h['plus_one'], h['dietary'], h['notes'])
    g.id = h['id']
    g.rsvp = h['rsvp'] || 'Pending'
    g.added = h['added'] || Time.now.iso8601
    g
  end
end

class Event
  attr_accessor :id, :name, :date, :venue, :guests, :created

  def initialize(name, date = '', venue = '')
    @id = SecureRandom.hex(4)
    @name = name
    @date = date
    @venue = venue
    @guests = []
    @created = Time.now.iso8601
  end

  def to_hash
    { id: @id, name: @name, date: @date, venue: @venue,
      guests: @guests.map(&:to_hash), created: @created }
  end

  def self.from_hash(h)
    e = new(h['name'], h['date'], h['venue'])
    e.id = h['id']
    e.guests = h['guests'].map { |g| Guest.from_hash(g) }
    e.created = h['created'] || Time.now.iso8601
    e
  end
end

class Organizer
  attr_reader :events

  def initialize
    @events = []
    load
  end

  def load
    if File.exist?(DATA_FILE)
      data = JSON.parse(File.read(DATA_FILE))
      @events = data.map { |h| Event.from_hash(h) }
    end
  end

  def save
    File.write(DATA_FILE, JSON.pretty_generate(@events.map(&:to_hash)))
  end

  def get_event(id)
    @events.find { |e| e.id == id }
  end

  def get_guest(event, id)
    event.guests.find { |g| g.id == id }
  end

  def create(name, date = '', venue = '')
    e = Event.new(name, date, venue)
    @events << e
    save
    puts "✅ Event created: #{e.name} (ID: #{e.id})"
  end

  def add_guest(event_id, name, email, plus_one = false, dietary = '', notes = '')
    e = get_event(event_id)
    unless e
      puts "Event #{event_id} not found."
      return
    end
    g = Guest.new(name, email, plus_one, dietary, notes)
    e.guests << g
    save
    puts "✅ Guest added: #{g.name} (ID: #{g.id})"
  end

  def list_guests(event_id, status = nil)
    e = get_event(event_id)
    unless e
      puts "Event #{event_id} not found."
      return
    end
    guests = e.guests
    guests = guests.select { |g| g.rsvp == status } if status
    if guests.empty?
      puts "No guests."
      return
    end
    puts "\n📋 Event: #{e.name}"
    puts "Date: #{e.date}" unless e.date.empty?
    puts "Venue: #{e.venue}" unless e.venue.empty?
    puts "\n👤 Guests (#{guests.size}):"
    guests.each_with_index do |g, i|
      plus = g.plus_one ? ' (plus‑one: Yes)' : ' (plus‑one: No)'
      dietary = g.dietary.empty? ? '' : " 🥗 #{g.dietary}"
      notes = g.notes.empty? ? '' : " 📝 #{g.notes}"
      puts "  #{i+1}. #{g.name} (#{g.email}) – #{g.rsvp}#{plus}#{dietary}#{notes}"
    end
  end

  def rsvp(event_id, guest_id, status)
    e = get_event(event_id)
    unless e
      puts "Event #{event_id} not found."
      return
    end
    g = get_guest(e, guest_id)
    unless g
      puts "Guest #{guest_id} not found."
      return
    end
    valid_statuses = ['Pending', 'Attending', 'Declined']
    unless valid_statuses.include?(status)
      puts "Invalid status. Choose: #{valid_statuses.join(', ')}"
      return
    end
    g.rsvp = status
    save
    puts "✅ #{g.name} RSVP updated to: #{status}"
  end

  def stats(event_id)
    e = get_event(event_id)
    unless e
      puts "Event #{event_id} not found."
      return
    end
    total = e.guests.size
    attending = e.guests.count { |g| g.rsvp == 'Attending' }
    pending = e.guests.count { |g| g.rsvp == 'Pending' }
    declined = e.guests.count { |g| g.rsvp == 'Declined' }
    plus_ones = e.guests.count(&:plus_one)
    puts "\n📊 Event: #{e.name}"
    puts "  Total guests: #{total}"
    puts "  Attending: #{attending}"
    puts "  Pending: #{pending}"
    puts "  Declined: #{declined}"
    puts "  Plus‑ones: #{plus_ones}"
  end

  def search(event_id, term)
    e = get_event(event_id)
    unless e
      puts "Event #{event_id} not found."
      return
    end
    lower = term.downcase
    results = e.guests.select { |g| g.name.downcase.include?(lower) || g.email.downcase.include?(lower) }
    if results.empty?
      puts "No matches."
      return
    end
    puts "\n🔍 Found #{results.size} guest(s):"
    results.each_with_index do |g, i|
      puts "  #{i+1}. #{g.name} (#{g.email}) – #{g.rsvp}"
    end
  end

  def export(event_id, filename)
    e = get_event(event_id)
    unless e
      puts "Event #{event_id} not found."
      return
    end
    CSV.open(filename, 'w') do |csv|
      csv << ['ID', 'Name', 'Email', 'RSVP', 'Plus-One', 'Dietary', 'Notes', 'Added']
      e.guests.each do |g|
        csv << [g.id, g.name, g.email, g.rsvp, g.plus_one, g.dietary, g.notes, g.added]
      end
    end
    puts "✅ Exported #{e.guests.size} guests to #{filename}"
  end
end

if ARGV.empty?
  puts "Usage: event_guests.rb <command> [options]"
  exit
end

app = Organizer.new
cmd = ARGV.shift

case cmd
when 'create'
  name = ARGV.shift
  if name.nil?
    puts "create <name> [--date DATE] [--venue VENUE]"
    exit
  end
  date = ''
  venue = ''
  while ARGV.any?
    case ARGV[0]
    when '--date'
      ARGV.shift
      date = ARGV.shift || ''
    when '--venue'
      ARGV.shift
      venue = ARGV.shift || ''
    else
      break
    end
  end
  app.create(name, date, venue)

when 'add-guest'
  if ARGV.size < 3
    puts "add-guest <event_id> <name> <email> [--plus-one] [--dietary DIET] [--notes NOTES]"
    exit
  end
  event_id = ARGV.shift
  name = ARGV.shift
  email = ARGV.shift
  plus_one = false
  dietary = ''
  notes = ''
  while ARGV.any?
    case ARGV[0]
    when '--plus-one'
      ARGV.shift
      plus_one = true
    when '--dietary'
      ARGV.shift
      dietary = ARGV.shift || ''
    when '--notes'
      ARGV.shift
      notes = ARGV.shift || ''
    else
      break
    end
  end
  app.add_guest(event_id, name, email, plus_one, dietary, notes)

when 'list-guests'
  event_id = ARGV.shift
  if event_id.nil?
    puts "list-guests <event_id> [--status STATUS]"
    exit
  end
  status = nil
  if ARGV.include?('--status')
    idx = ARGV.index('--status')
    status = ARGV[idx+1] if idx
  end
  app.list_guests(event_id, status)

when 'rsvp'
  if ARGV.size < 3
    puts "rsvp <event_id> <guest_id> <status>"
    exit
  end
  event_id = ARGV.shift
  guest_id = ARGV.shift
  status = ARGV.shift
  app.rsvp(event_id, guest_id, status)

when 'stats'
  event_id = ARGV.shift
  if event_id.nil?
    puts "stats <event_id>"
    exit
  end
  app.stats(event_id)

when 'search'
  if ARGV.size < 2
    puts "search <event_id> <term>"
    exit
  end
  event_id = ARGV.shift
  term = ARGV.shift
  app.search(event_id, term)

when 'export'
  event_id = ARGV.shift
  if event_id.nil?
    puts "export <event_id> [--filename FILE]"
    exit
  end
  filename = 'guests.csv'
  if ARGV.include?('--filename')
    idx = ARGV.index('--filename')
    filename = ARGV[idx+1] if idx
  end
  app.export(event_id, filename)

else
  puts "Unknown command. Use create, add-guest, list-guests, rsvp, stats, search, export."
end
