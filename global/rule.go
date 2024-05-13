package global

type Target string

const (
	TargetMongoDB Target = "mongodb"
	TargetRedis   Target = "redis"
)

type Rule struct {
	Schema        string `mapstructure:"schema" validate:"required" `
	Table         string `mapstructure:"table"  validate:"required"`
	Target        Target `mapstructure:"target"  validate:"required,oneof=mongodb redis"`
	OrderByColumn string `mapstructure:"order_by_column"`

	IncludeColumnsConfig []string `mapstructure:"include_columns"`
	ExcludeColumnsConfig []string `mapstructure:"exclude_columns"`
	ValueEncoder         string   `mapstructure:"value_encoder"`

	// ------------------- MongoDB -----------------
	MongodbDatabase      string           `mapstructure:"mongodb_database" validate:"required_if=Target mongodb"`
	MongodbCollection    string           `mapstructure:"mongodb_collection" validate:"required_if=Target mongodb"`
	ColumnMappingsConfig []*ColumnMapping `mapstructure:"column_mappings" validate:"dive"`
	NewColumnsConfig     []*NewColumn     `mapstructure:"new_columns" validate:"dive"`
	// ------------------- Redis -----------------
	RedisStructure   string     `mapstructure:"redis_structure"`
	RedisKeyConfig   RedisKey   `mapstructure:"redis_key"`
	RedisValueConfig RedisValue `mapstructure:"redis_value"`
}

type ColumnMapping struct {
	Source string `mapstructure:"source" validate:"required"`
	Target string `mapstructure:"target" validate:"required"`
}

type NewColumn struct {
	Name  string `mapstructure:"name" validate:"required"`
	Type  string `mapstructure:"type" validate:"required,oneof=int bool string"`
	Templ bool   `mapstructure:"templ"`
	Value string `mapstructure:"value"`
}

type RedisKey struct {
	Templ bool   `mapstructure:"templ"`
	Value string `mapstructure:"value"`
}

type RedisValue struct {
	Templ bool   `mapstructure:"templ"`
	Value string `mapstructure:"value"`
}
