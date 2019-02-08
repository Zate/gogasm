package main

import (
	"log"

	"github.com/dustin/seriesly/serieslyclient"
)

func initDB() {
	S, err := serieslyclient.New("http://127.0.0.1:3133")
	if err != nil {
		log.Fatalf("Something seriously fucked with the seriesly db: %v", err)
	}

	for _, r := range LiveRealms {
		err := S.Create(r)
		if err != nil {
			log.Fatalf("1 Something seriously fucked with the seriesly db: %v", err)
		}
		// SDBLive := S.DB(r)
		// dbinfo, err := SDBLive.Info()
		// if err != nil {
		// 	if strings.Contains(err.Error(), "no such file or directory") {
		// 		S.Create(r)
		// 	}
		// 	log.Fatalf("1 Something seriously fucked with the seriesly db: %v", err)
		// }

		log.Printf("DB created for %v %v", r, err)

	}

}
