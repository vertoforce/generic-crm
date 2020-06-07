package crm

import (
	"crm/backends/airtable"
	"crm/backends/backend"
	"crm/backends/googlesheet"
)

var backends = []backend.Backend{
	&googlesheet.Client{},
	&airtable.Client{},
}
