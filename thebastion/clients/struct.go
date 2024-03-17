package clients

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResponseBastion struct {
	ErrorMessage string      `json:"error_message"`
	ErrorCode    string      `json:"error_code"`
	Command      string      `json:"command"`
	Value        interface{} `json:"value"`
}

type Account struct {
	UID      int64  `json:"uid"`
	Name     string `json:"name"`
	IsActive int64  `json:"is_active"`
}

type ResponseBastionAccountList struct {
	ErrorMessage string             `json:"error_message"`
	ErrorCode    string             `json:"error_code"`
	Command      string             `json:"command"`
	Value        map[string]Account `json:"value"`
}

type ResponseBastionListIngressKeys struct {
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code"`
	Command      string `json:"command"`
	Value        struct {
		Keys []struct {
			Family      string        `json:"family"`
			Validity    string        `json:"validity"`
			Comment     string        `json:"comment"`
			Size        int64         `json:"size"`
			Id          int64         `json:"id"`
			Base64      string        `json:"base64"`
			Line        string        `json:"line"`
			Fingerprint string        `json:"fingerprint"`
			Prefix      string        `json:"prefix"`
			Typecode    string        `json:"typecode"`
			Mtime       interface{}   `json:"mtime"`
			From_list   []interface{} `json:"from_list"`
		} `json:"keys"`
		Account string `json:"account"`
	} `json:"value"`
}

// Struct that about info command on thebastion
type ResponseBastionInfo struct {
	Command      string `json:"command"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Value        struct {
		Account                 string   `json:"account"`
		AccountExpirationDays   string   `json:"account_expiration_days"`
		AdminAccounts           []string `json:"adminAccounts"`
		AllowedNetworksList     []any    `json:"allowed_networks_list"`
		BastionAliasCommand     string   `json:"bastion_alias_command"`
		BastionName             string   `json:"bastion_name"`
		EgressIPList            []string `json:"egress_ip_list"`
		EgressRsaMaxSize        string   `json:"egress_rsa_max_size"`
		EgressRsaMinSize        string   `json:"egress_rsa_min_size"`
		EgressSSHKeyAlgorithms  []string `json:"egress_ssh_key_algorithms"`
		ForbiddenNetworksList   []any    `json:"forbidden_networks_list"`
		Fortune                 string   `json:"fortune"`
		Hostname                string   `json:"hostname"`
		IdleKillTimeout         string   `json:"idle_kill_timeout"`
		IdleLockTimeout         string   `json:"idle_lock_timeout"`
		IngressKeysFromIPList   []any    `json:"ingress_keys_from_ip_list"`
		IngressRsaMaxSize       string   `json:"ingress_rsa_max_size"`
		IngressRsaMinSize       string   `json:"ingress_rsa_min_size"`
		IngressSSHKeyAlgorithms []string `json:"ingress_ssh_key_algorithms"`
		InteractiveModeAllowed  int      `json:"interactive_mode_allowed"`
		IsAdmin                 int      `json:"is_admin"`
		IsAuditor               int      `json:"is_auditor"`
		IsSuperowner            int      `json:"is_superowner"`
		MfaPasswordBypass       int      `json:"mfa_password_bypass"`
		MfaPasswordConfigured   int      `json:"mfa_password_configured"`
		MfaPasswordRequired     int      `json:"mfa_password_required"`
		MfaTotpBypass           int      `json:"mfa_totp_bypass"`
		MfaTotpConfigured       int      `json:"mfa_totp_configured"`
		MfaTotpRequired         int      `json:"mfa_totp_required"`
		MoshAllowed             int      `json:"mosh_allowed"`
		OsRelease               string   `json:"os_release"`
		OsSystem                string   `json:"os_system"`
		RegisteredAccounts      int      `json:"registered_accounts"`
		RegisteredGroups        int      `json:"registered_groups"`
		SlaveMode               int      `json:"slave_mode"`
		SuperOwnerAccounts      []any    `json:"superOwnerAccounts"`
		Uptime                  string   `json:"uptime"`
	} `json:"value"`
}

type ResponseBastionGroupList struct {
	Command      string `json:"command"`
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code"`
	Value        map[string]struct {
		Flags []string `json:"flags"`
	} `json:"value"`
}

type ResponseBastionCreateGroup struct {
	Command      string `json:"command"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Value        struct {
		Group     string `json:"group"`
		Owner     string `json:"owner"`
		PublicKey struct {
			Base64      string `json:"base64"`
			Comment     string `json:"comment"`
			Family      string `json:"family"`
			Filename    string `json:"filename"`
			Fingerprint string `json:"fingerprint"`
			FromList    []any  `json:"fromList"`
			Fullpath    string `json:"fullpath"`
			ID          string `json:"id"`
			Line        string `json:"line"`
			Mtime       int    `json:"mtime"`
			Prefix      string `json:"prefix"`
			Size        int    `json:"size"`
			Typecode    string `json:"typecode"`
		} `json:"public_key"`
	} `json:"value"`
}

type ResponseBastionListServer struct {
	Value []struct {
		IP            string `json:"ip"`
		Expiry        any    `json:"expiry"`
		AddedBy       string `json:"addedBy"`
		Port          string `json:"port"`
		ForcePassword any    `json:"forcePassword"`
		AddedDate     string `json:"addedDate"`
		ForceKey      any    `json:"forceKey"`
		ReverseDNS    any    `json:"reverseDns"`
		User          string `json:"user"`
		UserComment   string `json:"userComment"`
		Comment       string `json:"comment"`
	} `json:"value"`
	ErrorMessage string `json:"error_message"`
	ErrorCode    string `json:"error_code"`
	Command      string `json:"command"`
}

type ResponseBastionGroupInfo struct {
	Command      string `json:"command"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Value        struct {
		Guests_accesses struct{} `json:"guests_accesses"`
		Members         []string `json:"members"`
		Owners          []string `json:"owners"`
		FullMembers     []string `json:"full_members"`
		Inactive        []string `json:"inactive"`
		Group           string   `json:"group"`
		Gatekeepers     []string `json:"gatekeepers"`
		Keys            map[string]struct {
			Prefix      string   `json:"prefix"`
			Typecode    string   `json:"typecode"`
			Size        int64    `json:"size"`
			Mtime       int      `json:"mtime"`
			Id          string   `json:"id"`
			Filename    string   `json:"filename"`
			Fingerprint string   `json:"fingerprint"`
			Comment     string   `json:"comment"`
			FromList    []string `json:"fromList"`
			Base64      string   `json:"base64"`
			Fullpath    string   `json:"fullpath"`
			Line        string   `json:"line"`
			Family      string   `json:"family"`
		} `json:"keys"`
		Guests         []string `json:"guests"`
		Aclkeepers     []string `json:"aclkeepers"`
		PartialMembers []string `json:"partial_members"`
	} `json:"value"`
}

type ServerModel struct {
	Host        types.String `tfsdk:"host"`
	User        types.String `tfsdk:"user"`
	Port        types.Int64  `tfsdk:"port"`
	UserComment types.String `tfsdk:"user_comment"`
}
