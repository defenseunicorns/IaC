#---------------------------------------------------------------
# EKS Add-Ons
#---------------------------------------------------------------

module "eks_blueprints_kubernetes_addons" {
  source = "git::https://github.com/aws-ia/terraform-aws-eks-blueprints.git//modules/kubernetes-addons?ref=v4.24.0"

  depends_on = [module.eks_blueprints]

  eks_cluster_id           = module.eks_blueprints.eks_cluster_id
  eks_cluster_endpoint     = module.eks_blueprints.eks_cluster_endpoint
  eks_oidc_provider        = module.eks_blueprints.oidc_provider
  eks_cluster_version      = module.eks_blueprints.eks_cluster_version
  auto_scaling_group_names = module.eks_blueprints.self_managed_node_group_autoscaling_groups

  # EKS Managed Add-ons
  # VPC CNI - This needs to be done outside of the blueprints module
  # enable_amazon_eks_vpc_cni = var.enable_amazon_eks_vpc_cni
  # amazon_eks_vpc_cni_config = var.amazon_eks_vpc_cni_config

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
