# CRM

[![Go Report Card](https://goreportcard.com/badge/github.com/vertoforce/generic-crm)](https://goreportcard.com/report/github.com/vertoforce/generic-crm)
[![Documentation](https://godoc.org/github.com/vertoforce/generic-crm?status.svg)](https://godoc.org/github.com/vertoforce/generic-crm)

This library allows you to

* Add/Update items
* Remove items
* Get items

From a generic interface supporting the following backends:

* airtable
* google sheet

You can use it to store and keep track of generic items without worrying about the implementation details, or easily swap out google sheet for airtable, etc.

## Usage

First set up the CRM

* Airtable: Set up your airtable such that the columns contain the names of the fields you will use in your CRM.
* Google Sheets: Set the first row of columns to the names of the fields you will use in your CRM.

Then you can create a client and use it.  There is an example in the godoc.
