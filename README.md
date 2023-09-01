# Step 1

Here you will find 2 folders `order-service` and `user-service`

They are our besic golang projects. We will use them as component of our "micro service".

> The example code is not production ready. 

## Overview

![step-1-overview](/images/step-1-overview.png)

## Start the server

### start `user-service`
```shell
cd user-service 
go run main.go
```

The user service server will run on port :5011

try request to the server
```shell
curl localhost:5011/users
```

> you can use Postman too if you like.

The server should response with json like this
```shell
{"users":[{"id":"user-a","name":"A"},{"id":"user-b","name":"B"},{"id":"user-c","name":"C"}]}
```

try call another endpoint on the user service
get `/user/user-a`
```shell
curl localhost:5011/user/user-a
```

You will get error msg
```shell
{"msg":"Get \"http://localhost:5010/order/user-a\": dial tcp [::1]:5010: connect: connection refused"}
```

Because when call `user-service/user/:user-id`, it will request that user orders from the `order-service`.
If we didn't start `order-service` server, it will error on `user-service` side.

### Let's start `order-service`

open new terminal
then run
```shell
cd order-service
go run main.go
```

The `order-service` will run on port `:5010`

Now try get user-a from the `user-service` again
```shell
curl localhost:5011/user/user-a
```

the `user-service` will response with
```shell
{"user":{"id":"user-a","name":"A","orders":[{"ID":"order-A-1","UserID":"user-a","Items":[{"Name":"itemA","Price":10,"TotalPrice":20,"Qty":2},{"Name":"itemB","Price":20,"TotalPrice":20,"Qty":1}]}]}}
```

You can also to request to order-service too 
```shell
curl localhost:5010/order/user-a
```

## Next step

Now we can run our service locally.

Let's go to [step 2](https://gitlab.com/thanabutjlmwn/workshop/-/tree/main/step-2-containerlize-application?ref_type=heads)
