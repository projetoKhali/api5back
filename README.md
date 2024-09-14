# api5back

## Install [`Go`](https://golang.org/dl/)

- Linux (Ubuntu / Mint)

```command
sudo apt install golang-go
```

- Linux (Arch / Garuda)

```command
sudo pacman -S go
```

- Windows  
  https://go.dev/doc/install

## Install [`Air`](https://github.com/air-verse/air)

Air is a tool that provides a live reload server for Go applications, updating
the running server when the source code changes, similar to `nodemon` for Node.js.

```command
go install github.com/air-verse/air@latest
```

## Setup [`Husky`](https://github.com/automation-co/husky)

Husky is a tool that provides a simple way to manage the project's Git hooks.

The hooks are used to ensure code style consistency, run tests, and semantic
commit messages.

### Install Husky

This will setup the `husky` command globally.

```command
go install github.com/automation-co/husky@latest
```

### Activate the hooks

This will copy the hooks to the `.git/hooks` directory from the `.husky/hooks` directory.

```command
husky install
```
