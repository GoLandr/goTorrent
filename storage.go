package main

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

type TorrentLocal struct { //local storage of the torrents for readd on server restart
	Hash          string
	DateAdded     string
	StoragePath   string
	TorrentName   string
	TorrentStatus string
	TorrentType   string //magnet or .torrent file
}

func readInTorrents(torrentStorage *bolt.DB) (TorrentLocalArray []*TorrentLocal) {

	TorrentLocalArray = []*TorrentLocal{}

	torrentStorage.View(func(tx *bolt.Tx) error {

		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			torrentLocal := new(TorrentLocal) //create a struct to store to an array
			var Dateadded []byte
			var StoragePath []byte
			var Hash []byte
			var TorrentName []byte
			var TorrentStatus []byte
			Dateadded = b.Get([]byte("Date"))
			if Dateadded == nil {
				fmt.Println("Date added error!")
				Dateadded = []byte(time.Now().Format("Jan _2 2006"))
			}
			StoragePath = b.Get([]byte("StoragePath"))
			if StoragePath == nil {
				fmt.Println("StoragePath error!")
				StoragePath = []byte("downloads")
			}
			Hash = b.Get([]byte("InfoHash"))
			if Hash == nil {
				fmt.Println("Hash error!")
			}
			TorrentName = b.Get([]byte("TorrentName"))
			if TorrentName == nil {
				fmt.Println("Torrent Name not found")
				TorrentName = []byte("Not Found!")
			}
			TorrentStatus = b.Get([]byte("TorrentStatus"))
			if TorrentStatus == nil {
				fmt.Println("Torrent Status not found in local storage")
				TorrentStatus = []byte("")
			}

			torrentLocal.DateAdded = string(Dateadded)
			torrentLocal.StoragePath = string(StoragePath)
			torrentLocal.Hash = string(Hash) //Converting the byte slice back into the full hash
			torrentLocal.TorrentName = string(TorrentName)
			torrentLocal.TorrentStatus = string(TorrentStatus)

			fmt.Println("Torrentlocal list: ", torrentLocal)
			TorrentLocalArray = append(TorrentLocalArray, torrentLocal) //dumping it into the array
			return nil
		})
		return nil
	})

	return TorrentLocalArray //all done, return the entire Array to add to the torrent client
}

func addTorrentLocalStorage(torrentStorage *bolt.DB, local *TorrentLocal) {
	println("Adding Local storage information")

	torrentStorage.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(local.Hash)) //translating hash into bytes for storage
		if err != nil {
			return fmt.Errorf("create bucket %s", err)
		}
		err = b.Put([]byte("Date"), []byte(local.DateAdded)) //TODO error checking  marshall into JSON
		if err != nil {
			return err
		}
		err = b.Put([]byte("StoragePath"), []byte(local.StoragePath))
		if err != nil {
			return err
		}
		err = b.Put([]byte("InfoHash"), []byte(local.Hash))
		if err != nil {
			return err
		}
		err = b.Put([]byte("TorrentName"), []byte(local.TorrentName))
		if err != nil {
			return err
		}

		return nil
	})
}