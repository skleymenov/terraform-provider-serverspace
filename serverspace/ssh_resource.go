package serverspace

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/serverspace/ssclient"
)

func resourceSSH() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSHCreate,
		ReadContext:   resourceSSHRead,
		DeleteContext: resourceSSHDelete,
		Schema:        sshSchema,
	}
}

func resourceSSHCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	publicSSH := d.Get("public_key").(string)
	sshKey, err := client.CreateSSHKey(name, publicSSH)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(sshKey.ID))
	resourceSSHRead(ctx, d, m)
	return diags
}

func resourceSSHRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := m.(*ssclient.SSClient)
	sshID, _ := strconv.Atoi(d.Id())

	sshKey, err := client.GetSSHKey(sshID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", sshKey.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("public_key", sshKey.PublicKey); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(sshID))
	return diags
}

func resourceSSHDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	sshID, _ := strconv.Atoi(d.Id())

	err := client.DeleteSSHKey(sshID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
