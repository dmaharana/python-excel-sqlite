# go-excel-sqlite

This project is a simple backup of excel data into a SQLite database. The backup process involves reading an Excel file and writing the data into a SQLite database.

## Usage

1. Make sure you have Go installed.
2. Clone the repository.
3. Install the required dependencies by running `go mod tidy` in the project directory.
4. Run the backup script by running `go run main.go` in the project directory.

## Requirements

- Go version 1.16 or higher
- github.com/tealeg/xlsx
- github.com/uptrace/bun
- github.com/uptrace/bun/dialect/sqlitedialect

## Project Structure

The project is structured as follows:

- `main.go`: contains the code to read Excel data and write it to a SQLite database.
- `User`: a struct representing a user in the SQLite database.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
