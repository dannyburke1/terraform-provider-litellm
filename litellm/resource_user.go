package litellm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	endpointUserNew    = "/user/new"
	endpointUserInfo   = "/user/info"
	endpointUserUpdate = "/user/update"
	endpointUserDelete = "/user/delete"
)

func ResourceLiteLLMUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiteLLMUserCreate,
		Read:   resourceLiteLLMUserRead,
		Update: resourceLiteLLMUserUpdate,
		Delete: resourceLiteLLMUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_role": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"proxy_admin",
					"proxy_admin_viewer",
					"internal_user",
					"internal_user_viewer",
					"team",
					"customer",
				}, false),
			},
			"max_budget": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"models": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tpm_limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rpm_limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_create_key": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"send_user_invite": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"teams": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceLiteLLMUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	userID := uuid.New().String()
	userData := buildUserData(d, userID)

	log.Printf("[DEBUG] Create user request payload: %+v", userData)

	resp, err := MakeRequest(client, "POST", endpointUserNew, userData)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	defer resp.Body.Close()

	if err := handleResponse(resp, "creating user"); err != nil {
		return err
	}

	d.SetId(userID)
	log.Printf("[INFO] User created with ID: %s", userID)

	/*
		There is a bug when creating users that the email invitation doesn't contain the URL correctly.
		I believe this could be due to a race condition and DB checks that happen on user creation.
	*/
	time.Sleep(2 * time.Second)

	return resourceLiteLLMUserRead(d, m)
}

func resourceLiteLLMUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	log.Printf("[INFO] Reading user with ID: %s", d.Id())

	resp, err := MakeRequest(client, "GET", fmt.Sprintf("%s?user_id=%s", endpointUserInfo, d.Id()), nil)
	if err != nil {
		return fmt.Errorf("error reading user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] user with ID %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	var userResp UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return fmt.Errorf("error decoding user info response: %w", err)
	}

	// Set API response so we can import users into the state.
	d.Set("user_email", userResp.userEmail)
	d.Set("user_alias", userResp.userAlias)
	d.Set("user_role", userResp.userRole)
	d.Set("max_budget", userResp.MaxBudget)
	d.Set("models", userResp.Models)
	d.Set("tpm_limit", userResp.TPMLimit)
	d.Set("rpm_limit", userResp.RPMLimit)
	d.Set("auto_create_key", userResp.autoCreateKey)
	d.Set("send_user_invite", userResp.sendEmailInvite)
	d.Set("teams", userResp.Teams)

	return nil
}

func resourceLiteLLMUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	userData := buildUserData(d, d.Id())
	log.Printf("[DEBUG] Update user request payload: %+v", userData)

	resp, err := MakeRequest(client, "POST", endpointUserUpdate, userData)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	defer resp.Body.Close()

	if err := handleResponse(resp, "updating user"); err != nil {
		return err
	}

	log.Printf("[INFO] Successfully updated user with ID: %s", d.Id())
	return resourceLiteLLMUserRead(d, m)
}

func resourceLiteLLMUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	log.Printf("[INFO] Deleting user with ID: %s", d.Id())

	deleteData := map[string]interface{}{
		"user_ids": []string{d.Id()},
	}

	resp, err := MakeRequest(client, "POST", endpointUserDelete, deleteData)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	defer resp.Body.Close()

	if err := handleResponse(resp, "deleting user"); err != nil {
		return err
	}

	log.Printf("[INFO] Successfully deleted user with ID: %s", d.Id())
	d.SetId("")
	return nil
}

func buildUserData(d *schema.ResourceData, userID string) map[string]interface{} {
	userData := map[string]interface{}{
		"user_id":    userID,
		"user_alias": d.Get("user_alias").(string),
	}

	for _, key := range []string{"user_email", "user_alias", "key_alias", "user_role", "max_budget", "models", "tpm_limit", "rpm_limit", "auto_create_keys", "send_email_invite"} {
		if v, ok := d.GetOk(key); ok {
			userData[key] = v
		}
	}

	return userData
}
