# TODO - go-api

## Folder/File structure
### Main
- `cmd/main.go`: start + bootstrap project
### /internal/app = app logic, other is vendors
- `handler` - route
- `service` - business logic
- `config`

### MongoDB
- `model` => type stuct : how json look like { username, password }
- `database` => type interface : what function we have CreateUser, ListUser, GetUser etc.
- `dbmongo` => client mongo
### redis
- `cache` - interface / repo
- `cacheredis` => client redis

### Flow
`handler` => `middleware` => `service` => `repo` => `database`

//>> Composite มันมีอะไร << go ไม่มี inherit ... inheritanace มันเปนอะไร
//>> Concrete <<

pkg/common => common packages, util, หรือเป็น package กลาง ให้ microservice อื่นๆ ใช้
vendor => dependencies

go mod pkg
go mod vender

ตั้ง mongoDB

// tour-go, 

pagekage เดียวกัน ชื่อ function ไม่ควรซ้ำกัน