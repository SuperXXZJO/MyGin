# MyGin

## example

``

```go
func main() {

   r:=MyGin.NewEngine()

   r.GET("/", func(context *MyGin.Context) {
      context.String(200,"hello world")
   })
   r.Run(":8080")


}
```

