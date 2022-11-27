package main

import (
	"context"

	"gorm.io/gen"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/repo/storage/po"
)

type CronTemplateDAO interface {
	// SELECT * FROM @@table WHERE id=@id
	FindByID(id uint) (*gen.T, error)

	// UPDATE @@table
	//  {{set}}
	//    {{if status > 0}} status=@status, {{end}}
	//  {{end}}
	// WHERE id=@id
	UpdateStatus(ctx context.Context, id uint, status pb.TriggerStatus) (gen.RowsAffected, error)
}

func main() {

	g := gen.NewGenerator(gen.Config{
		OutPath: "../sql",
		Mode:    gen.WithDefaultQuery,
	})

	g.ApplyBasic(po.CronTriggerTemplate{})

	g.ApplyInterface(func(CronTemplateDAO) {}, po.CronTriggerTemplate{})
	g.Execute()

}
