package types

import (
	log "github.com/sirupsen/logrus"
	"regexp"
	"time"
)

const (
	Delete RuleResult = iota
	Close
	Save
)

type RuleResult int

type DeleteRule struct {
	DeleteIntervalDay time.Duration `mapstructure:"delete-interval-day"`
	CloseIntervalDay  time.Duration `mapstructure:"close-interval-day"`
	Indices           struct {
		Include Rule
		Exclude Rule
	}
}

type Rule struct {
	Keys      string `mapstructure:"keys"`
	Regex     string `mapstructure:"regex"`
	KeysRegex string
}

func (d *DeleteRule) NeedDeleteOrClose(index Index) RuleResult {
	indexDuration := time.Now().Sub(index.CreateTime)
	switch {
	case indexDuration <= time.Hour*24*d.CloseIntervalDay:
		{
			log.Infof("创建日期在 %d 天内，保留索引 %s", d.CloseIntervalDay, index.Name)
			return Save
		}
	case indexDuration <= time.Hour*24*d.DeleteIntervalDay:
		{
			log.Infof("创建日期在 %d-%d 天内，关闭索引 %s", d.CloseIntervalDay, d.DeleteIntervalDay, index.Name)
			return Close
		}
	case indexDuration >= time.Hour*24*90:
		{
			log.Infof("创建日期超过 90 天，强制删除索引 %s", index.Name)
			return Delete
		}
	default:
		{
			log.Infof("创建日期超过 %d 天，校验是否删除索引 %s", d.DeleteIntervalDay, index.Name)
			break
		}
	}

	if ruleCheck(&d.Indices.Exclude, index.Name) {
		log.Infof("符合排除条件，保留索引 %s", index.Name)
		return Save
	}
	if ruleCheck(&d.Indices.Include, index.Name) {
		log.Infof("符合过滤条件，删除索引 %s", index.Name)
		return Delete
	}

	log.Infof("不符合过滤条件，保留索引 %s", index.Name)
	return Save
}

func ruleCheck(rule *Rule, indexName string) bool {
	return regularCheck(rule.Regex, indexName) || regularCheck(rule.KeysRegex, indexName)
}

// 正则匹配返回 true
func regularCheck(reg, src string) bool {
	match, err := regexp.MatchString(reg, src)
	if err != nil {
		log.Error(err)
	}
	return match
}
