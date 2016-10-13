package main

import (
	"github.com/willf/bloom"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"time"
	"io/ioutil"
	"log"
)

var u *Uniq

func today() int{
	return int(time.Now().Unix()) / 86400
}

func Open(workdir string, n uint, maxdays int) (u *Uniq, err error) {
	u = &Uniq{filters: make([]filter, maxdays), maxdays: maxdays, day: today(), workdir: workdir, n:n}
	for i:=0; i<maxdays; i++ {

		f := filter{path: fmt.Sprintf("%s/%d.blf" , workdir, u.day-i)}
		f.open(n)

		u.filters[i] = f
	}
	go u.sync_worker()
	return
}

type Uniq struct {
	filters []filter
	maxdays int
	n uint
	day int
	workdir string
}

func (u *Uniq)sync_worker(){

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigs:
			log.Println("exiting...")
			for i:=0; i<u.maxdays; i++ {
				u.filters[i].sync()
			}
			os.Exit(0)
		case <-time.NewTimer(time.Second * 300).C:
			for i:=0; i<u.maxdays; i++ {
				u.filters[i].sync()
			}

			//换天了
			if today() != u.day {
				u.day = today()
				f := filter{path: fmt.Sprintf("%s/%d.blf" , u.workdir, u.day)}
				f.open(u.n)

				u.filters = append([]filter{f},   u.filters[:1]...)
			}
		}
	}
}

func (u *Uniq) Write(test []byte, days int) {
	if days < u.maxdays {
		u.filters[days].Add(test)
		u.filters[days].is_synced = false
	}
}

func (u *Uniq) Test(test []byte, days int) (ok bool) {
	if days >= u.maxdays || days<=0{
		days = u.maxdays-1
	}
	for i:=0;i<=days;i++{
		if u.filters[i].Test(test) {
			return true
		}
	}
	return false
}

type filter struct {
	*bloom.BloomFilter
	path string
	is_synced bool
}

func (f *filter)open(n uint) (err error) {
	fd, err := os.OpenFile(f.path, os.O_CREATE | os.O_RDWR, 0644)
	if err != nil {
		log.Println( "openerror error", f.path, err)
		return err
	}

	defer fd.Close()
	buf, err := ioutil.ReadAll(fd)
	f.BloomFilter = bloom.New(n, 10)

	if err == nil {
		err = f.GobDecode(buf)
		if err != nil {
			log.Println( "decode error", f.path, err)
		}
	}else{
		log.Println( "read error", f.path, err)
	}
	f.is_synced = true
	return
}

func (f *filter) sync() {
	if f.is_synced == false{
		log.Println("writing "+f.path+"...")
		buf, err := f.GobEncode()
		if err == nil {
			//copy on write
			fd, err := os.OpenFile(f.path+".tmp", os.O_CREATE | os.O_RDWR, 0644)
			if err!= nil {
				log.Println("write error", f.path+".tmp", err)
			}
			defer  fd.Close()
			fd.Truncate(0)
			fd.Seek(0, 0)
			_, err = fd.Write(buf)
			if err == nil {
				f.is_synced = true
				fd.Sync()
				defer func(){
					os.Remove(f.path)
					os.Rename(f.path+".tmp", f.path)					
				}()
			}else{
				log.Println( "write error", f.path, err)
			}
		}else{
			log.Println( "encode error", f.path, err)
		}
	}
}
