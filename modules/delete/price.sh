#!/bin/bash

# get lowest price spot price for t4g.nano in us-east-1a
PRICE=$(aws --region=us-east-1 ec2 describe-spot-price-history --instance-types t4g.nano --start-time=$(date +%s) --product-descriptions="Linux/UNIX" --query 'SpotPriceHistory[*].{az:AvailabilityZone, price:SpotPrice}' | jq -r ".[] | select(.az == \"us-east-1a\") | .price")
jq -n --arg price "$PRICE" '{"price":$price}'

# DAY / MONTH

#   T3
#     NANO .5Gb
# .1128 / 3.384

#   T4
#     NANO .5Gb
# .1008 / 3.024 

#     MICRO 1Gb
# .1824 / 5.472

#     SMALL 2Gb
# .3192 / 9.576

#     MEDIUM 4GB
# .5064 / 15.192


