# update
1. set latest go version in `go.mod`
2. run `go get -u .` from inside `src` folder
3. run `go mod tidy`

# environment
```json
{
  "action": "", # "manual", !=
  "secret": "", # match with what's in .env
  "test": "", # "true", !=
  "body": "" # any
}
```
