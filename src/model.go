package src

type DcPropertyDto struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	NamespaceId string `json:"namespaceId"`
	ServiceId   string `json:"serviceId"`
	Active      bool   `json:"active"`
}
