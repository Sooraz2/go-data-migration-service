package models

type Config struct {
    Source      DBConfig      `yaml:"source"`
    Destination DBConfig      `yaml:"destination"`
    Migration   MigrationConfig `yaml:"migration"`
}

type DBConfig struct {
    Connection string `yaml:"connection"`
    Host      string `yaml:"host"`
    Port      int    `yaml:"port"`
    Database  string `yaml:"database"`
    Username  string `yaml:"username"`
    Password  string `yaml:"password"`
}

type MigrationConfig struct {
    MainTable     TableConfig   `yaml:"main_table"`
    RelatedTables []TableConfig `yaml:"related_tables"`
}

type TableConfig struct {
    Name       string   `yaml:"name"`
    ForeignKey string   `yaml:"foreign_key,omitempty"`
    Filter     Filter   `yaml:"filter,omitempty"`
    Columns    []string `yaml:"columns"`
}

type Filter struct {
    Column string `yaml:"column"`
    Value  string `yaml:"value"`
}