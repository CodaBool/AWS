use `aws --region=us-east-1 ec2 describe-spot-price-history --instance-types t4g.nano --start-time=$(date +%s) --product-descriptions="Linux/UNIX" --query 'SpotPriceHistory[*].{az:AvailabilityZone, price:SpotPrice}'` to find best price spot instances. Make sure to change instance type to match what you want to search


#### Find latest al2 ami
- `aws ec2 describe-images --owners amazon --filters "Name=name,Values=amzn*" "Name=architecture,Values=arm64" --query 'sort_by(Images, &CreationDate)[-1].ImageId' --output text`


# TODO
- look into new metadata
- look into better monitoring setting