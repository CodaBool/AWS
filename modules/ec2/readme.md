use `aws ec2 describe-spot-price-history --instance-types t4g.nano --filters Name=spot-price,Values=0.001400 | jq -r '.SpotPriceHistory[] | "\(.SpotPrice) @ \(.AvailabilityZone)"'` to find best price spot instances


or this one, but it seems wrong

- `aws --region=us-east-1 ec2 describe-spot-price-history --instance-types t4g.micro --start-time=$(date +%s) --product-descriptions="Linux/UNIX" --query 'SpotPriceHistory[*].{az:AvailabilityZone, price:SpotPrice}'`


#### Find latest al2 ami
- `aws ec2 describe-images --owners amazon --filters "Name=name,Values=amzn*" "Name=architecture,Values=arm64" --query 'sort_by(Images, &CreationDate)[-1].ImageId' --output text`