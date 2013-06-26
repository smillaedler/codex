![librarian](http://i.imgur.com/lvQmuIY.png)

## About

Librarian is **NOT** intended to be an ORM, but a relation algebra inspired by [Arel](http://www.github.com/rails/arel). Project is still under development.

## Installation

With Google's [Go](http://www.golang.org) installed on your machine:

    $ go get -u github.com/chuckpreslar/librarian

## Usage

To show a little sample of what using Libraian looks like, lets assume you have a table in your database (we'll call it `users`) and you want to select all the records it contains.  The SQL looks a little like this:

```sql
SELECT "users".* FROM "users"
```

Now using Librarian:

```go
import (
  /* Maybe some other imports. */
  l "github.com/chuckpreslar/librarian"
)


users := l.NewTable("users")
sql := users.ToSql()

```

Now that wasn't too bad, was it?

#### Projections

You can, of course, speed up your queries by only selecting the columns you need:

```go

// ...

users := l.NewTable("users")
sql := users.Project("id", "email", "first_name", "last_name").ToSql()

```

#### Filtering

An example of how to search for records that meet a specified criteria:

```go

// ...

users := l.NewTable("users")
sql := users.Where(users("id").Eq(1).Or(users("email").Eq("test@example.com"))).ToSql()

```

#### Joins

```go

// ...

users := l.NewTable("users")
orders := l.NewTable("orders")
sql := users.InnerJoin(orders.On(orders("id").Eq(users("order_id")))).ToSql()

// SELECT "users".* FROM "users" INNER JOIN "orders" ON "orders"."id" = "users"."order_id"

```

#### Updates

```go

// ...

users := l.NewTable("users")
sql := users.Set(users("admin").Eq(true)).Where(users("id").Eq(1)).ToSql()

// UPDATE "users" SET "users"."admin" = 'true' WHERE "users"."id" = 1

```

## License

> The MIT License (MIT)

> Copyright (c) 2013 Chuck Preslar

> Permission is hereby granted, free of charge, to any person obtaining a copy
> of this software and associated documentation files (the "Software"), to deal
> in the Software without restriction, including without limitation the rights
> to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
> copies of the Software, and to permit persons to whom the Software is
> furnished to do so, subject to the following conditions:

> The above copyright notice and this permission notice shall be included in
> all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
> THE SOFTWARE.
