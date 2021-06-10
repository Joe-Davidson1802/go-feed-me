resource "google_cloud_run_service" "run_service" {
    for_each = toset(local.services)

    name = "gfm-${each.value}-run"
    location = "europe-west2"
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
    for_each = toset(local.services)
    location    = google_cloud_run_service.run_service[each.value].location
    project     = google_cloud_run_service.run_service[each.value].project
    service     = google_cloud_run_service.run_service[each.value].name

    policy_data = data.google_iam_policy.noauth.policy_data
}