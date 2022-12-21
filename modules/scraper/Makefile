build:
  # Make sure you ran awsume for your account
	aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 919759177803.dkr.ecr.us-east-1.amazonaws.com
	docker build -t scraper:latest .
run:
  # Make sure you ran awsume for your account
	docker run \
		-e AWS_ACCESS_KEY_ID=$$AWS_ACCESS_KEY_ID \
		-e AWS_SECRET_ACCESS_KEY=$$AWS_SECRET_ACCESS_KEY \
		-e AWS_SESSION_TOKEN=$$AWS_SESSION_TOKEN \
		--env-file ./.env \
		-p 9000:8080 \
		--name scraper \
		scraper:latest
test:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'
clean:
	docker rm scraper