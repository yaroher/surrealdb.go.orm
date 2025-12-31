package protobuf

//go:generate protoc -I ../../ -I . --go_out=. --go_opt=paths=source_relative --surreal-orm_out=. --surreal-orm_opt=paths=source_relative models.proto
