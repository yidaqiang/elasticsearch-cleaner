package action

import (
	log "github.com/sirupsen/logrus"
	yitypes "github.com/yidaqiang/elasticsearch-manager/pkg/types"
)

func RunDeleteIndices(manager *yitypes.ElasticSearchManager) {
	indices := manager.Indices()
	indicesCount := len(indices)
	log.Infof("获取到索引数量：%d", indicesCount)
	for i, index := range indices {
		log.Infof("当前处理第 %d/%d 索引：%s", i+1, indicesCount, index.Name)
		switch manager.NeedDeleteOrClose(index) {
		case yitypes.Delete:
			{
				manager.DeleteIndex(index.Name)
				break
			}
		case yitypes.Save:
			{
				// Do noting
				break
			}
		case yitypes.Close:
			{
				manager.CloseIndex(index.Name)
				break
			}
		default:
			{
				log.Warn("未知的操作码")
			}
		}
	}
}
