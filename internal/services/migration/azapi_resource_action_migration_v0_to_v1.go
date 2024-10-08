package migration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AzapiResourceActionMigrationV0ToV1(ctx context.Context) resource.StateUpgrader {
	return resource.StateUpgrader{
		PriorSchema: &schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed: true,
				},

				"type": schema.StringAttribute{
					Required: true,
				},

				"resource_id": schema.StringAttribute{
					Required: true,
				},

				"action": schema.StringAttribute{
					Optional: true,
				},

				"method": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},

				"body": schema.StringAttribute{
					Optional: true},

				"when": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},

				"locks": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},

				"response_export_values": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},

				"output": schema.StringAttribute{
					Computed: true,
				},
			},

			Blocks: map[string]schema.Block{
				"timeouts": timeouts.Block(ctx, timeouts.Opts{
					Create: true,
					Update: true,
					Read:   true,
					Delete: true,
				}),
			},
			Version: 0,
		},
		StateUpgrader: func(ctx context.Context, request resource.UpgradeStateRequest, response *resource.UpgradeStateResponse) {
			type OldModel struct {
				ID                   types.String   `tfsdk:"id"`
				Type                 types.String   `tfsdk:"type"`
				ResourceId           types.String   `tfsdk:"resource_id"`
				Action               types.String   `tfsdk:"action"`
				Method               types.String   `tfsdk:"method"`
				Body                 types.String   `tfsdk:"body"`
				When                 types.String   `tfsdk:"when"`
				Locks                types.List     `tfsdk:"locks"`
				ResponseExportValues types.List     `tfsdk:"response_export_values"`
				Output               types.String   `tfsdk:"output"`
				Timeouts             timeouts.Value `tfsdk:"timeouts"`
			}
			type newModel struct {
				ID                   types.String   `tfsdk:"id"`
				Type                 types.String   `tfsdk:"type"`
				ResourceId           types.String   `tfsdk:"resource_id"`
				Action               types.String   `tfsdk:"action"`
				Method               types.String   `tfsdk:"method"`
				Body                 types.Dynamic  `tfsdk:"body"`
				When                 types.String   `tfsdk:"when"`
				Locks                types.List     `tfsdk:"locks"`
				ResponseExportValues types.List     `tfsdk:"response_export_values"`
				Output               types.Dynamic  `tfsdk:"output"`
				Timeouts             timeouts.Value `tfsdk:"timeouts"`
			}

			var oldState OldModel
			if response.Diagnostics.Append(request.State.Get(ctx, &oldState)...); response.Diagnostics.HasError() {
				return
			}

			when := oldState.When
			if when.IsNull() {
				when = types.StringValue("apply")
			}

			body := types.DynamicNull()
			if !oldState.Body.IsNull() && oldState.Body.ValueString() != "" {
				body = types.DynamicValue(types.StringValue(oldState.Body.ValueString()))
			}

			newState := newModel{
				ID:                   oldState.ID,
				Type:                 oldState.Type,
				ResourceId:           oldState.ResourceId,
				Action:               oldState.Action,
				Method:               oldState.Method,
				Body:                 body,
				When:                 when,
				Locks:                oldState.Locks,
				ResponseExportValues: oldState.ResponseExportValues,
				Output:               types.DynamicValue(types.StringValue(oldState.Output.ValueString())),
				Timeouts:             oldState.Timeouts,
			}

			response.Diagnostics.Append(response.State.Set(ctx, newState)...)
		},
	}
}
