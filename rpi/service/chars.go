package service

import (
	"fmt"
	"log"
	"time"

	"github.com/currantlabs/ble"
)

var (
	TestSvcUUID    = NewUUID(0x00010000)
	PlayCharUUID   = NewUUID(0x00020000)
	VolumeCharUUID = NewUUID(0x00030000)
	LoopCharUUID   = NewUUID(0x00040000)
	EQCharUUID     = NewUUID(0x00050000)
	MusicCharUUID  = NewUUID(0x00060000)
	StatusCharUUID = NewUUID(0x00070000)
	TestCharUUID   = NewUUID(0x00080000)
)

func NewUUID(n uint32) ble.UUID {
	return ble.MustParse(fmt.Sprintf("%08x-0011-1000-8000-00805F9B34FB", n))
}

// NewPlayChar ...
func NewPlayChar(s *Server) *ble.Characteristic {
	c := ble.NewCharacteristic(PlayCharUUID)

	c.HandleRead(ble.ReadHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		state := s.df.readState()
		fmt.Printf("playchar-read: state: %d\n", state)
		rsp.Write([]byte{byte(state)})
	}))

	c.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		log.Printf("playchar-write: %s", string(req.Data()))
		str := string(req.Data())
		switch str[0] {
		case 'n':
			s.df.next()
		case 'p':
			s.df.previous()
		case 'r':
			s.df.start()
		case 's':
			s.df.pause()
		case 'f':
			folder := uint8(str[1])
			file := uint8(str[2])
			s.df.playFolder(folder, file)
		case 'd':

		}
	}))
	return c
}

// NewVolumeChar ...
func NewVolumeChar(s *Server) *ble.Characteristic {
	c := ble.NewCharacteristic(VolumeCharUUID)

	c.HandleRead(ble.ReadHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		vol := s.df.readVolume()
		fmt.Printf("volume-read: vol: %d\n", vol)
		rsp.Write([]byte{byte(vol)})
	}))

	c.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		log.Printf("volume-write: %s", string(req.Data()))
		str := string(req.Data())
		switch str[0] {
		case 'u':
			s.df.volumeUp()
		case 'd':
			s.df.volumeDown()
		default:
			if len(str) >= 2 {
				vol := uint16(str[0]) << 8
				vol += uint16(str[1])
				s.df.volume(vol)
			}
		}
	}))
	return c
}

// NewLoopChar ...
func NewLoopChar(s *Server) *ble.Characteristic {
	c := ble.NewCharacteristic(LoopCharUUID)

	c.HandleRead(ble.ReadHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		state := s.df.readPlaybackMode()
		fmt.Printf("loop-read: vol: %d\n", state)
		rsp.Write([]byte{byte(state)})
	}))

	c.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		log.Printf("loop-write: %s", string(req.Data()))
		str := string(req.Data())
		switch str[0] {
		case 'e':
			s.df.enableLoop()
		case 'd':
			s.df.disableLoop()
		case 'r':
			s.df.randomAll()
		}
	}))
	return c
}

// NewEQChar ...
func NewEQChar(s *Server) *ble.Characteristic {
	c := ble.NewCharacteristic(EQCharUUID)

	c.HandleRead(ble.ReadHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		state := s.df.readEQ()
		fmt.Printf("eq-read: vol: %d\n", state)
		rsp.Write([]byte{byte(state)})
	}))

	c.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		log.Printf("loop-write: %s", string(req.Data()))
		str := string(req.Data())
		v := uint16(str[0]) << 8
		v += uint16(str[1])
		s.df.EQ(v)
	}))
	return c
}

// NewMusicChar ...
func NewMusicChar(s *Server) *ble.Characteristic {
	c := ble.NewCharacteristic(MusicCharUUID)
	ch := make(chan uint16, 10)

	c.HandleRead(ble.ReadHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		folders := byte(s.df.readFolderCounts())
		totalFiles := s.df.readFileCounts(0)
		fmt.Printf("music-read: vol: %d\n", folders)
		fmt.Printf("music-read: req: %d\n", totalFiles)

		if folders > 100 {
			folders = 100
		}
		b := make([]byte, folders)
		for i := byte(1); i <= folders; i++ {
			files := byte(s.df.readFileCountsInFolder(uint16(i)))
			b[i-1] = files
		}

		rsp.Write(b)
	}))
	c.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		log.Printf("music: Wrote %s", string(req.Data()))
		v := uint16(req.Data()[0]) << 8
		v += uint16(req.Data()[1])
		folder := v
		ch <- folder
	}))
	c.HandleNotify(ble.NotifyHandlerFunc(func(req ble.Request, n ble.Notifier) {
		log.Printf("test: Notification subscribed")
		for {
			select {
			case <-n.Context().Done():
				log.Printf("count: Notification unsubscribed")
				return
			case folder := <-ch:
				files := s.df.readFileCountsInFolder(folder)
				log.Printf("count: Notify files: %d", files)
				n.Write([]byte{byte(files)})
			}
		}
	}))
	return c
}

// NewCountChar ...
func NewTestChar(s *Server) *ble.Characteristic {
	n := 0
	c := ble.NewCharacteristic(TestCharUUID)

	c.HandleRead(ble.ReadHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		//fmt.Fprintf(rsp, "test: Read %d", n)
		count := 0

		//s.df.reset()
		//state := s.df.readState()
		//fmt.Printf("state: %d\n", state)

		//s.df.sendStack(0x46, 9)
		//fmt.Printf("state: %d\n", version)

		count = s.df.readState()
		fmt.Printf("state: %d\n", count)
		count = s.df.readCurrentFileNumber(0)
		fmt.Printf("num: %d\n", count)
		/*
			count = s.df.readFileCountsInFolder(1)
			fmt.Printf("files in 1: %d\n", count)
			count = s.df.readFileCountsInFolder(2)
			fmt.Printf("files in 2: %d\n", count)
		*/

		log.Printf("test: Read %d\n", count)
		rsp.Write([]byte{byte(count)})
		n++
	}))

	c.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		log.Printf("test: Wrote %s", string(req.Data()))
		s.SetValue(uint32(req.Data()[0]))

		str := string(req.Data()[0])
		if str == "u" {
			s.df.volumeUp()
		} else if str == "d" {
			s.df.volumeDown()
		} else {
			s.df.playFolder(1, 1)
		}
		//s.df.Write([]byte{0x7E, 0xFF, 0x06, 0x0F, 0x00, 0x01, 0x01, 0xef, 0xda, 0xEF}) //req.Data())
	}))

	c.HandleNotify(ble.NotifyHandlerFunc(func(req ble.Request, n ble.Notifier) {
		cnt := 0
		log.Printf("test: Notification subscribed")
		for {
			select {
			case <-n.Context().Done():
				log.Printf("count: Notification unsubscribed")
				return
			case <-time.After(time.Second):
				log.Printf("count: Notify: %d", cnt)
				if _, err := fmt.Fprintf(n, "Count: %d", cnt); err != nil {
					// Client disconnected prematurely before unsubscription.
					log.Printf("count: Failed to notify : %s", err)
					return
				}
				cnt++
			}
		}
	}))

	c.HandleIndicate(ble.NotifyHandlerFunc(func(req ble.Request, n ble.Notifier) {
		cnt := 0
		log.Printf("test: Indication subscribed")
		for {
			select {
			case <-n.Context().Done():
				log.Printf("test: Indication unsubscribed")
				return
			case <-time.After(time.Second):
				log.Printf("test: Indicate: %d", cnt)
				if _, err := fmt.Fprintf(n, "Count: %d", cnt); err != nil {
					// Client disconnected prematurely before unsubscription.
					log.Printf("test: Failed to indicate : %s", err)
					return
				}
				cnt++
			}
		}
	}))
	return c
}
