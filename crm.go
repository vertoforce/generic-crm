package crm

import (
	"github.com/vertoforce/generic-crm/backends/airtablecrm"
	"github.com/vertoforce/generic-crm/backends/googlesheet"

	"github.com/vertoforce/generic-crm/backends/crm"
)

var backends = []crm.CRM{
	&googlesheet.Client{},
	&airtablecrm.Client{},
}
