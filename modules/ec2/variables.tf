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

variable "ssh_ip" {
  type    = string
}

variable "instance_type" {
  type    = string
  default = "t4g.nano"
}

variable "price" {
  type    = string
  default = ".0024"
}