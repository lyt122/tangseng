package woker

import (
	"context"
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/app/index_platform/analyzer"
	"github.com/CocaineCong/tangseng/app/index_platform/repository/db/dao"
	"github.com/CocaineCong/tangseng/app/index_platform/rpc"
	"github.com/CocaineCong/tangseng/app/index_platform/service/input_data_mr"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/mysql/db"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	db.InitDB()
	analyzer.InitSeg()
	rpc.Init()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestWorker(t *testing.T) {
	ctx := context.Background()
	dao.InitMysqlDirectUpload(ctx)
	Worker(ctx, input_data_mr.Map, input_data_mr.Reduce)
}