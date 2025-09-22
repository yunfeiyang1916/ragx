package internal

// 生成entity(与表结构一一对应)及gorm query
//go:generate gentoolplus -c gentoolplus_config.json

// 生成biz与data
//go:generate gen --connstr "postgres://root:root@localhost:5432/ragx?sslmode=disable" --sqltype=postgres --database ragx --gorm --model=entity --dao=ragx  --exec=../pkg/codegen/template/repo.gen --out ./data --templateDir=../pkg/codegen/template --json --json-fmt=lower_camel   --exclude=schema_migrations --overwrite
