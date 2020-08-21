# Go GraphQL Server

A basic GraphQL server based on following Schema.

```GraphQL
author {
    Id
    Name
    ISBN No
    Books [book]
}

book {
    Id
    Title
    Description
    Authors [author]
}
```

#### Table Of Contents

- [Server Setup](#server-setup)
- [Mutations](#mutation)
- [Query](#query)
- [Solving N+1 Problem](#solving-n1-problem)

## Server Setup

Database `ArangoDB`. Spin up database

```DockerFile
docker-compose up -d
```

Get project dependency

```golang
go get -v ./...
```

Run

```golang
go run main.go
```

## Mutation

### Create Author

```cURL
curl -g 'http://localhost:8080/graphql?query=mutation+_{author(name:"pavel",isbn_no:"19923"){id,name,isbn_no}}'
```

or in GraphiQL

```GraphQL
mutation {
  author(name:"pavel", isbn_no:"19923"){
    id
    name
    isbn_no
  }
}
```

### Create Book

```
curl -g 'http://localhost:8080/graphql?query=mutation+_{book(title:"A+Book+of+Fire",description:"A+World+Famous+Book",authors:["154"]){id,title,description}}'
```

or in GraphiQL

```GraphQL
mutation {
  book(title:"A Book of Fire", description:"A world famous book", authors:["154"]) {
    id
    title
    description
  }
}
```

## Query

### Get Books of Author

```
curl -g 'http://localhost:8080/graphql?query={book{id,title,description,author(name:"tarique"){id,name,isbn_no}}}'
```

or in GraphQL

```GraphQL
query GetBooksOfAuthor {
  book {
    id
    title
    description
    author(name: "John") {
      id
      name
    }
  }
}
```

### Get All Books

```
curl -g 'http://localhost:8080/graphql?query={books{id,title,description}}'
```

or in GraphQL

```GraphQL
query {
  books {
    id
    title
    description
  }
}
```

### Get All Authors

```
curl -g 'http://localhost:8080/graphql?query={authors{id,name,isbn_no}}'
```

or in GraphQL

```GraphQL
query {
  authors {
    id
    name
    isbn_no
  }
}
```

Supports nested query

```GraphQL
query {
  authors {
    id
    name
    isbn_no
    books {
        id
        title
        description
    }
  }
}
```

## Solving N+1 Problem

The N+1 problem means that server executes multiple nnecessary round trips to datastores for nested data. Lets say we are querying an author and all of his books. If this particular author has `N` books the server will hit datastore 1 time for author and N times to retrieve N books. Hence it is called N+1 problem.

In our simple server this particular problem happens. To solve this issue we have handy tool called `DataLoader`. Essentially what it does is wait for all your resolvers to load in their individual keys. Once it has them, it hits the DB once with the keys, and returns a promise that resolves an array of the values. It batches our queries instead of making one at a time.

We used Golang DataLoader lib [github.com/graph-gophers/dataloader](https://github.com/graph-gophers/dataloader)
