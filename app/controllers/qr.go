package controllers

import (
	"bytes"
	"code.google.com/p/rsc/qr"
	"github.com/revel/revel"
	"time"
)

func (c App) QR(hostname string) revel.Result {
	code, err := qr.Encode(hostname, qr.H)
	if err != nil {
		panic(err)
	}
	png := code.PNG()
	modtime := time.Now()
	c.Response.Out.Header().Set("Content-Type", "image/png")
	return &revel.BinaryResult{
		Reader:  bytes.NewReader(png),
		Name:    "qr.png",
		ModTime: modtime,
	}
}
