#!/usr/bin/env ruby

require 'sysstat'
require 'influxdb'

servers = ['server1', 'server2', 'server3']
db_host = 'influxdb_server'
db_name = 'metrics_db'
threshold = 90

client = InfluxDB::Client.new db_name, host: db_host

while true
  servers.each do |server|
    stats = Sysstat.new(server)

    cpu_usage = stats.cpu.usage
    memory_usage = stats.memory.percent_used
    disk_space = stats.disk('/').percent_full

    data_point = {
      values: { cpu_usage: cpu_usage, memory_usage: memory_usage, disk_space: disk_space },
      timestamp: Time.now.to_i
    }

    client.write_point('system_metrics', data_point)

    if cpu_usage > threshold || memory_usage > threshold || disk_space > threshold
      send_alert(server, cpu_usage, memory_usage, disk_space)
    end
  end

  sleep(60) # Collect data every minute
end

def send_alert(server, cpu, memory, disk)
  # Send an alert message via your preferred notification method
  puts "ALERT: High resource utilization on #{server} - CPU: #{cpu}%, Memory: #{memory}%, Disk: #{disk}%"
end
