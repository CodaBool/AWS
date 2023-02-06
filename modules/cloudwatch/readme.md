# variables
- {instance_id} = i-00f6856ce0836e175
- {hostname} = ip-172-31-95-166.ec2.internal
- {local_hostname} = ip-172-31-95-166.ec2.internal

# cmd
- sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl -a fetch-config -m ec2 -s -c file:/opt/aws/agent.json