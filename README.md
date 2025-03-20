# Go Data Migration Tool

A robust data migration tool built in Go that allows for selective data migration between MySQL databases with support for related tables and custom filters.

## Features

- Configurable source and destination database connections
- Support for main table migration with custom filters
- Related tables migration with foreign key relationships
- Environment-based configuration
- Detailed logging of migration process

## Prerequisites

- Go 1.23 or higher
- MySQL/MariaDB
- Access to source and destination databases

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/go-data-migration.git
cd go-data-migration
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment configuration:
```bash
cp config/.env.example config/.env
```

## Configuration

### Environment Variables

1. Edit `config/.env` file with your database credentials:
```env
SOURCE_HOST=your_source_host
SOURCE_PORT=3306
SOURCE_DATABASE=your_source_db
SOURCE_USERNAME=your_username
SOURCE_PASSWORD=your_password

DEST_HOST=your_dest_host
DEST_PORT=3306
DEST_DATABASE=your_dest_db
DEST_USERNAME=your_username
DEST_PASSWORD=your_password
```

### Migration Configuration

Edit `config/config.yaml` to specify your migration settings:

```yaml
migration:
  main_table:
    name: "your_main_table"
    filter:
      column: "your_filter_column"
      value: "your_filter_value"
    columns: ["column1", "column2", "column3"]
  
  related_tables:
    - name: "related_table1"
      foreign_key: "fk_column"
      columns: ["column1", "column2"]
```

## Usage

1. Run the migration tool:
```bash
go run main.go
```

2. Monitor the logs for migration progress and any potential errors.

## Migration Process

1. The tool first migrates the main table based on the specified filter
2. It then migrates related tables using the foreign key relationships
3. Tables are created in the destination database if they don't exist
4. Data is copied while maintaining referential integrity

## Project Structure

```
.
├── config/
│   ├── config.yaml         # Migration configuration
│   ├── .env               # Database credentials (not in git)
│   └── .env.example       # Example environment configuration
├── models/
│   └── config.go          # Configuration structures
├── services/
│   ├── database.go        # Database connection handling
│   └── migration.go       # Migration logic
├── main.go                # Application entry point
├── go.mod                 # Go modules file
└── README.md             # This file
```

## Important Notes

- Never commit `.env` file to version control
- Always backup your destination database before running migrations
- Ensure proper database permissions for both source and destination
- Check logs for any errors during migration

## Error Handling

The tool provides detailed error logging. Common error messages and solutions:

- `Error connecting to database`: Check your database credentials and network connectivity
- `Table doesn't exist`: Verify table names in configuration
- `Permission denied`: Ensure proper database user permissions
- `Error during migration`: Check the SQL error message in logs for specific details

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details

## Support

For support, please open an issue in the GitHub repository or contact the maintainers.