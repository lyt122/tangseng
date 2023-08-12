package index

import (
	"time"

	"github.com/CocaineCong/tangseng/app/search_engine/engine"
	"github.com/CocaineCong/tangseng/app/search_engine/recall"
	"github.com/CocaineCong/tangseng/app/search_engine/types"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// Recall 召回
type Recall struct {
	*recall.Recall
}

// NewRecallServ 创建召回服务
func NewRecallServ(meta *engine.Meta) *Recall {
	r := recall.NewRecall(meta)
	return &Recall{r}
}

func SearchRecall(query string) (res []*types.SearchItem, err error) {
	meta, err := engine.ParseMeta()
	if err != nil {
		log.LogrusObj.Errorf("SearchRecall-ParseMeta:%+v", err)
		return
	}

	// 定时同步meta数据
	ticker := time.NewTicker(time.Second * 10)
	go meta.SyncByTicker(ticker)
	recallService := NewRecallServ(meta)
	res, err = recallService.Search(query)
	if err != nil {
		log.LogrusObj.Errorf("SearchRecall-NewRecallServ:%+v", err)
		return
	}
	// close
	// func() {
	// 	// 最后同步元数据至文件
	// 	fmt.Println("close")
	// 	meta.SyncMeta()
	// 	fmt.Println("close")
	// 	ticker.Stop()
	// 	fmt.Println("close")
	// }()

	return
}