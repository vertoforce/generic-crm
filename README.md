# CRM

[![Go Report Card](https://goreportcard.com/badge/github.com/vertoforce/generic-crm)](https://goreportcard.com/report/github.com/vertoforce/generic-crm)
[![Documentation](https://godoc.org/github.com/vertoforce/generic-crm?status.svg)](https://godoc.org/github.com/vertoforce/generic-crm)

This library acts as a generic interface to multiple backend "CRMs"

It allows you to get, create, update, and remove generic "items" (`map[string]interface{}`) from a single interface supporting the following backends:

* airtable
* google sheet
* MySQL

You can use it to store and keep track of generic items without worrying about the implementation details, or easily swap out google sheet for airtable, etc.

It also has functionality to synchronize the CRM based on a stream of new incoming items.  This is useful if you are crawling a website or have a stream of updates and you want to apply them to a CRM.

## Usage

First set up the CRM

* Airtable: Set up your airtable such that the columns contain the names of the fields you will use in your CRM.
* Google Sheets: Set the first row of columns to the names of the fields you will use in your CRM.
* MySQL: Set up the table to match the columns to the keys in the item map

Then you can create a client and use it.  There is an example in the [godoc](https://godoc.org/github.com/vertoforce/generic-crm).
