# Author: Denis Komnenovic

require 'json'
require 'net/http'
require 'uri'

class VoIPProvisioner
  def initialize(username, password, endpoint)
    @username = username
    @password = password
    @endpoint = endpoint
  end

  def provision_account(account)
    account_data = {
      'id' => account[:id],
      'phone_number' => account[:phone_number],
      'username' => @username,
      'password' => @password,
      'settings' => account[:settings]
    }

    response = send_provision_request(account_data)

    if response['success'] == true
      puts 'Account provisioned successfully.'
    else
      puts "Account provisioning failed. Error: #{response['error']}"
    end
  end

  def send_provision_request(data)
    endpoint_url = "#{@endpoint}/provision"
    headers = { 'Content-Type' => 'application/json' }
    uri = URI.parse(endpoint_url)

    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = (uri.scheme == 'https')

    request = Net::HTTP::Post.new(uri.request_uri, headers)
    request.body = data.to_json

    response = http.request(request)

    JSON.parse(response.body)
  end
end

class VoIPAccount
  attr_accessor :id, :phone_number, :settings

  def initialize(id, phone_number, settings)
    @id = id
    @phone_number = phone_number
    @settings = settings
  end
end

class VoIPSettings
  attr_accessor :codec, :quality, :call_forwarding, :voicemail

  def initialize(codec, quality, call_forwarding, voicemail)
    @codec = codec
    @quality = quality
    @call_forwarding = call_forwarding
    @voicemail = voicemail
  end
end

username = 'your_username'
password = 'your_password'
endpoint = 'https://voip-provider-api.com'

provisioner = VoIPProvisioner.new(username, password, endpoint)

account_id = '12345'
phone_number = '+1234567890'

settings = VoIPSettings.new('G.729', 'Standard', true, false)
account = VoIPAccount.new(account_id, phone_number, settings)

begin
  provisioner.provision_account(account)
rescue StandardError => e
  puts "Error provisioning account: #{e}"
end

