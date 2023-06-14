variable "name" {
  default = "scraper"
  type    = string
}

variable "tag" {
  description = "Tag to use for deployed Docker image"
  default     = "latest"
}