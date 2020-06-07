package crm

import (
	"crm/backends/backend"
	"crm/backends/googlesheet"
)

var backends = []backend.Backend{
	&googlesheet.Client{},
}
