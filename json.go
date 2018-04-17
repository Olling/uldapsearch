package main

import (
	"fmt"
	"encoding/json"
	"gopkg.in/ldap.v2"
)

type JsonEntry struct {
	DN string
	Attributes []*JsonEntryAttribute
}

type JsonEntryAttribute struct {
	Name string
	Values []string
}


func toJsonEntry (input ldap.Entry) (output JsonEntry){
	output.DN = input.DN
	output.Attributes = toJsonEntryAttributes(input.Attributes)
	return output
}

func toJsonEntryAttributes (attributes []*ldap.EntryAttribute)(output []*JsonEntryAttribute) {
	for _,attribute := range attributes {
		jsonattribute := toJsonEntryAttribute(*attribute)
		output = append(output,&jsonattribute)
	}
	return output
}

func toJsonEntryAttribute (input ldap.EntryAttribute)(output JsonEntryAttribute) {
	output.Name = input.Name
	output.Values = input.Values
	return output
}

func toJson(s interface{}) (string,error){
        bytes, marshalErr := json.MarshalIndent(s,"","  ")
        if marshalErr != nil {
		fmt.Println("Could not convert struct to bytes", marshalErr)
                return "",marshalErr
        }
        return string(bytes),nil
}
