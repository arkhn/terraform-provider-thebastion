package clients

type ResponseBastion struct {
	ErrorMessage string                 `json:"error_message"`
	ErrorCode    string                 `json:"error_code"`
	Command      string                 `json:"command"`
	Value        map[string]interface{} `json:"value"`
}

// Struct that hold information about singular user from accountList command
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

// Struct that hold information about singular user from accountListIngressKeys command
type AccountListIngressKeysValue struct {
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
}

type AccountListIngressKeys struct {
	Keys    []AccountListIngressKeysValue `json:"keys"`
	Account string                        `json:"account"`
}

type ResponseBastionListIngressKeys struct {
	ErrorMessage string                 `json:"error_message"`
	ErrorCode    string                 `json:"error_code"`
	Command      string                 `json:"command"`
	Value        AccountListIngressKeys `json:"value"`
}
