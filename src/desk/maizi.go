package desk

import (
	"lib/utils"
	"protocol"
	"github.com/golang/glog"
)

func (t *Desk) MaiZi(seat, count uint32) int32 {
	t.Lock() //房间加锁
	defer t.Unlock()
	glog.Infoln(seat,count)
	if !t.data.MaiZi {
		return int32(protocol.Error_BuyAlready)
	}
	if _, ok := t.ready[seat]; !ok {
		return int32(protocol.Error_BuyAlready)
	}

	if _, ok := t.maizi[seat]; ok {
		return int32(protocol.Error_BuyAlready)
	}

	t.maizi[seat] = count
	go func() {
		utils.Sleep(2) //延迟2秒
		t.Diceing()    //主动打骰
	}()

	stoc := &protocol.SMaiZi{
		Seat:&seat,
		Count:&count,
	}
	t.broadcast(stoc)
	glog.Infoln(stoc.String())
	return 0
}
