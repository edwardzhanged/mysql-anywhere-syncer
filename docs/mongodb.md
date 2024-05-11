# MongoDB Configuration
Add following configurations to your yaml file
## MongoDB Info
```yaml
mongodb_host: 10.136.134.218
mongodb_port: 27017
mongodb_username: root
mongodb_password: 123456
```

## Rules
```yaml
rules:
  - schema: notify # MySQL schema
    table: students 
    target: mongodb # Tag this rule as mongodb
    mongodb_database: test
    mongodb_collection: students
    include_columns: # Not necessary; can pick needed columns to target only
      - name
      - age
      - enrollment_date
    exclude_columns: address # Not necessary; can exclude unneeded columns
    column_mappings: # Not necessary; map column to a new name
      - source: mobile
        target: phone
    new_columns:  # Not necessary; add new column with default value, support int/bool/string
      - name: status
        type: int
        value: 1
      - name: male
        type: bool
        value: false
  - schema: notify2
  ....more rules
```