package password_reset_message_text

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/encoding/protojson"

	textpb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/text"

	"github.com/zitadel/terraform-provider-zitadel/gen/github.com/zitadel/zitadel/pkg/grpc/text"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

const (
	orgIDVar    = "org_id"
	languageVar = "language"
)

var (
	_ resource.Resource = &passwordResetMessageTextResource{}
)

func NewPasswordResetMessageText() resource.Resource {
	return &passwordResetMessageTextResource{}
}

type passwordResetMessageTextResource struct {
	clientInfo *helper.ClientInfo
}

func (r *passwordResetMessageTextResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_reset_message_text"
}

func (r *passwordResetMessageTextResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return text.GenSchemaMessageCustomText(ctx)
}

func (r *passwordResetMessageTextResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clientInfo = req.ProviderData.(*helper.ClientInfo)
}

func (r *passwordResetMessageTextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	orgID, language := getPlanAttrs(ctx, req.Plan, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := textpb.MessageCustomText{}
	resp.Diagnostics.Append(text.CopyMessageCustomTextFromTerraform(ctx, plan, &obj)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jsonpb := &runtime.JSONPb{
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
	data, err := jsonpb.Marshal(obj)
	if err != nil {
		resp.Diagnostics.AddError("failed to marshal", err.Error())
		return
	}
	zReq := &management.SetCustomPasswordResetMessageTextRequest{}
	if err := jsonpb.Unmarshal(data, zReq); err != nil {
		resp.Diagnostics.AddError("failed to unmarshal", err.Error())
		return
	}
	zReq.Language = language

	client, err := helper.GetManagementClient(r.clientInfo, orgID)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.SetCustomPasswordResetMessageText(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to create", err.Error())
		return
	}

	setID(plan, orgID, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *passwordResetMessageTextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state types.Object
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, language := getID(ctx, state)

	client, err := helper.GetManagementClient(r.clientInfo, orgID)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	zResp, err := client.GetCustomPasswordResetMessageText(ctx, &management.GetCustomPasswordResetMessageTextRequest{Language: language})
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}
	if zResp.CustomText.IsDefault {
		resp.Diagnostics.AddError("Error while reading", "No custom texts existing")
	}

	resp.Diagnostics.Append(text.CopyMessageCustomTextToTerraform(ctx, *zResp.CustomText, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setID(state, orgID, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *passwordResetMessageTextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	orgID, language := getPlanAttrs(ctx, req.Plan, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := textpb.MessageCustomText{}
	resp.Diagnostics.Append(text.CopyMessageCustomTextFromTerraform(ctx, plan, &obj)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jsonpb := &runtime.JSONPb{
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
	data, err := jsonpb.Marshal(obj)
	if err != nil {
		resp.Diagnostics.AddError("failed to marshal", err.Error())
		return
	}
	zReq := &management.SetCustomPasswordResetMessageTextRequest{}
	if err := jsonpb.Unmarshal(data, zReq); err != nil {
		resp.Diagnostics.AddError("failed to unmarshal", err.Error())
		return
	}
	zReq.Language = language

	client, err := helper.GetManagementClient(r.clientInfo, orgID)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.SetCustomPasswordResetMessageText(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to update", err.Error())
		return
	}

	setID(plan, orgID, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *passwordResetMessageTextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	orgID, language := getStateAttrs(ctx, req.State, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := helper.GetManagementClient(r.clientInfo, orgID)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.ResetCustomPasswordResetMessageTextToDefault(ctx, &management.ResetCustomPasswordResetMessageTextToDefaultRequest{Language: language})
	if err != nil {
		resp.Diagnostics.AddError("failed to delete", err.Error())
		return
	}
}

func setID(obj types.Object, orgID string, language string) {
	attrs := obj.Attributes()
	attrs["id"] = types.StringValue(orgID + "_" + language)
	attrs[orgIDVar] = types.StringValue(orgID)
	attrs[languageVar] = types.StringValue(language)
}

func getID(ctx context.Context, obj types.Object) (string, string) {
	id := helper.GetStringFromAttr(ctx, obj.Attributes(), "id")
	parts := strings.Split(id, "_")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return helper.GetStringFromAttr(ctx, obj.Attributes(), orgIDVar), helper.GetStringFromAttr(ctx, obj.Attributes(), languageVar)
}

func getPlanAttrs(ctx context.Context, plan tfsdk.Plan, diag diag.Diagnostics) (string, string) {
	var orgID string
	diag.Append(plan.GetAttribute(ctx, path.Root(orgIDVar), &orgID)...)
	if diag.HasError() {
		return "", ""
	}
	var language string
	diag.Append(plan.GetAttribute(ctx, path.Root(languageVar), &language)...)
	if diag.HasError() {
		return "", ""
	}

	return orgID, language
}

func getStateAttrs(ctx context.Context, state tfsdk.State, diag diag.Diagnostics) (string, string) {
	var orgID string
	diag.Append(state.GetAttribute(ctx, path.Root(orgIDVar), &orgID)...)
	if diag.HasError() {
		return "", ""
	}
	var language string
	diag.Append(state.GetAttribute(ctx, path.Root(languageVar), &language)...)
	if diag.HasError() {
		return "", ""
	}

	return orgID, language
}
