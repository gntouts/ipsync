package main

import (
	"time"
)

const TIMEOUT = 20 // seconds

func main() {
	log_info("Started monitoring IP address", "main")

	target := config().get_dns_target()

	netlify := NewNetlifyClient()
	ip := NewIpClient()

	zone := netlify.get_dns_zone(target)
	record_id, record_ip := netlify.get_dns_record(zone, target)
	log_info("Netlify IP is set "+record_ip, "main")

	for {
		current := ip.get_ip()

		if current.Ip != record_ip {
			netlify.delete_dns_record(zone, record_id)
			changed := netlify.create_dns_record(zone, target, current.Ip)
			log_info("Local IP changed to "+current.Ip, "main")
			var msg string
			if changed {
				msg = "Updated IP address to " + current.Ip
			} else {
				msg = "Failed to update IP address to " + current.Ip
			}
			log_info(msg, "main")
			record_id, record_ip = netlify.get_dns_record(zone, target)
		}
		duration := time.Duration(config().get_timeout()) * time.Second
		time.Sleep(duration)
	}
}
