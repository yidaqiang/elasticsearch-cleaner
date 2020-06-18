package types

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ElasticSearchManager struct {
	Server     string
	Username   string
	password   string
	DryRun     bool `mapstructure:"dry-run"`
	DeleteRule `mapstructure:"delete-policy"`
	Client     *elastic.Client
}

func (m *ElasticSearchManager) Init() {
	var err error

	m.Client, err = elastic.NewClient(elastic.SetURL(m.Server), elastic.SetSniff(false))

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// 根据关键字生成正则表达式
	m.DeleteRule.Indices.Include.KeysRegex = generateRegexByKeys(m.DeleteRule.Indices.Include.Keys)
	m.DeleteRule.Indices.Exclude.KeysRegex = generateRegexByKeys(m.DeleteRule.Indices.Exclude.Keys)

	// 值为空时设置为无法匹配到的字符
	if len(m.DeleteRule.Indices.Include.Regex) == 0 {
		m.DeleteRule.Indices.Include.Regex = ":"
	}
	if len(m.DeleteRule.Indices.Include.KeysRegex) == 0 {
		m.DeleteRule.Indices.Include.KeysRegex = ":"
	}
	if len(m.DeleteRule.Indices.Exclude.Regex) == 0 {
		m.DeleteRule.Indices.Exclude.Regex = ":"
	}
	if len(m.DeleteRule.Indices.Exclude.KeysRegex) == 0 {
		m.DeleteRule.Indices.Exclude.KeysRegex = ":"
	}
}

func (m *ElasticSearchManager) Indices() (result Indices) {
	idx, err := m.Client.CatIndices().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	var indicesName []string
	for _, i := range idx {
		indicesName = append(indicesName, i.Index)
	}
	if len(indicesName) == 0 {
		return nil
	}

	return m.indicesObjs(indicesName)
}

func (m *ElasticSearchManager) DeleteIndex(index string) {
	if m.DryRun {
		log.Infof("跳过删除索引操作")
		return
	}
	res, err := m.Client.DeleteIndex(index).Do(context.Background())
	if err != nil {
		log.Error(err)
	}
	if res.Acknowledged {
		log.Infof("删除索引 %s 成功", index)
	} else {
		log.Infof("删除索引 %s失败", index)
	}
}

func (m *ElasticSearchManager) CloseIndex(index string) {
	if m.DryRun {
		log.Infof("跳过关闭索引操作")
		return
	}
	res, err := m.Client.CloseIndex(index).Do(context.Background())
	if err != nil {
		log.Error(err)
	}
	if res.Acknowledged {
		log.Infof("关闭索引 %s 成功", index)
	} else {
		log.Infof("关闭索引 %s失败", index)
	}
}

// 通过关键字生成正则表达式
func generateRegexByKeys(keys string) string {
	if strings.HasPrefix(keys, ",") {
		keys = keys[1:]
	}
	if strings.HasSuffix(keys, ",") {
		keys = keys[:len(keys)-1]
	}
	if len(keys) == 0 || len(strings.Replace(keys, ",", "", -1)) == 0 {
		return ":"
	}
	return fmt.Sprintf(".*%s.*", strings.Replace(keys, ",", ".*|.*", -1))
}

func (m *ElasticSearchManager) indicesObjs(indicesName []string) Indices {
	indexObjs := make([]Index, 0)
	for _, idx := range indicesName {
		index, err := m.Client.IndexGetSettings(idx).Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		creationDate := index[idx].Settings["index"].(map[string]interface{})["creation_date"]

		d, _ := strconv.Atoi(creationDate.(string))
		// creation_date 表示毫秒的 unix 时间
		t := time.Unix(int64(d/1000), 0)
		newIndex := Index{
			Name:       idx,
			CreateTime: t,
		}
		indexObjs = append(indexObjs, newIndex)
	}
	sort.Sort(Indices(indexObjs))

	return indexObjs
}
