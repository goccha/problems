# Problems
Problem Details for HTTP APIs [RFC7807](https://tools.ietf.org/html/rfc7807)

## Simple Usage
```go
problem := problems.New("/users").NotFound("user not found")
```

## Bad Request
```go
if err := validate.Struct(s); err != nil {
    problem := problems.ClientProblemOf(context.TODO, "/validate", err)	
}
```

##  Conversion to error
```go
err := problems.New("").Unauthorized("password mismatch").Wrap()
problem := problems.ServerProblemOf(context.TODO, "/login", err)
```

## Json Response
```go
func(w http.ResponseWriter, req *http.Request) {
    ctx := req.Context()
    problems.New("/users", problems.NewBadRequest(err)).BadRequest("user not found").JSON(ctx, w)
}
```

