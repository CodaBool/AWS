# EC2
> check back in [cost explorer](https://us-east-1.console.aws.amazon.com/cost-management/home?region=us-east-1#/cost-explorer) in February

## commands
- `sudo yum update`
- `sudo amazon-linux-extras install -y nginx1 && sudo systemctl enable nginx.service && sudo systemctl start nginx.service` start nginx
- `aws --region=us-east-1 ec2 describe-spot-price-history --instance-types t4g.micro --start-time=$(date +%s) --product-descriptions="Linux/UNIX" --query 'SpotPriceHistory[*].{az:AvailabilityZone, price:SpotPrice}'`
- `aws ec2 describe-images --owners amazon --filters "Name=name,Values=al2*" "Name=architecture,Values=arm64" --query 'sort_by(Images, &CreationDate)[-1].ImageId' --output text`

### Memory
- `top` use shift + m to sort by memory
- `free -m`
- `cat /proc/meminfo`
- `sudo slabtop`

> spicy
- `echo 2 > /proc/sys/vm/drop_caches` requires being sudo
  - this freed 85MB / 20% of memory, could be cron every 12 hours if you want to make Linus unhappy. I noticed after an hour it went back to its regular value of 126M / 34% usage
- `echo vm.vfs_cache_pressure=1000 >> /etc/sysctl.conf` requires being sudo
- `sudo sysctl -w vm.swappiness=10` 
- `sudo sysctl -w vm.vfs_cache_pressure=200`

## Notes
- AL2 Arm t4g.nano has 169MB free, 79MB used / 435MB. (0MB used in swap, if this is high then it needs to be expanded)
- in theory socket.io can handle ~1,000 (200MB) concurrent connections on this machine
- redis should only be used on machines with at least 4GB RAM, t4g medium would cost $7

### Costs
- .5GB = $1.01 / month
- 1GB  = $1.80
- 2GB  = $4.03
- 4GB  = $7.27