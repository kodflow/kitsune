terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
    }
  }
}

# Définition du fournisseur Docker
provider "docker" {
  host = "unix:///private/var/run/docker.sock"
}

locals {
  current_directory = abspath("../../")
  services_folder   = abspath("${local.current_directory}/services")
  services_dirs     = distinct(flatten([for _, v in flatten(fileset(path.module, "../../services/**")) : regex("../../services/([^/]+)", dirname(v))]))
}

resource "docker_image" "kitsune" {
  count = length(local.services_dirs)
  name  = "kodmain/kitsune.${local.services_dirs[count.index]}"

  build {
	platform   = "linux/arm64,linux/amd64,windows/amd64,darwin/arm64"
    context    = "${local.current_directory}"
    dockerfile = "deploy/publish/Dockerfile"
    build_args  = { 
		BINARY_NAME = "${local.services_dirs[count.index]}"
	}
    tag        = [
		"kodmain/kitsune.${local.services_dirs[count.index]}:latest",
		"kodmain/kitsune.${local.services_dirs[count.index]}:${formatdate("YYMMDD", timestamp())}",
	]
  }
}

/*
resource "null_resource" "push_to_dockerhub" {
  count = length(local.services_dirs)

  triggers = {
    image_id = docker_image.kitsune[count.index].id
  }

  provisioner "local-exec" {
    command = "docker push kodmain/kitsune.${local.services_dirs[count.index]}:latest"
  }
}
*/

