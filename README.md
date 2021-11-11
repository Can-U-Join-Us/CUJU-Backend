# CUJU-Backend


- port 4000
- db   3333 : 3306

`go run main.go`

--------
DB Schema table

- User ( pk : UID )
  - Join ( pk : UID - fk )

- Post ( pk : PID )
  - Image ( pk : PID - fk )
  - Member ( pk : PID - fk )


