# guacinator/cmd

The `cmd` package provides guacamole CLI utilities.

---

## Table of contents

- [Functions](#functions)
- [Installation](#installation)
- [Usage](#usage)
- [Tests](#tests)
- [Contributing](#contributing)
- [License](#license)

---

## Functions

### Execute()

```go
Execute()
```

Execute runs the root cobra command. It checks for errors and exits
the program if any are encountered.

---

### GuacServiceImpl.CreateAdminUser(string)

```go
CreateAdminUser(string) error
```

CreateAdminUser creates a new admin user in
Guacamole with the specified
username and password.

**Parameters:**

user: A string representing the desired username for the new admin user.

password: A string representing the desired password for the new admin user.

**Returns:**

error: An error if the admin user cannot be created.

---

### GuacServiceImpl.CreateGuacamoleConnection(VncHost)

```go
CreateGuacamoleConnection(VncHost) error
```

CreateGuacamoleConnection establishes a new connection

in Guacamole using the provided VncHost information.
CreatePackageDocs generates documentation for all Go packages in the current
directory and its subdirectories. It traverses the file tree using a provided
afero.Fs and Repo to create a new README.md file in each directory containing
a Go package. It uses a specified template file for generating the README files.

**Parameters:**

vncHost: A VncHost struct containing the necessary information for the connection.

**Returns:**

error: An error if the connection cannot be created.

---

### GuacServiceImpl.DeleteGuacUser(string)

```go
DeleteGuacUser(string) error
```

DeleteGuacUser removes a specified Guacamole user.

**Parameters:**

user: A string representing the username of the Guacamole user to be deleted.

**Returns:**

error: An error if the specified user cannot be deleted.

---

## Installation

To use the guacinator/cmd package, you first need to install it.
Follow the steps below to install via go install.

```bash
go install github.com/cowdogmoo/guacinator/cmd@latest
```

---

## Usage

After installation, you can import the package in your Go project
using the following import statement:

```go
import "github.com/cowdogmoo/guacinator/cmd"
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
License - see the [LICENSE](https://github.com/CowDogMoo/guacinator/blob/main/LICENSE)
file for details.
