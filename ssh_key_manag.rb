#!/usr/bin/env ruby

require 'net/ssh'
require 'optparse'

options = {}
OptionParser.new do |opts|
  opts.banner = "Usage: ssh_key_management.rb [options]"
  opts.on("-a", "--add KEY", "Add an SSH key") { |key| options[:add] = key }
  opts.on("-r", "--remove KEY", "Remove an SSH key") { |key| options[:remove] = key }
  opts.on("-h", "--host HOST", "Target host") { |host| options[:host] = host }
end.parse!

host = options[:host]
key_to_add = options[:add]
key_to_remove = options[:remove]

Net::SSH.start(host, 'admin', password: 'your_password') do |ssh|
  if key_to_add
    ssh.exec!("echo #{key_to_add} >> ~/.ssh/authorized_keys")
  end

  if key_to_remove
    ssh.exec!("sed -i '/#{key_to_remove}/d' ~/.ssh/authorized_keys")
  end
end
