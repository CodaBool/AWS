variable "name" {
  type    = string
}

variable "key_name" {
  type    = string
  default = "win"
}

variable "profile_name" {
  type    = string
  default = "ec2_profile"
}

variable "instance_type" {
  type    = string
  default = "t4g.nano"
}

# variable "ami_name" {
#   type    = string
# }

variable "price" {
  type    = string
  default = ".0024"
}