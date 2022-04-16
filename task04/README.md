# Docker compose

## Usage

start up the docker compose:

   ```shell
   docker-compose up
   ```

insert data to database:

   ```shell
   curl -X PUT -H "Content-Type: application/json" -d '{"Name":"bageta","Price":22,"Ammount":123}' http://localhost:9000/product
   ```

get data from database:

   ```shell
   curl http://localhost:9000/product

   curl http://localhost:9000/product/${ID}
   ```

API recoginses the following HTTP methods:

* GET
  * /product - get all products
  * /product/{id} - get product by id
* PUT
  * /product - insert product
* DELETE
  * /product/{id} - delete product by id
* PATCH
  * /product/{id} - update product by id

## Assignment

1. Bundle application from assignment 3 with database as docker compose file
   1. All images must be public
2. User volumes for data storage
3. Use separate docker network bridge

## Example way hot to build container containing Go binary

1. (Windows Only) [Cross compile to linux](https://stackoverflow.com/a/43945772)

   ```shell
   set GOARCH=amd64
   set GOOS=linux
   ```

2. Build

   ```shell
   go build -o myapp
   ```

3. Dockerfile

   ```Dockerfile
   FROM alpine:3.15.4
   COPY myapp myapp
   ENTRYPOINT ["./myapp"]
   ```

4. Build image and push

   ```shell
   docker build -t mydockeraccount/myapp:1.0.0 .
   docker push mydockeraccount/myapp:1.0.0
   ```
