# WVC-Sync
service which is intended to run as cronjob and sync all users bidirectional between easyverein and wordpress

# Usage
## locally
- Before running the service you need to create a config.yaml  
- The content of this yaml should look exactly like the example-config.yaml  
- Replace "EASYVEREIN-API-TOKEN" with your easyverein token  
- Either place the config.yaml to ./configs/config.yaml or specify the path via -config flag

Run `go run cmd/service/main.go (-config path/to/config.yaml)` to execute the service

# Curlz used for the service API
### Easyverein
```sh
curl https://easyverein.com/api/stable/contact-details\?limit\=100&page=1 -H "Authorization: Token <TOKEN>"
```

# TODOs:
- [x] Add config tests
- [ ] Fetch users via wordpress api
- [ ] Implement sync(compare) algorithm
- [ ] Add the new users (result of sync algorithm) to wordpress via API
