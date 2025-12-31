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

// orm:node table=post schemaless=true permissions="FOR select WHERE published = true OR user = $auth.id FOR create, update WHERE user = $auth.id FOR delete WHERE user = $auth.id OR $auth.admin = true"
type Post struct {
	// orm:field name=id type=record<post>
	ID        string
	Title     string
	Body      string
	Published bool
	User      string
}
