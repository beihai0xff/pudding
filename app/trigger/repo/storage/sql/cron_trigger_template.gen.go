// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package sql

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gen/helper"

	"gorm.io/plugin/dbresolver"

	"github.com/beihai0xff/pudding/app/trigger/repo/storage/po"
)

func newCronTriggerTemplate(db *gorm.DB, opts ...gen.DOOption) cronTriggerTemplate {
	_cronTriggerTemplate := cronTriggerTemplate{}

	_cronTriggerTemplate.cronTriggerTemplateDo.UseDB(db, opts...)
	_cronTriggerTemplate.cronTriggerTemplateDo.UseModel(&po.CronTriggerTemplate{})

	tableName := _cronTriggerTemplate.cronTriggerTemplateDo.TableName()
	_cronTriggerTemplate.ALL = field.NewAsterisk(tableName)
	_cronTriggerTemplate.ID = field.NewUint(tableName, "id")
	_cronTriggerTemplate.CreatedAt = field.NewTime(tableName, "created_at")
	_cronTriggerTemplate.UpdatedAt = field.NewTime(tableName, "updated_at")
	_cronTriggerTemplate.DeletedAt = field.NewField(tableName, "deleted_at")
	_cronTriggerTemplate.CronExpr = field.NewString(tableName, "cron_expr")
	_cronTriggerTemplate.Topic = field.NewString(tableName, "topic")
	_cronTriggerTemplate.Payload = field.NewBytes(tableName, "payload")
	_cronTriggerTemplate.LastExecutionTime = field.NewTime(tableName, "last_execution_time")
	_cronTriggerTemplate.ExceptedEndTime = field.NewTime(tableName, "excepted_end_time")
	_cronTriggerTemplate.ExceptedLoopTimes = field.NewUint64(tableName, "excepted_loop_times")
	_cronTriggerTemplate.LoopedTimes = field.NewUint64(tableName, "looped_times")
	_cronTriggerTemplate.Status = field.NewInt(tableName, "status")

	_cronTriggerTemplate.fillFieldMap()

	return _cronTriggerTemplate
}

type cronTriggerTemplate struct {
	cronTriggerTemplateDo cronTriggerTemplateDo

	ALL               field.Asterisk
	ID                field.Uint
	CreatedAt         field.Time
	UpdatedAt         field.Time
	DeletedAt         field.Field
	CronExpr          field.String
	Topic             field.String
	Payload           field.Bytes
	LastExecutionTime field.Time
	ExceptedEndTime   field.Time
	ExceptedLoopTimes field.Uint64
	LoopedTimes       field.Uint64
	Status            field.Int

	fieldMap map[string]field.Expr
}

func (c cronTriggerTemplate) Table(newTableName string) *cronTriggerTemplate {
	c.cronTriggerTemplateDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c cronTriggerTemplate) As(alias string) *cronTriggerTemplate {
	c.cronTriggerTemplateDo.DO = *(c.cronTriggerTemplateDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *cronTriggerTemplate) updateTableName(table string) *cronTriggerTemplate {
	c.ALL = field.NewAsterisk(table)
	c.ID = field.NewUint(table, "id")
	c.CreatedAt = field.NewTime(table, "created_at")
	c.UpdatedAt = field.NewTime(table, "updated_at")
	c.DeletedAt = field.NewField(table, "deleted_at")
	c.CronExpr = field.NewString(table, "cron_expr")
	c.Topic = field.NewString(table, "topic")
	c.Payload = field.NewBytes(table, "payload")
	c.LastExecutionTime = field.NewTime(table, "last_execution_time")
	c.ExceptedEndTime = field.NewTime(table, "excepted_end_time")
	c.ExceptedLoopTimes = field.NewUint64(table, "excepted_loop_times")
	c.LoopedTimes = field.NewUint64(table, "looped_times")
	c.Status = field.NewInt(table, "status")

	c.fillFieldMap()

	return c
}

func (c *cronTriggerTemplate) WithContext(ctx context.Context) *cronTriggerTemplateDo {
	return c.cronTriggerTemplateDo.WithContext(ctx)
}

func (c cronTriggerTemplate) TableName() string { return c.cronTriggerTemplateDo.TableName() }

func (c cronTriggerTemplate) Alias() string { return c.cronTriggerTemplateDo.Alias() }

func (c *cronTriggerTemplate) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *cronTriggerTemplate) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 12)
	c.fieldMap["id"] = c.ID
	c.fieldMap["created_at"] = c.CreatedAt
	c.fieldMap["updated_at"] = c.UpdatedAt
	c.fieldMap["deleted_at"] = c.DeletedAt
	c.fieldMap["cron_expr"] = c.CronExpr
	c.fieldMap["topic"] = c.Topic
	c.fieldMap["payload"] = c.Payload
	c.fieldMap["last_execution_time"] = c.LastExecutionTime
	c.fieldMap["excepted_end_time"] = c.ExceptedEndTime
	c.fieldMap["excepted_loop_times"] = c.ExceptedLoopTimes
	c.fieldMap["looped_times"] = c.LoopedTimes
	c.fieldMap["status"] = c.Status
}

func (c cronTriggerTemplate) clone(db *gorm.DB) cronTriggerTemplate {
	c.cronTriggerTemplateDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c cronTriggerTemplate) replaceDB(db *gorm.DB) cronTriggerTemplate {
	c.cronTriggerTemplateDo.ReplaceDB(db)
	return c
}

type cronTriggerTemplateDo struct{ gen.DO }

//SELECT * FROM @@table WHERE id=@id
func (c cronTriggerTemplateDo) FindByID(id uint) (result *po.CronTriggerTemplate, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("SELECT * FROM cron_trigger_template WHERE id=? ")

	var executeSQL *gorm.DB

	executeSQL = c.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result)
	err = executeSQL.Error
	return
}

//UPDATE @@table
// {{set}}
//   {{if status > 0}} status=@status, {{end}}
// {{end}}
//WHERE id=@id
func (c cronTriggerTemplateDo) UpdateStatus(ctx context.Context, id uint, status int) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("UPDATE cron_trigger_template ")
	var setSQL0 strings.Builder
	if status > 0 {
		params = append(params, status)
		setSQL0.WriteString("status=?, ")
	}
	helper.JoinSetBuilder(&generateSQL, setSQL0)
	params = append(params, id)
	generateSQL.WriteString("WHERE id=? ")

	var executeSQL *gorm.DB

	executeSQL = c.UnderlyingDB().Exec(generateSQL.String(), params...)
	err = executeSQL.Error
	return
}

func (c cronTriggerTemplateDo) Debug() *cronTriggerTemplateDo {
	return c.withDO(c.DO.Debug())
}

func (c cronTriggerTemplateDo) WithContext(ctx context.Context) *cronTriggerTemplateDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c cronTriggerTemplateDo) ReadDB() *cronTriggerTemplateDo {
	return c.Clauses(dbresolver.Read)
}

func (c cronTriggerTemplateDo) WriteDB() *cronTriggerTemplateDo {
	return c.Clauses(dbresolver.Write)
}

func (c cronTriggerTemplateDo) Session(config *gorm.Session) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Session(config))
}

func (c cronTriggerTemplateDo) Clauses(conds ...clause.Expression) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c cronTriggerTemplateDo) Returning(value interface{}, columns ...string) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c cronTriggerTemplateDo) Not(conds ...gen.Condition) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c cronTriggerTemplateDo) Or(conds ...gen.Condition) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c cronTriggerTemplateDo) Select(conds ...field.Expr) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c cronTriggerTemplateDo) Where(conds ...gen.Condition) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c cronTriggerTemplateDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *cronTriggerTemplateDo {
	return c.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (c cronTriggerTemplateDo) Order(conds ...field.Expr) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c cronTriggerTemplateDo) Distinct(cols ...field.Expr) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c cronTriggerTemplateDo) Omit(cols ...field.Expr) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c cronTriggerTemplateDo) Join(table schema.Tabler, on ...field.Expr) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c cronTriggerTemplateDo) LeftJoin(table schema.Tabler, on ...field.Expr) *cronTriggerTemplateDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c cronTriggerTemplateDo) RightJoin(table schema.Tabler, on ...field.Expr) *cronTriggerTemplateDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c cronTriggerTemplateDo) Group(cols ...field.Expr) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c cronTriggerTemplateDo) Having(conds ...gen.Condition) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c cronTriggerTemplateDo) Limit(limit int) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c cronTriggerTemplateDo) Offset(offset int) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c cronTriggerTemplateDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c cronTriggerTemplateDo) Unscoped() *cronTriggerTemplateDo {
	return c.withDO(c.DO.Unscoped())
}

func (c cronTriggerTemplateDo) Create(values ...*po.CronTriggerTemplate) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c cronTriggerTemplateDo) CreateInBatches(values []*po.CronTriggerTemplate, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c cronTriggerTemplateDo) Save(values ...*po.CronTriggerTemplate) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c cronTriggerTemplateDo) First() (*po.CronTriggerTemplate, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*po.CronTriggerTemplate), nil
	}
}

func (c cronTriggerTemplateDo) Take() (*po.CronTriggerTemplate, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*po.CronTriggerTemplate), nil
	}
}

func (c cronTriggerTemplateDo) Last() (*po.CronTriggerTemplate, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*po.CronTriggerTemplate), nil
	}
}

func (c cronTriggerTemplateDo) Find() ([]*po.CronTriggerTemplate, error) {
	result, err := c.DO.Find()
	return result.([]*po.CronTriggerTemplate), err
}

func (c cronTriggerTemplateDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*po.CronTriggerTemplate, err error) {
	buf := make([]*po.CronTriggerTemplate, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c cronTriggerTemplateDo) FindInBatches(result *[]*po.CronTriggerTemplate, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c cronTriggerTemplateDo) Attrs(attrs ...field.AssignExpr) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c cronTriggerTemplateDo) Assign(attrs ...field.AssignExpr) *cronTriggerTemplateDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c cronTriggerTemplateDo) Joins(fields ...field.RelationField) *cronTriggerTemplateDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c cronTriggerTemplateDo) Preload(fields ...field.RelationField) *cronTriggerTemplateDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c cronTriggerTemplateDo) FirstOrInit() (*po.CronTriggerTemplate, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*po.CronTriggerTemplate), nil
	}
}

func (c cronTriggerTemplateDo) FirstOrCreate() (*po.CronTriggerTemplate, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*po.CronTriggerTemplate), nil
	}
}

func (c cronTriggerTemplateDo) FindByPage(offset int, limit int) (result []*po.CronTriggerTemplate, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c cronTriggerTemplateDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c cronTriggerTemplateDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c cronTriggerTemplateDo) Delete(models ...*po.CronTriggerTemplate) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *cronTriggerTemplateDo) withDO(do gen.Dao) *cronTriggerTemplateDo {
	c.DO = *do.(*gen.DO)
	return c
}