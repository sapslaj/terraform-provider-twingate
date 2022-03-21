package twingate

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/twingate/go-graphql-client"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Groups are how users are authorized to access Resources. For more information, see Twingate's [documentation](https://docs.twingate.com/docs/groups).",
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		DeleteContext: resourceGroupDelete,
		UpdateContext: resourceGroupUpdate,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Group",
			},
			"active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the group is active or deactivated.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the group",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceGroupCreate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	groupName := resourceData.Get("name").(string)
	group, err := client.createGroup(ctx, graphql.String(groupName))

	if err != nil {
		return diag.FromErr(err)
	}

	resourceData.SetId(group.ID.(string))
	log.Printf("[INFO] Group %s created with id %s", groupName, resourceData.Id())

	return resourceGroupRead(ctx, resourceData, meta)
}

func resourceGroupUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	groupID := resourceData.Id()
	groupName := resourceData.Get("name").(string)

	if resourceData.HasChange("name") {
		log.Printf("[INFO] Updating group id %s", groupID)

		err := client.updateGroup(ctx, graphql.ID(groupID), graphql.String(groupName))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGroupRead(ctx, resourceData, meta)
}

func resourceGroupDelete(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	var diags diag.Diagnostics

	groupID := resourceData.Id()

	err := client.deleteGroup(ctx, groupID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleted group id %s", resourceData.Id())

	return diags
}

func resourceGroupRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	var diags diag.Diagnostics

	groupID := resourceData.Id()
	group, err := client.readGroup(ctx, groupID)

	if err != nil {
		return diag.FromErr(err)
	}

	err = resourceData.Set("name", group.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = resourceData.Set("active", group.IsActive)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
