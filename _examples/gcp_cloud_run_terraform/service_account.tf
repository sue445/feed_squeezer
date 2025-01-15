resource "google_service_account" "feed_squeezer" {
  account_id = var.service_account_name
}
