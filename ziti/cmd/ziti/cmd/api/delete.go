package api

import (
	"github.com/fatih/color"
	"github.com/openziti/foundation/storage/boltz"
	"github.com/openziti/ziti/ziti/cmd/ziti/util"
	"net/url"
	"strings"
)

func DeleteEntitiesOfType(api util.API, o *Options, entityType string, ids []string, body string) error {
	for _, id := range ids {
		err := util.ControllerDelete(api, entityType, id, body, o.Out, o.OutputJSONRequest, o.OutputJSONResponse, o.Timeout, o.Verbose)
		if err != nil {
			o.Printf("delete of %v with id %v: %v\n", boltz.GetSingularEntityType(entityType), id, color.New(color.FgRed, color.Bold).Sprint("FAIL"))
			return err
		}
		o.Printf("delete of %v with id %v: %v\n", boltz.GetSingularEntityType(entityType), id, color.New(color.FgGreen, color.Bold).Sprint("OK"))
	}
	return nil
}

// DeleteEntityOfTypeWhere implements the commands to delete various entity types
func DeleteEntityOfTypeWhere(api util.API, options *Options, entityType string, body string) error {
	filter := strings.Join(options.Args, " ")

	params := url.Values{}
	params.Add("filter", filter)

	children, pageInfo, err := ListEntitiesOfType(api, entityType, params, options.OutputJSONResponse, options.Out, options.Timeout, options.Verbose)
	if err != nil {
		return err
	}

	options.Printf("filter returned ")
	pageInfo.Output(options)

	var ids []string
	for _, entity := range children {
		id, _ := entity.Path("id").Data().(string)
		ids = append(ids, id)
	}

	return DeleteEntitiesOfType(api, options, entityType, ids, body)
}

func GetPlural(entityType string) string {
	if strings.HasSuffix(entityType, "y") {
		return strings.TrimSuffix(entityType, "y") + "ies"
	}
	return entityType + "s"
}