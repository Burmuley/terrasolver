dependency "iam_roles" {
  config_path = "../../global/iam-roles"
}

dependency "load-balancers" {
  config_path = "../load-balancers"
}

dependency "target-groups" {
  config_path = "../target-groups"
}

dependency "ecs-cluster" {
  config_path = "../ecs-clusters"
}
