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
  book(title:"A Book of Fire", description:"A world famous book", authors:["154"]){
    id
    title
    description
  }
}
```

curl -g 'http://localhost:8080/graphql?query={authors{id,name}}'
