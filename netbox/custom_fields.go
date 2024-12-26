package netbox

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const customFieldsKey = "custom_fields"

var customFieldsSchema = &schema.Schema{
	Type:     schema.TypeMap,
	Optional: true,
	Computed: true,
	Default:  nil,
	Elem: &schema.Schema{
		Type:    schema.TypeString,
		Default: nil,
	},
}

var customFieldsSchemaRead = &schema.Schema{
	Type:     schema.TypeMap,
	Computed: true,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
}

func normalizeCustomFields(cfm map[string]interface{}) map[string]interface{} {
	newcfm := make(map[string]interface{})

	for k, v := range cfm {
		if v != nil && v != "" {
			newcfm[k] = v
		}
	}

	return newcfm
}

func mergeCustomFields(oldcfm, newcfm map[string]interface{}) map[string]interface{} {
	if newcfm == nil {
		newcfm = make(map[string]interface{})
	}

	for k, v := range newcfm {
		if v == nil {
			newcfm[k] = ""
		}
	}

	for k := range oldcfm {
		if _, ok := newcfm[k]; !ok {
			newcfm[k] = ""
		}
	}

	return newcfm
}

func customFieldsDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	cfg := d.GetRawConfig().GetAttr(customFieldsKey)
	cfm, ok := d.Get(customFieldsKey).(map[string]interface{})

	if cfg.IsNull() || !ok {
		d.SetNew(customFieldsKey, nil)
	} else {
		newcfm := normalizeCustomFields(cfm)
		d.SetNew(customFieldsKey, newcfm)
	}

	return nil
}

func computeCustomFieldsAttr(cf interface{}) map[string]interface{} {
	cfm, _ := cf.(map[string]interface{})
	return normalizeCustomFields(cfm)
}

func computeCustomFieldsModel(d *schema.ResourceData) interface{} {
	oldcf, newcf := d.GetChange(customFieldsKey)

	oldcfm, _ := oldcf.(map[string]interface{})
	newcfm, _ := newcf.(map[string]interface{})
	return mergeCustomFields(oldcfm, newcfm)
}

type CustomFieldParams struct {
	params runtime.ClientRequestWriter
	cfm    map[string]interface{}
}

func (o *CustomFieldParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {
	if err := o.params.WriteToRequest(r, reg); err != nil {
		return err
	}

	for k, v := range o.cfm {
		if vs, ok := v.(string); ok {
			if err := r.SetQueryParam(fmt.Sprintf("cf_%s", url.QueryEscape(k)), vs); err != nil {
				return err
			}
		}
	}

	return nil
}

func WithCustomFieldParamsOption(cfm map[string]interface{}) func(*runtime.ClientOperation) {
	if cfm == nil {
		cfm = make(map[string]interface{})
	}

	return func(co *runtime.ClientOperation) {
		co.Params = &CustomFieldParams{
			params: co.Params,
			cfm:    cfm,
		}
	}
}
