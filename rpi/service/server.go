package service

var Musics = map[int]map[int]string{
	1: map[int]string{
		1: "d",
	},
}
var Value uint32 = 0

// Server ;
type Server struct {
	stopChannel chan bool
	// いらない？
	gatt *GATTServer
	df   *DFPlayerMini
}

func NewServer() *Server {
	s := &Server{}
	s.stopChannel = make(chan bool)
	return s
}

func (s *Server) StartGATTServer() {
	gatt := &GATTServer{}
	gatt.Start(s)
	s.gatt = gatt
	if s.df != nil {
		s.df.Close()
	}
}
func (s *Server) StartUARTService() {
	df := NewDFPlayer()
	df.Connect(s)
	s.df = df
}

func Close() {

}

func (*Server) GetMusic() {
}
func (*Server) GetValue() uint32 {
	return Value
}
func (*Server) SetValue(v uint32) {
	Value = v
}
