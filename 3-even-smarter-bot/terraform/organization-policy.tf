# resource "google_org_policy_policy" "allowed_policy_member_domains" {
#   name   = "${data.google_organization.main.name}/policies/iam.allowedPolicyMemberDomains"
#   parent = data.google_organization.main.name
#   spec {
#     rules {
#       values {
#         allowed_values = [
#           "allUsers",
#           "C01vrxtc4", #craigdickson.altostrat.com
#           "C02h8e9nw", #google.com
#         ]
#       }
#     }
#   }

#   depends_on = [
#     google_project_service.main // need to wait for the org policy api to be enabled on the project
#   ]
# }
