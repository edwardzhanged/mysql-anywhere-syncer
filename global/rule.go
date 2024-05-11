package global

type Target string

const (
	TargetMongoDB Target = "mongodb"
	TargetMySQL   Target = "mysql"
)

type Rule struct {
	Schema        string `mapstructure:"schema" validate:"required" `
	Table         string `mapstructure:"table"  validate:"required"`
	Target        Target `mapstructure:"target"  validate:"required,oneof=mongodb mysql"`
	OrderByColumn string `mapstructure:"order_by_column"`

	IncludeColumnsConfig []string         `mapstructure:"include_columns"`
	ExcludeColumnsConfig []string         `mapstructure:"exclude_columns"`
	ColumnMappingsConfig []*ColumnMapping `mapstructure:"column_mappings" validate:"dive"`
	NewColumnsConfig     []*NewColumn     `mapstructure:"new_columns" validate:"dive"`
	ValueEncoder         string           `mapstructure:"value_encoder"`
	ValueFormatter       string           `mapstructure:"value_formatter"`

	// ------------------- MONGODB -----------------
	MongodbDatabase   string `mapstructure:"mongodb_database" validate:"required_if=Target mongodb"`
	MongodbCollection string `mapstructure:"mongodb_collection" validate:"required_if=Target mongodb"`
	// ------------------- MYSQL -----------------
	MysqlDatabase string `mapstructure:"mysql_database" validate:"required_if=Target mysql"`
	MysqlTable    string `mapstructure:"mysql_table" validate:"required_if=Target mysql"`
}

type ColumnMapping struct {
	Source string `mapstructure:"source" validate:"required"`
	Target string `mapstructure:"target" validate:"required"`
}

type NewColumn struct {
	Name  string `mapstructure:"name" validate:"required"`
	Type  string `mapstructure:"type" validate:"required,oneof=int bool string"`
	Value string `mapstructure:"value"`
}
