# guacinator/cmd

The `cmd` package is a collection of utility functions
designed to simplify common cmd tasks.

Table of contents:

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### CreateAdminUser

```go
CreateAdminUser(string) error
```

CreateAdminUser creates a new admin user in Guacamole
with the specified username and password.

**Parameters:**

user: A string representing the desired username for the new admin user.
password: A string representing the desired password for the new admin user.

**Returns:**

error: An error if the admin user cannot be created.

### CreateGuacamoleConnection

```go
CreateGuacamoleConnection(VncHost) error
```

CreateGuacamoleConnection establishes a new connection
in Guacamole using the provided VncHost information.

**Parameters:**

vncHost: A VncHost struct containing the necessary information for the connection.

**Returns:**

error: An error if the connection cannot be created.

### DeleteGuacUser

```go
DeleteGuacUser(string) error
```

DeleteGuacUser removes a specified Guacamole user.

**Parameters:**

user: A string representing the username of the Guacamole user to be deleted.

**Returns:**

error: An error if the specified user cannot be deleted.

### Execute

```go
Execute()
```

Execute runs the root cobra command

---

## Installation

To use the guacinator/cmd package, you first need to install it.
Follow the steps below to install via go get.

```bash
go get github.com/guacinator/cowdogmoo/cmd
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/guacinator/cowdogmoo/cmd"
```

---

## Tests

To ensure the package is working correctly, run the following
command to execute the tests for `guacinator/cmd`:

```bash
go test -v
```

---

## Contributing

Pull requests are welcome. For major changes,
please open an issue first to discuss what
you would like to change.

---

## License

This project is licensed under the MIT
License - see the [LICENSE](../LICENSE)
file for details.
