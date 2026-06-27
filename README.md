# problems
~~Problem Details for HTTP APIs [RFC7807](https://tools.ietf.org/html/rfc7807)~~

golang library for handling HTTP API errors defined in RFC9457([HTTP JSON Problem Details](https://tools.ietf.org/html/rfc9457))

RFC9457([HTTP JSON Problem Details](https://tools.ietf.org/html/rfc9457))で定義されたHTTP APIのエラーを扱うためのGO言語のライブラリーです。

Can handle http api error details as golang error

http apiのエラー詳細をgolangのerrorとして処理できます。

## Simple Usage
```go
problem := problems.New(problems.Instance("/users")).NotFound("user not found")
```

## Bad Request (RFC7807)
```go
if err := validate.Struct(s); err != nil {
    problems.New(problems.Path(req), problems.InvalidParams(err)).BadRequest("Invalid Parameters").JSON(ctx, req.Writer)
	return
}
```

## Bad Request (RFC9457)
```go
if err := validate.Struct(s); err != nil {
    problems.New(problems.Path(req), problems.ValidationErrors(err)).BadRequest("Invalid Parameters").JSON(ctx, req.Writer)
	return
}
```

##  Conversion to error
```go
err := problems.New().Unauthorized("password mismatch").Wrap()
problems.Of(context.TODO, "/login", err).JSON(ctx, req.Writer)
```
