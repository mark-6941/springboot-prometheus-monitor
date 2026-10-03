Hcl

variable "vsphere_user" {
  type    = string
  default = "administrator@vsphere.local"
}

variable "vsphere_password" {
  type    = string
  default = "your-password"
}

variable "vsphere_server" {
  type    = string
  default = "192.168.1.100"
}