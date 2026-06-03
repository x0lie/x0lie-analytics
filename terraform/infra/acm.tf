resource "aws_acm_certificate" "analytics" {
  domain_name       = "analytics.x0lie.com"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "cloudflare_dns_record" "acm_validation" {
  for_each = {
    for dvo in aws_acm_certificate.analytics.domain_validation_options : dvo.domain_name => dvo
  }

  zone_id = var.cloudflare_zone_id
  name    = each.value.resource_record_name
  type    = each.value.resource_record_type
  content = each.value.resource_record_value
  ttl     = 60
}

resource "aws_acm_certificate_validation" "analytics" {
  certificate_arn         = aws_acm_certificate.analytics.arn
  validation_record_fqdns = [for dvo in aws_acm_certificate.analytics.domain_validation_options : dvo.resource_record_name]
  depends_on              = [cloudflare_dns_record.acm_validation]
}
