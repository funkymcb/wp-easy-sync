# WP-Easy-Sync
Service which is intended to run as cronjob to sync all users bidirectional between easyverein and wordpress

# Prerequisites
- Install WP-Plugin 'miniOrange API Authentication' and enable Basic Auth

# Usage
## locally
- Before running the service you need to create a config.yaml  
- The content of this yaml should look exactly like the example-config.yaml  
- Replace "EASYVEREIN-API-TOKEN" with your easyverein token  
- Either place the config.yaml to ./configs/config.yaml or specify the path via -config flag

Run the following command to execute the service
```sh
go run cmd/service/main.go (-config path/to/config.yaml)
```

## as docker container

# Curlz
### Easyverein GET users
```sh
curl https://easyverein.com/api/stable/contact-details\?limit\=100&page=1 -H "Authorization: Token <TOKEN>"
```
### Wordpress GET users
```sh
curl https://your-wp-domain.com/wp-json/wp/v2/users\?per_page\=100\&page\=1\&context\=edit -H "Authorization:Basic base64encoded(username:password)"
```

# Docs
[Easyverein API Documentation](https://easyverein.com/api/documentation/)  
[Wordpress API Documentation](https://developer.wordpress.org/rest-api/)  
[Wordpress API List Users](https://developer.wordpress.org/rest-api/reference/users/#list-users)

# TODOs:
- [x] Add config tests
- [x] Fetch users via wordpress api
- [x] Implement sync(compare) algorithm
- [ ] Add the new users (result of sync algorithm) to wordpress via API
- [ ] Optimize sync algorithm
- [ ] Add verbose debug output (by using --debug or -v -vv -vvv flags)
