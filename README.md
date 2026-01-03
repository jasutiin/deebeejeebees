got curious about how a database actually works under the hood so why not make my own. will make it easier for me to understand why doing a bunch of JOINS is inefficient, or why doing 'SELECT \*' is slow.

what it does so far:

- lexical analysis, which is converting the sql query into tokens
- building a parse tree (cst)

todo:

- turning parse tree into abstract syntax tree (ast)
- semantic analysis, ensures the query is meaningful and correct
  - name resolution to check if the identifiers exist in the schema
  - type checking
  - constraint checking
- query planning
  - turn ast into relational algebra
  - optimization, how we actually do the query
- query execution
