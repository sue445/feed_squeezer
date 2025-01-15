resource "google_cloud_run_service" "feed_squeezer" {
  name     = var.service_name
  location = var.location

  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/minScale" = var.min_instance_count
        "autoscaling.knative.dev/maxScale" = var.max_instance_count
      }
    }

    spec {
      service_account_name = google_service_account.feed_squeezer.email

      containers {
        # c.f. https://console.cloud.google.com/artifacts/docker/feed-squeezer/us/feed-squeezer/app
        image = "us-docker.pkg.dev/feed-squeezer/feed-squeezer/app:${var.tag}"

        resources {
          limits = {
            cpu    = "1",
            memory = "128Mi",
          }
        }

        ports {
          container_port = 8080
        }

        env {
          name  = "SENTRY_RELEASE"
          value = var.tag
        }

        dynamic "env" {
          for_each = length(var.sentry_dsn) > 0 ? [var.sentry_dsn] : []

          content {
            name  = "SENTRY_DSN"
            value = var.sentry_dsn
          }
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  autogenerate_revision_name = true
}

resource "google_cloud_run_service_iam_member" "feed_squeezer_allow_public_access" {
  for_each = toset([
    "roles/run.invoker",
  ])

  location = google_cloud_run_service.feed_squeezer.location
  project  = google_cloud_run_service.feed_squeezer.project
  service  = google_cloud_run_service.feed_squeezer.name
  role     = each.key
  member   = "allUsers"
}
