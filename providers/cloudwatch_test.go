package providers

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func TestCwTagsToMap(t *testing.T) {
	testData := []struct {
		inputTags  map[string]*string
		outputTags map[string]string
	}{
		{map[string]*string{}, map[string]string{}},
		{map[string]*string{"foo": aws.String("bar")}, map[string]string{"foo": "bar"}},
		{map[string]*string{"foo": aws.String("bar"), "Aerosmith": aws.String("rocks")}, map[string]string{"foo": "bar", "Aerosmith": "rocks"}},
	}
	for _, d := range testData {
		p := CwProcessor{}

		res := p.TagsToMap(d.inputTags)
		if !reflect.DeepEqual(res, d.outputTags) {
			t.Errorf("Expecting to get tags: %v\nGot: %v\n", d.outputTags, res)
		}
	}
}
