# MVP Match test task

Based on Swag example: https://github.com/swaggo/swag/tree/master/example/celler

Assumptions and differences from the Exercise brief: 
1. Additional role "admin" was introduced to assign the permissions for not user-specific CRUD operations, such as Get Users.
2. Based on field names definition it is assumed that each user have only one role
3. In Bonus section it was not clear if logging in when being logged in already should produce an error. Logically existence of /logout/all API assumes that it should be so.
4. The in-memory data repository without persistence was implemented for simplicity.
5. As the parameters for logout/all were not specified, I assumed we should authenticate the user, and as there is no known session exist at that moment, the only mean to do it is passing the username/password, similar to Login

Generate doc

```console
$ go get -u github.com/swaggo/swag/cmd/swag
$ swag init
```

Run app

```console
$ go run main.go
```

Run tests

```console
$ go test ./test/... -v -count=1
```

[open swagger](http://localhost:8081/swagger/index.html)

