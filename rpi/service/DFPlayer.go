package service

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

const (
	DFPLAYER_EQ_NORMAL  = 0
	DFPLAYER_EQ_POP     = 1
	DFPLAYER_EQ_ROCK    = 2
	DFPLAYER_EQ_JAZZ    = 3
	DFPLAYER_EQ_CLASSIC = 4
	DFPLAYER_EQ_BASS    = 5

	DFPLAYER_DEVICE_U_DISK = 1
	DFPLAYER_DEVICE_SD     = 2
	DFPLAYER_DEVICE_AUX    = 3
	DFPLAYER_DEVICE_SLEEP  = 4
	DFPLAYER_DEVICE_FLASH  = 5

	DFPLAYER_RECEIVED_LENGTH = 10
	DFPLAYER_SEND_LENGTH     = 10

	//_DEBUG

	// message types
	TimeOut               = 0
	WrongStack            = 1
	DFPlayerCardInserted  = 2
	DFPlayerCardRemoved   = 3
	DFPlayerCardOnline    = 4
	DFPlayerPlayFinished  = 5
	DFPlayerError         = 6
	DFPlayerUSBInserted   = 7
	DFPlayerUSBRemoved    = 8
	DFPlayerUSBOnline     = 9
	DFPlayerCardUSBOnline = 10
	DFPlayerFeedBack      = 11

	Busy             = 1
	Sleeping         = 2
	SerialWrongStack = 3
	CheckSumNotMatch = 4
	FileIndexOut     = 5
	FileMismatch     = 6
	Advertise        = 7

	// format
	Stack_Header    = 0
	Stack_Version   = 1
	Stack_Length    = 2
	Stack_Command   = 3
	Stack_ACK       = 4
	Stack_Parameter = 5
	Stack_CheckSum  = 7
	Stack_End       = 9
)

const (
	_timeOutDuration = 500
)

var config = &serial.Config{Name: "/dev/ttyS0", Baud: 9600}

// DFPlayerMini :
type DFPlayerMini struct {
	CommandHandler func()
	port           *serial.Port
	alive          bool

	_timeOutTimer    uint64
	_timeOutDuration uint64
	_isAvailable     bool
	_isSending       bool

	_sending         [DFPLAYER_SEND_LENGTH]byte
	_received        [DFPLAYER_RECEIVED_LENGTH]byte
	_receivedIndex   uint8
	_handleType      uint8
	_handleCommand   uint8
	_handleParameter uint16
}

func NewDFPlayer() *DFPlayerMini {
	df := DFPlayerMini{}
	df._sending = [DFPLAYER_RECEIVED_LENGTH]byte{0x7E, 0xFF, 06, 00, 01, 00, 00, 00, 00, 0xEF}
	df._timeOutDuration = 10000
	return &df
}

func (df *DFPlayerMini) Connect(server *Server) {
	s, err := serial.OpenPort(config)
	df.port = s
	if err != nil {
		log.Fatal(err)
	}
	df.alive = true
	go df.readLoop(server.stopChannel)
}
func (df *DFPlayerMini) Close() {
	df.port.Close()
	df.alive = false
}
func (df *DFPlayerMini) readLoop(stopChannel chan bool) {
	buf := make([]byte, 128)
	for {
		n, err := df.port.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		df.readMsg(buf[:n])
		//log.Println(buf[:n])
		if df.CommandHandler != nil {
			df.CommandHandler()
		}
	}
}
func (df *DFPlayerMini) readMsg(buf []byte) bool {
	for _, b := range buf {
		if df._receivedIndex == 0 {
			df._received[Stack_Header] = b
			if df._received[Stack_Header] == 0x7E {
				df._receivedIndex++
			}
		} else {
			df._received[df._receivedIndex] = b
			switch df._receivedIndex {
			case Stack_Version:
				if df._received[df._receivedIndex] != 0xFF {
					return df.handleError(WrongStack, 0)
				}
			case Stack_Length:
				if df._received[df._receivedIndex] != 0x06 {
					return df.handleError(WrongStack, 0)
				}
			case Stack_End:
				if df._received[df._receivedIndex] != 0xEF {
					return df.handleError(WrongStack, 0)
				}
				if df.validateStack() {
					df._receivedIndex = 0
					df.parseStack()
					return df._isAvailable
				}
				return df.handleError(WrongStack, 0)
			default:
			}
			df._receivedIndex++
		}
	}
	return df._isAvailable
}
func (df *DFPlayerMini) Write(b []byte) {
	n, err := df.port.Write(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("write:", n)
	logArray(b[:n])
}
func (df *DFPlayerMini) sendStack(command uint8, arg uint16) {
	df._sending[Stack_Command] = command
	uint16ToArray(arg, df._sending[Stack_Parameter:])
	uint16ToArray(calculateCheckSum(&df._sending), df._sending[Stack_CheckSum:])

	df.Write(df._sending[:])
	msgtype := getType(command, arg)
	fmt.Printf("sendStack: cmd: %x, type: %d, params: %x\n", command, msgtype, arg)

	df._timeOutTimer = millis()
	df._isSending = df._sending[Stack_ACK] == 1 && true || false
	if df._sending[Stack_ACK] == 0 {
		delay(10)
	}
}

func (df *DFPlayerMini) parseStack() {
	fmt.Println("parseStack")

	handleCommand := df._received[Stack_Command]
	if handleCommand == 0x41 {
		df._isSending = false
		return
	}
	df._handleCommand = handleCommand
	df._handleParameter = arrayToUint16(df._received[Stack_Parameter:])

	msgtype := getType(df._handleCommand, df._handleParameter)
	if msgtype == WrongStack {
		df.handleError(msgtype, 0)
	} else {
		df.handleMessage(msgtype, df._handleParameter)
	}

}

func (df *DFPlayerMini) waitAvailable(duration uint64) bool {
	timer := millis()
	if duration == 0 {
		duration = df._timeOutDuration
	}
	for !df._isAvailable {
		fmt.Println("wait available:", df._receivedIndex, df._received)
		if millis()-timer > duration {
			return false
		}
		delay(1000)
	}
	return true
}

/* ------------------------- */
// helper

func logArray(array []byte) {
	for _, v := range array {
		fmt.Printf("%X ", v)
	}
}

func millis() uint64 {
	return uint64(time.Nanosecond * time.Duration(time.Now().UnixNano()) / time.Millisecond)
}
func delay(ms uint) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func uint16ToArray(v uint16, array []byte) {
	array[0] = byte(v >> 8)
	array[1] = byte(v)
}
func arrayToUint16(array []byte) uint16 {
	return uint16(array[0])<<8 + uint16(array[1])
}

func calculateCheckSum(buffer *[DFPLAYER_SEND_LENGTH]byte) uint16 {
	var sum uint16 = 0
	for i := Stack_Version; i < Stack_CheckSum; i++ {
		sum += uint16(buffer[i])
	}
	return -sum
}

func getParam(a1 uint8, a2 uint8) uint16 {
	return uint16(a1)<<8 | uint16(a2)
}

func getType(cmd uint8, params uint16) uint8 {
	switch cmd {
	case 0x3D:
		return DFPlayerPlayFinished
	case 0x3F:
		if params&0x01 > 0 {
			return DFPlayerUSBOnline
		} else if params&0x02 > 0 {
			return DFPlayerCardOnline
		} else if params&0x03 > 0 {
			return DFPlayerCardUSBOnline
		}
	case 0x3A:
		if params&0x01 > 0 {
			return DFPlayerUSBInserted
		} else if params&0x02 > 0 {
			return DFPlayerCardInserted
		}
	case 0x3B:
		if params&0x01 > 0 {
			return DFPlayerUSBRemoved
		} else if params&0x02 > 0 {
			return DFPlayerCardRemoved
		}
	case 0x40:
		return DFPlayerError
	case 0x3C:
		fallthrough
	case 0x3E:
		fallthrough
	case 0x42:
		fallthrough
	case 0x43:
		fallthrough
	case 0x44:
		fallthrough
	case 0x45:
		fallthrough
	case 0x46:
		fallthrough
	case 0x47:
		fallthrough
	case 0x48:
		fallthrough
	case 0x49:
		fallthrough
	case 0x4B:
		fallthrough
	case 0x4C:
		fallthrough
	case 0x4D:
		fallthrough
	case 0x4E:
		fallthrough
	case 0x4F:
		return DFPlayerFeedBack
	}
	return WrongStack
}

/* ------------------------- */

func (df *DFPlayerMini) setTimeOut(timeOutDuration uint64) {
	df._timeOutDuration = timeOutDuration
}

func (df *DFPlayerMini) enableACK() {
	df._sending[Stack_ACK] = 0x01
}
func (df *DFPlayerMini) disableACK() {
	df._sending[Stack_ACK] = 0x00
}

/* ------------------------- */

func (df *DFPlayerMini) validateStack() bool {
	return calculateCheckSum(&df._received) == arrayToUint16(df._received[Stack_CheckSum:])
}
func (df *DFPlayerMini) read() uint16 {
	df._isAvailable = false
	return df._handleParameter
}
func (df *DFPlayerMini) readType() uint8 {
	df._isAvailable = false
	return df._handleType
}
func (df *DFPlayerMini) readCommand() uint8 {
	df._isAvailable = false
	return df._handleCommand
}
func (df *DFPlayerMini) handleMessage(msgtype uint8, parameter uint16) bool {
	df._receivedIndex = 0
	df._handleType = msgtype
	df._handleParameter = parameter
	df._isAvailable = true

	fmt.Printf("handleMessage: cmd: %x, type: %d, params: %x\n", df._handleCommand, msgtype, parameter)
	fmt.Print("handleMessage-received: ")
	logArray(df._received[:])

	cs := [10]byte{}
	uint16ToArray(calculateCheckSum(&df._received), cs[Stack_CheckSum:])

	return df._isAvailable
}
func (df *DFPlayerMini) handleError(msgtype uint8, parameter uint16) bool {
	df.handleMessage(msgtype, parameter)
	df._isSending = false
	return false
}

// sendstack, waitAvailable, readType, read
func (df *DFPlayerMini) next() {
	df.sendStack(0x01, 0)
}

func (df *DFPlayerMini) previous() {
	df.sendStack(0x02, 0)
}

func (df *DFPlayerMini) play(fileNumber uint16) {
	df.sendStack(0x03, fileNumber)
}

func (df *DFPlayerMini) volumeUp() {
	df.sendStack(0x04, 0)
}

func (df *DFPlayerMini) volumeDown() {
	df.sendStack(0x05, 0)
}

func (df *DFPlayerMini) volume(volume uint16) {
	df.sendStack(0x06, volume)
}

func (df *DFPlayerMini) EQ(eq uint16) {
	df.sendStack(0x07, eq)
}

func (df *DFPlayerMini) loop(fileNumber uint16) {
	df.sendStack(0x08, fileNumber)
}

func (df *DFPlayerMini) outputDevice(device uint16) {
	df.sendStack(0x09, device)
	delay(200)
}

func (df *DFPlayerMini) sleep() {
	df.sendStack(0x0A, 0)
}

func (df *DFPlayerMini) reset() {
	df.sendStack(0x0C, 0)
}

func (df *DFPlayerMini) start() {
	df.sendStack(0x0D, 0)
}

func (df *DFPlayerMini) pause() {
	df.sendStack(0x0E, 0)
}

func (df *DFPlayerMini) playFolder(folderNumber uint8, fileNumber uint8) {
	df.sendStack(0x0F, getParam(folderNumber, fileNumber))
}

func (df *DFPlayerMini) outputSetting(enable bool, gain uint8) {
	var e uint8
	if enable {
		e = 1
	} else {
		e = 0
	}
	df.sendStack(0x10, getParam(e, gain))
}

func (df *DFPlayerMini) enableLoopAll() {
	df.sendStack(0x11, 0x01)
}

func (df *DFPlayerMini) disableLoopAll() {
	df.sendStack(0x11, 0x00)
}

func (df *DFPlayerMini) playMp3Folder(fileNumber uint16) {
	df.sendStack(0x12, fileNumber)
}

func (df *DFPlayerMini) advertise(fileNumber uint16) {
	df.sendStack(0x13, fileNumber)
}

func (df *DFPlayerMini) playLargeFolder(folderNumber uint8, fileNumber uint16) {
	df.sendStack(0x14, (uint16(folderNumber)<<12)|fileNumber)
}

func (df *DFPlayerMini) stopAdvertise() {
	df.sendStack(0x15, 0)
}

func (df *DFPlayerMini) stop() {
	df.sendStack(0x16, 0)
}

func (df *DFPlayerMini) loopFolder(folderNumber uint16) {
	df.sendStack(0x17, folderNumber)
}

func (df *DFPlayerMini) randomAll() {
	df.sendStack(0x18, 0)
}

func (df *DFPlayerMini) enableLoop() {
	df.sendStack(0x19, 0x00)
}

func (df *DFPlayerMini) disableLoop() {
	df.sendStack(0x19, 0x01)
}

func (df *DFPlayerMini) enableDAC() {
	df.sendStack(0x1A, 0x00)
}

func (df *DFPlayerMini) disableDAC() {
	df.sendStack(0x1A, 0x01)
}

func (df *DFPlayerMini) readState() int {
	df.sendStack(0x42, 0)
	if df.waitAvailable(0) {
		if df.readType() == DFPlayerFeedBack {
			return int(df.read())
		}
		return -1
	}
	return -1
}

func (df *DFPlayerMini) readVolume() int {
	df.sendStack(0x43, 0)
	if df.waitAvailable(0) {
		return int(df.read())
	}
	return -1
}

func (df *DFPlayerMini) readEQ() int {
	df.sendStack(0x44, 0)
	if df.waitAvailable(0) {
		if df.readType() == DFPlayerFeedBack {
			return int(df.read())
		}
		return -1
	}
	return -1
}
func (df *DFPlayerMini) readPlaybackMode() int {
	df.sendStack(0x45, 0)
	if df.waitAvailable(0) {
		if df.readType() == DFPlayerFeedBack {
			return int(df.read())
		}
		return -1
	}
	return -1
}

func (df *DFPlayerMini) readCurrentFileNumber(device uint8) int {
	if device == 0 {
		device = DFPLAYER_DEVICE_SD
	}
	switch device {
	case DFPLAYER_DEVICE_U_DISK:
		df.sendStack(0x4B, 0)
		break
	case DFPLAYER_DEVICE_SD:
		df.sendStack(0x4C, 0)
		break
	case DFPLAYER_DEVICE_FLASH:
		df.sendStack(0x4D, 0)
		break
	default:
		break
	}
	if df.waitAvailable(0) {
		if df.readType() == DFPlayerFeedBack {
			return int(df.read())
		}
		return -1
	}
	return -1
}

func (df *DFPlayerMini) readFileCountsInFolder(folderNumber uint16) int {
	df.sendStack(0x4E, folderNumber)
	if df.waitAvailable(0) {
		if df.readType() == DFPlayerFeedBack {
			return int(df.read())
		}
		return -1
	}
	return -1
}

func (df *DFPlayerMini) readFolderCounts() int {
	df.sendStack(0x4F, 0)
	if df.waitAvailable(0) {
		fmt.Println("foldercount:", df.read())
		if df.readType() == DFPlayerFeedBack {
			return int(df.read())
		}
		return -1
	}
	return -1
}

func (df *DFPlayerMini) readFileCounts(device uint8) int {
	if device == 0 {
		device = DFPLAYER_DEVICE_SD
	}
	switch device {
	case DFPLAYER_DEVICE_U_DISK:
		df.sendStack(0x47, 0)
		break
	case DFPLAYER_DEVICE_SD:
		df.sendStack(0x48, 0)
		break
	case DFPLAYER_DEVICE_FLASH:
		df.sendStack(0x49, 0)
		break
	default:
		break
	}

	if df.waitAvailable(0) {
		if df.readType() == DFPlayerFeedBack {
			return int(df.read())
		}
		return -1
	}
	return -1
}
