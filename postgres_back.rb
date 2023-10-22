#!/usr/bin/env ruby

require 'date'

backup_dir = '/var/backup'
db_name = 'your_db'
username = 'db_user'
backup_file = "#{backup_dir}/db_backup_#{Date.today.strftime('%Y%m%d')}.sql"

# Backup the database
puts "Backing up database..."
`pg_dump -U #{username} -d #{db_name} -f #{backup_file}`

# Restore the database
puts "Restoring database..."
`psql -U #{username} -d #{db_name} -f #{backup_file}`

puts "Backup and restore completed."
