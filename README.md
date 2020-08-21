# Go GraphQL Server

A basic GraphQL server based on following Schema.

```
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

## Mutation

Create Author

```
curl -g 'http://localhost:8080/graphql?query=mutation+_{author(name:"pavel",isbn_no:"19923"){id,name,isbn_no}}'
```

or

```
mutation {
  author(name:"tarique", isbn_no:"5"){
    id
    name
    isbn_no
  }
}
```

curl -g 'http://localhost:8080/graphql?query={authors{id,name}}'
