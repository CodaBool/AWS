# Emailer
### Commands
> assumes you have a Docker daemon and have ran `awsume`
- `make build` builds the docker container
- `make run` starts a simulated lambda on port 9000
- `make test` sends a request to the simulated lambda which will potential send emails
- `make clean` removes the container, necessary before building again

<details>
<summary style="font-size: 1.4rem;">Details</summary>

### Build
> make sure you are in the same directory as the Dockerfile

`docker build -t emailer:latest .`

### Run
> make sure you have assumed role recently
```bash
docker run \
  -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
  -e AWS_SESSION_TOKEN=$AWS_SESSION_TOKEN \
  -p 9000:8080 \
  --name emailer \
  emailer:latest emailer
```

### Test
> In a different terminal run a curl against the running container. You can edit the data json to test different input

`curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'`

### Remove
`docker rm emailer`

### Alternative
Building docker containers between changes can be time consuming. For this reason the JavaScript file can be ran locally without Docker by running `npm start`


</details>