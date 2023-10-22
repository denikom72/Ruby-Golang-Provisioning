#!/usr/bin/env ruby

require 'date'

logs_dir = '/var/log/your_app_logs'
report_file = 'log_report.txt'

# Define criteria and patterns for log entry analysis
error_pattern = /ERROR/
date_range = (Date.today - 5)..Date.today

report = File.open(report_file, 'w')

Dir.glob("#{logs_dir}/*.log").each do |log_file|
  File.readlines(log_file).each do |line|
    if line.match?(error_pattern)
      log_date = Date.parse(line.match(/\d{4}-\d{2}-\d{2}/).to_s)
      if date_range.include?(log_date)
        report.puts(line)
      end
    end
  end
end

report.close

# Generate a summary report
error_count = `grep -c 'ERROR' #{report_file}`.to_i
puts "Log analysis completed. Found #{error_count} errors in the last 5 days."
