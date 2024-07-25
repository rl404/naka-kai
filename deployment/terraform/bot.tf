resource "kubernetes_deployment" "bot" {
  metadata {
    name = var.gke_deployment_name
    labels = {
      app = var.gke_deployment_name
    }
  }

  spec {
    replicas = 1
    selector {
      match_labels = {
        app = var.gke_deployment_name
      }
    }
    template {
      metadata {
        labels = {
          app = var.gke_deployment_name
        }
      }
      spec {
        container {
          name    = var.gke_deployment_name
          image   = var.gcr_image_name
          command = ["./naka-kai"]
          args    = ["bot"]
          env {
            name  = "NAKA_KAI_DISCORD_TOKEN"
            value = var.naka_kai_discord_token
          }
          env {
            name  = "NAKA_KAI_DISCORD_DELETE_TIME"
            value = var.naka_kai_discord_delete_time
          }
          env {
            name  = "NAKA_KAI_DISCORD_QUEUE_LIMIT"
            value = var.naka_kai_discord_queue_limit
          }
          env {
            name  = "NAKA_KAI_DB_DIALECT"
            value = var.naka_kai_db_dialect
          }
          env {
            name  = "NAKA_KAI_DB_ADDRESS"
            value = var.naka_kai_db_address
          }
          env {
            name  = "NAKA_KAI_DB_NAME"
            value = var.naka_kai_db_name
          }
          env {
            name  = "NAKA_KAI_DB_USER"
            value = var.naka_kai_db_user
          }
          env {
            name  = "NAKA_KAI_DB_PASSWORD"
            value = var.naka_kai_db_password
          }
          env {
            name  = "NAKA_KAI_YOUTUBE_KEY"
            value = var.naka_kai_youtube_key
          }
          env {
            name  = "NAKA_KAI_LOG_LEVEL"
            value = var.naka_kai_log_level
          }
          env {
            name  = "NAKA_KAI_LOG_JSON"
            value = var.naka_kai_log_json
          }
          env {
            name  = "NAKA_KAI_NEWRELIC_LICENSE_KEY"
            value = var.naka_kai_newrelic_license_key
          }
        }
      }
    }
  }
}
