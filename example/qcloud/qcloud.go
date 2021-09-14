// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package main

import (
	"context"
	"github.com/ysicing/cloudsdk"
	"github.com/ysicing/cloudsdk/drivers/qcloud"
	"log"
	"os"
)

func setupProvider() (cloudsdk.Provider, error) {
	return qcloud.New(qcloud.WithApi(os.Getenv("qkey"), os.Getenv("qsecret")),
		qcloud.WithRegion("ap-nanjing"))
}

func main() {
	ctx := context.Background()
	p, err := setupProvider()
	if err != nil {
		log.Println("Invalid or missing hosting provider: ", err)
		return
	}
	//instance, err := p.Create(ctx, cloudsdk.InstanceCreateOpts{})
	//if err != nil {
	//	log.Println("create instance err: ", err)
	//	return
	//}
	//log.Println(instance.Name)
	if err := p.Destroy(ctx, &cloudsdk.Instance{ID: "ins-kz20g1e6"}); err != nil {
		log.Println(err)
		return
	}
	log.Println("done")
}
