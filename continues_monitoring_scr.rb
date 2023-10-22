#!/usr/bin/env ruby

require 'net/ping'
require 'net/smtp'

servers = ['server1', 'server2', 'server3']
services = ['http://server1/service', 'http://server2/service']

def send_alert(message)
  from = 'alerts@example.com'
  to = 'admin@example.com'

  msg = <<MESSAGE_END
  From: #{from}
  To: #{to}
  Subject: Server/Service Alert

  #{message}
MESSAGE_END

  Net::SMTP.start('smtp.example.com') do |smtp|
    smtp.send_message msg, from, to
  end
end

while true
  servers.each do |server|
    unless Net::Ping::External.new(server).ping?
      send_alert("#{server} is down!")
    end
  end

  services.each do |service|
    response = Net::HTTP.get_response(URI.parse(service))
    unless response.code == '200'
      send_alert("#{service} is not responding correctly! Response code: #{response.code}")
    end
  end

  sleep(60) # Check every minute
end
