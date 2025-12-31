package gogenerate

//go:generate go run ../../cmd/surreal-orm-gen -dir .

// orm:node table=users rename_all=snake_case permissions=full
type User struct {
	// orm:field name=id type=record<users>
	ID        string
	FirstName string
	LastName  string
	Email     string
}

// orm:edge table=user_account in=User out=Account
type UserAccount struct {
	Note string
}

// orm:node table=account rename_all=snake_case permissions=full
type Account struct {
	// orm:field name=id type=record<account>
	ID   string
	Name string
}
