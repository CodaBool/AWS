variable "name" {
  type    = string
}

variable "key_name" {
  type    = string
  default = "win"
}

variable "ssh_ip" {
  type    = string
}

# variable "app_ports" {
#   type    = list(number)
# }

variable "instance_type" {
  type    = string
  default = "t4g.nano"
}