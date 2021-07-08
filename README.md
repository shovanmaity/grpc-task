## grpc-task

### Clone the repo
```bash
git clone --recurse-submodules --remote-submodules https://github.com/shovanmaity/grpc-task.git
cd grpc-task
```
### Run the app
```bash
make run
```
### Get profile without auth header
```bash
curl http://0.0.0.0:8090/api/v1/shovan/profile
```
```json
{"code":16,"message":"authorization token not found","details":[]}
```
### Update profile without auth header
```bash
curl -X PUT http://0.0.0.0:8090/api/v1/shovan/profile -d '{"name":"Shovan Maity", "email":"shovan.cse91+1@gmail.com"}'
```
```json
{"code":16,"message":"authorization token not found","details":[]}
```
### Registration
```bash
curl -X POST http://localhost:8090/api/v1/registration -d '{"username":"shovan", "password":"shovan"}'
```
```json
{"jwt":"JWT_TOKEN"}
```

### Login
```bash
curl -X POST http://localhost:8090/api/v1/login -d '{"username":"shovan", "password":"shovan"}'
```
```json
{"jwt":"JWT_TOKEN"}
```
### Get profile
```bash
curl http://0.0.0.0:8090/api/v1/shovan/profile -H "authorization: Token ${JWT_TOKEN}"
```
```json
{"username":"shovan", "name":"", "email":""}
```
### Update profile
```bash
curl -X PUT http://0.0.0.0:8090/api/v1/shovan/profile -H "Authorization: Token ${JWT_TOKEN}" -d '{"name":"Shovan Maity", "email":"shovan.cse91@gmail.com"}'
```
### Get profile after update
```bash
curl http://0.0.0.0:8090/api/v1/shovan/profile -H "Authorization: Token ${JWT_TOKEN}"
```
```json
{"username":"shovan", "name":"Shovan Maity", "email":"shovan.cse91@gmail.com"}
```
### Logout
```bash
curl -X DELETE http://localhost:8090/api/v1/logout -H "authorization: Token ${JWT_TOKEN}"
```
### Get profile after logout
```bash
curl http://0.0.0.0:8090/api/v1/shovan/profile -H "Authorization: Token ${JWT_TOKEN}"
```
```json
{"code":16, "message":"invalid session", "details":[]}
```
### Update profile after logout
```bash
curl -X PUT http://0.0.0.0:8090/api/v1/shovan/profile -H "Authorization: Token ${JWT_TOKEN}" -d '{"name":"Shovan Maity", "email":"shovan.cse91+1@gmail.com"}'
```
```json
{"code":16, "message":"invalid session", "details":[]}
```
