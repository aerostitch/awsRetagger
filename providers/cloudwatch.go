package providers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/sirupsen/logrus"

	"github.com/VEVO/awsRetagger/mapper"
)

// CwProcessor holds the cloudwatch-related actions
type CwProcessor struct {
	svc *cloudwatchlogs.CloudWatchLogs
}

// NewCwProcessor creates a new instance of CwProcessor containing an already
// initialized cloudwatchlogs client
func NewCwProcessor(sess *session.Session) *CwProcessor {
	return &CwProcessor{svc: cloudwatchlogs.New(sess)}
}

// TagsToMap transform the cw tags structure into a map[string]string for
// easier manipulations
func (p *CwProcessor) TagsToMap(tagsInput map[string]*string) map[string]string {
	tagsHash := make(map[string]string)
	for k, v := range tagsInput {
		tagsHash[k] = *v
	}
	return tagsHash
}

// SetTags sets tags on an cloudwatchLog resource
func (p *CwProcessor) SetTags(resourceID *string, tags []*mapper.TagItem) error {
	newTags := map[string]*string{}
	for _, tag := range tags {
		newTags[(*tag).Name] = aws.String((*tag).Value)
	}

	input := &cloudwatchlogs.TagLogGroupInput{
		LogGroupName: resourceID,
		Tags:         newTags,
	}
	_, err := p.svc.TagLogGroup(input)
	return err
}

// GetTags gets the tags allocated to an rds resource
func (p *CwProcessor) GetTags(resourceID *string) (map[string]*string, error) {
	input := &cloudwatchlogs.ListTagsLogGroupInput{
		LogGroupName: resourceID,
	}

	result, err := p.svc.ListTagsLogGroup(input)
	return result.Tags, err
}

// RetagLogGroups parses all running and stopped instances and retags them
func (p *CwProcessor) RetagLogGroups(m *mapper.Mapper) {
	err := p.svc.DescribeLogGroupsPages(&cloudwatchlogs.DescribeLogGroupsInput{},
		func(page *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
			for _, lg := range page.LogGroups {
				t, err := p.GetTags(lg.LogGroupName)
				if err != nil {
					log.WithFields(logrus.Fields{"error": err, "resource": *lg.Arn}).Fatal("Failed to get LogGroup tags")
				}

				tags := p.TagsToMap(t)
				keys := []string{}
				if lg.LogGroupName != nil {
					keys = append(keys, *lg.LogGroupName)
				}
				m.Retag(lg.LogGroupName, &tags, keys, p.SetTags)
			}
			return !lastPage
		})
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Fatal("DescribeLogGroups failed")
	}
}
