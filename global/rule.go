package global

type Target string

const (
	TargetMongoDB Target = "mongodb"
	TargetMySQL   Target = "mysql"
)

type Rule struct {
	Schema                   string `mapstructure:"schema" validate:"required" `
	Table                    string `mapstructure:"table"  validate:"required"`
	Target                   Target `mapstructure:"target"  validate:"required,oneof=mongodb mysql"`
	OrderByColumn            string `mapstructure:"order_by_column"`
	ColumnLowerCase          bool   `mapstructure:"column_lower_case"`          // 列名称转为小写
	ColumnUpperCase          bool   `mapstructure:"column_upper_case"`          // 列名称转为大写
	ColumnUnderscoreToCamel  bool   `mapstructure:"column_underscore_to_camel"` // 列名称下划线转驼峰
	IncludeColumnConfig      string `mapstructure:"include_columns"`            // 包含的列
	ExcludeColumnConfig      string `mapstructure:"exclude_columns"`            // 排除掉的列
	ColumnMappingConfigs     string `mapstructure:"column_mappings"`            // 列名称映射
	DefaultColumnValueConfig string `mapstructure:"default_column_values"`      // 默认的字段和值
	// #值编码，支持json、kv-commas、v-commas；默认为json；json形如：{"id":123,"name":"wangjie"} 、kv-commas形如：id=123,name="wangjie"、v-commas形如：123,wangjie
	ValueEncoder      string `mapstructure:"value_encoder"`
	ValueFormatter    string `mapstructure:"value_formatter"`    //格式化定义key,{id}表示字段id的值、{name}表示字段name的值
	LuaScript         string `mapstructure:"lua_script"`         //lua 脚本
	LuaFilePath       string `mapstructure:"lua_file_path"`      //lua 文件地址
	DateFormatter     string `mapstructure:"date_formatter"`     //date类型格式化， 不填写默认2006-01-02
	DatetimeFormatter string `mapstructure:"datetime_formatter"` //datetime、timestamp类型格式化，不填写默认RFC3339(2006-01-02T15:04:05Z07:00)

	// ------------------- MONGODB -----------------
	MongodbDatabase   string `mapstructure:"mongodb_database" validate:"required_if=Target mongodb"`
	MongodbCollection string `mapstructure:"mongodb_collection" validate:"required_if=Target mongodb"`
}
