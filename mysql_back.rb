#!/usr/bin/env ruby

require 'date'

database_name = 'your_db'
username = 'db_user'
password = 'db_password'
backup_directory = '/path/to/backup'
remote_server = 'remote_server'
remote_directory = '/backup_storage'

timestamp = DateTime.now.strftime('%Y%m%d%H%M%S')
backup_file = "#{backup_directory}/backup_#{timestamp}.sql.gz"

`mysqldump -u #{username} -p#{password} #{database_name} | gzip > #{backup_file}`

`scp #{backup_file} #{remote_server}:#{remote_directory}`

# Optionally, you can clean up old backups to save disk space.
`find #{backup_directory} -name 'backup_*' -mtime +7 -exec rm {} \\;`
