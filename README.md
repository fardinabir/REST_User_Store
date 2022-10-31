# REST_User_Store
This project is for implementing API based CRUD operations. It covers the basic concepts of REST API, Gorrila MUX and GORM. 

The main requests and responses sample are as like below:

### POST /users ###
```
Request:
{
    "firstName": <firstname>,
    "lastName": <lastname>,
    "password": <password>,
    "phone": <phone>
}
Response:
{
    "id": <id>
}
```
### GET /users/{id} ###
```
Response:
{
    "id": <id>,
    "name": <full name [firstname and lastname]>,
    "phone": <phone>
}
```

### POST /users/{id}/tags ###
```
Request:
{
    "tags": [<tag 1>, <tag 2>, ...],
    "expiry": <miliseconds>        
}
Response:
{}
```

### GET /users?tags=tag1,tag2... ###
```
Response:
{
    "users": [
        {
            "id": <id 1>
            "name": "<full name 1>"
            "tags": [<tag 1>, <tag 2>, ...]
        },
    ]
}
```
