resource "cloudflare_dns_record" "analytics" {
  zone_id = var.cloudflare_zone_id
  name    = "analytics"
  type    = "CNAME"
  content = aws_lb.main.dns_name
  ttl     = 1
  proxied = false
}
