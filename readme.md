# Warehouse controller

## Installation and running the application
* make sure kubernetes, helm, helm-s3, and golang installed on your machine.
* make sure you have access to `gudangada-bi-helm-charts`
* run `go get ./...`
* run `go run cmd/main.go`

## Using warehouse controller

### Deploying single helm chart

#### Release
POST `/chart`
```
{
    "release_name": string,
    "name": string,
    "version": string(optional),
    "values": JSON(optional)
}
```
Will return `HTTP 200` if success and `HTTP 400` if failed.
#### Get All Release
GET `/chart`  
Will return `HTTP 200` alongside with response body if success and `HTTP 400` if failed.  
response body:
```
["release-name-1", "release-name-2", ... , "release-name-N"]
```
#### Get Release Detail
GET `/chart/{release-name}`
Will return `HTTP 200` alongside with response body if success and `HTTP 400` if failed.  
response body:
```
{
    "release_name": string,
    "name": string,
    "version": string(optional),
    "values": JSON(optional)
}
```
#### Update release
PUT `/chart/{release-name}`
```
{
    "release_name": string,
    "name": string,
    "version": string(optional),
    "values": JSON(optional)
}
```
Will return `HTTP 200` if success and `HTTP 400` if failed.
#### Delete Release
DELETE `/chart/{release-name}`  
Will return `HTTP 200` if success and `HTTP 400` if failed.

### Using module

#### Add Module
POST `/module`
```
{
    "name": string,
    "version": string,
    "spec": {
        "chart": []chartspec
    }
}
```
Will return `HTTP 200` if success and `HTTP 400` if failed.
#### Add Module Release
POST `/module/release`
```
{
    "release": string,
    "module": string,
    "version": string,
    "values": {
        (key pair value needed as specified by module spec)
    }
}
```
Will return `HTTP 200` if success and `HTTP 400` if failed.
#### Update Module Release
PUT `/module/release/{release-name}`
```
{
    "release": string,
    "module": string,
    "version": string,
    "values": {
        (key pair value needed as specified by module spec)
    }
}
```
Will return `HTTP 200` if success and `HTTP 400` if failed.
#### Delete Module Release
DELETE `/module/release/{release-name}`  
Will return `HTTP 200` if success and `HTTP 400` if failed.