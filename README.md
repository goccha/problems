# Problems
Problem Details for HTTP APIs [RFC7807](https://tools.ietf.org/html/rfc7807)

## Simple Usage
```go
problem := problems.New("/users").NotFound("user not found")
```

## Bad Request
```go
if err := validate.Struct(s); err != nil {
    problems.New(req.URL.Path, problems.NewBadRequest(err)).BadRequest("Invalid Parameters").JSON(ctx, req.Writer)
	return
}
```

##  Conversion to error
```go
err := problems.New("").Unauthorized("password mismatch").Wrap()
problems.ServerProblemOf(context.TODO, "/login", err).JSON(ctx, req.Writer)
```
