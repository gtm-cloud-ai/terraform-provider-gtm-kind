terraform {
  required_providers {
    kind = {
      source = "gtm-cloud-ai/terraform-provider-gtm-kind"
      version = "~> 1.0"
    }
  }
}

# Provider will auto-detect available container runtime (podman preferred over docker)
provider "kind" {}

# Kind cluster using auto-detected container runtime
resource "kind_cluster" "auto_detect_example" {
  name           = "auto-detect-cluster"
  wait_for_ready = true

  kind_config = <<YAML
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
YAML
}

# Output the kubeconfig
output "kubeconfig" {
  value = kind_cluster.auto_detect_example.kubeconfig
  sensitive = true
}

# Output the cluster endpoint
output "endpoint" {
  value = kind_cluster.auto_detect_example.endpoint
}
