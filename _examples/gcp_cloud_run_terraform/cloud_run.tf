resource "google_cloud_run_service" "feed_squeezer" {
  name     = var.service_name
  location = var.location

  metadata {
    annotations = {
      "run.googleapis.com/invoker-iam-disabled" = true
    }
  }

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
