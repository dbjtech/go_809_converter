package handlers

import (
	"github.com/peifengll/go_809_converter/converter/handlers/converters"
	"github.com/tidwall/gjson"
	"log"
)

type distributor struct {
	converters map[string]converters.ConverterInterface
}
type Distributor interface {
	Handle(packet string)
	UplinkSend(res []byte)
	GetPacketType(item string) string
}

func NewDistributor(
	con map[string]converters.ConverterInterface) Distributor {
	return &distributor{
		converters: con,
	}
}

func (d *distributor) Handle(packet string) {
	packetType := d.GetPacketType(packet)
	converter, ok := d.converters[packetType]
	if !ok {
		log.Println("no converter for " + packetType)
		return
	}
	if value := gjson.Get(packet, "trace_id"); value.Exists() {
		converter.SetTraceID(value.String())
	}
	respPacket := converter.Handle(packet)
	d.UplinkSend(respPacket)
}

func (d *distributor) UplinkSend(res []byte) {
	if res == nil || len(res) == 0 {
		log.Println("no data need to send")
		return
	}
	if CsCenter.Uwriter == nil {
		log.Println("uplink disconnected")
	}
	// todo python里边是转化为16进制
	log.Println(string(res))
	CsCenter.Uwriter.Write(res)

}

func (d *distributor) GetPacketType(item string) string {
	s := gjson.Get(item, "packet_type").String()
	return s
}
