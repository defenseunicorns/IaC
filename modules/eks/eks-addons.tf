#---------------------------------------------------------------
# EKS Add-Ons
#---------------------------------------------------------------

module "eks_blueprints_kubernetes_addons" {
  source = "git::https://github.com/aws-ia/terraform-aws-eks-blueprints.git//modules/kubernetes-addons?ref=v4.27.0"

  eks_cluster_id           = module.aws_eks.cluster_name
  eks_cluster_endpoint     = module.aws_eks.cluster_endpoint
  eks_oidc_provider        = module.aws_eks.oidc_provider
  eks_cluster_version      = module.aws_eks.cluster_version
  auto_scaling_group_names = concat(lookup(module.aws_eks.self_managed_node_groups, "autoscaling_group_name", []), lookup(module.aws_eks.eks_managed_node_groups, "node_group_autoscaling_group_names", []))
  # EKS CoreDNS
  enable_amazon_eks_coredns = var.enable_amazon_eks_coredns
  amazon_eks_coredns_config = var.amazon_eks_coredns_config

  # EKS kube-proxy
  enable_amazon_eks_kube_proxy = var.enable_amazon_eks_kube_proxy
  amazon_eks_kube_proxy_config = var.amazon_eks_kube_proxy_config

  # EKS EBS CSI Driver
  enable_amazon_eks_aws_ebs_csi_driver = var.enable_amazon_eks_aws_ebs_csi_driver
  amazon_eks_aws_ebs_csi_driver_config = var.amazon_eks_aws_ebs_csi_driver_config


  # K8s Add-ons
  # EKS Metrics Server
  enable_metrics_server      = var.enable_metrics_server
  metrics_server_helm_config = var.metrics_server_helm_config

  # EKS AWS node termination handler
  enable_aws_node_termination_handler      = var.enable_aws_node_termination_handler
  aws_node_termination_handler_helm_config = var.aws_node_termination_handler_helm_config

  # EKS Cluster Autoscaler
  enable_cluster_autoscaler      = var.enable_cluster_autoscaler
  cluster_autoscaler_helm_config = var.cluster_autoscaler_helm_config
}
