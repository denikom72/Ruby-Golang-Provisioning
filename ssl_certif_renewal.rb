#!/usr/bin/env ruby

require 'openssl'
require 'fileutils'

cert_path = '/etc/ssl/certs/your_domain.crt'
key_path = '/etc/ssl/private/your_domain.key'
renewal_script = '/path/to/renew_certificate.sh'

def renew_certificate(cert_path, key_path, renewal_script)
  if File.exist?(renewal_script)
    puts "Running certificate renewal script..."
    `#{renewal_script}`
    puts "Certificate renewed."
  else
    puts "No renewal script found. Attempting automatic renewal..."
    current_cert = OpenSSL::X509::Certificate.new(File.read(cert_path))
    expiration_date = current_cert.not_after

    if expiration_date < Time.now + (30 * 24 * 60 * 60) # Renew if certificate expires within 30 days
      puts "Renewing SSL certificate..."
      `certbot renew`
      FileUtils.cp('/etc/letsencrypt/live/your_domain/fullchain.pem', cert_path)
      FileUtils.cp('/etc/letsencrypt/live/your_domain/privkey.pem', key_path)
      puts "Certificate renewed."
    else
      puts "Certificate renewal not required at this time."
    end
  end
end

renew_certificate(cert_path, key_path, renewal_script)
