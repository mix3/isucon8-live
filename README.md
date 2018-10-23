```
$ export NAME={portal_login_name}
$ export PASSWORD={portal_login_password}
$ go run main.go | tee -a log | xargs -I@ say @
```
