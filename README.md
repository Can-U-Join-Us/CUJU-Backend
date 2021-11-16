# CUJU-Backend


- router       4000
- Docker db    3333 : 3306
- Docker redis 6379 : 6379
`go run main.go`

--------
DB Schema table

- User ( pk : UID )
  - Join ( pk : UID - fk )

- Post ( pk : PID )
  - Image ( pk : PID - fk )
  - Member ( pk : PID - fk )


