# How to run scripts?

## Create new schemas

Create a new entity schema definition in the `src/schema` directory with the given name.

### Command:

```command
go run scripts/schema/main.go <NAME>
```

Replace `<NAME>` with the name of the schema. Example (creates a `HelloWorld` schema):

```command
go run scripts/schema/main.go HelloWorld
```

## Ent codegen

Generates the ent boilerplate: entity definitions, basic query CRUD functions, etc.

### Command:

```command
go run scripts/codegen/main.go
```

## Database migration

Migrates the database schema to the latest version, creating tables, indexes, etc.

### Command:

```command
go run scripts/migrate/main.go
```
