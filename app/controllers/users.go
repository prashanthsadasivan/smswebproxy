package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"smswebproxy/app/appconnections"
	"smswebproxy/app/models"
)

func (c App) Create(regid, num string) revel.Result {
	fmt.Printf("regId: %s\n num:%s\n postnum: %s\n", regid, num, c.Request.PostFormValue("num"))
	appconnections.New(regid, num)
	results, err := c.Txn.Select(models.SoleUser{}, "select * from SoleUser")
	if err != nil {
		panic(err)
	}
	fmt.Printf("results: %s\n", results)
	if len(results) >= 1 {
		numDeleted, err := c.Txn.Delete(results...)
		if err != nil {
			fmt.Printf("err: %s\n", err)
		}
		fmt.Printf("numDeleted: %d\n", numDeleted)
	}
	newSoleUser := &models.SoleUser{Number: num, GcmId: regid}
	insertErr := c.Txn.Insert(newSoleUser)
	if insertErr != nil {
		fmt.Printf("here err: %s\n", insertErr)
	}
	c.Response.Status = 201
	return c.RenderText("created")
}
