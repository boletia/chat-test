package bot

import log "github.com/sirupsen/logrus"

func (b bot) ReadMessage() {
	for {
		select {
		case <-b.quit:
			return
		default:
			msgType, msg, err := b.conn.ReadMessage()

			log.WithFields(log.Fields{
				"bot":   b.NickName,
				"error": err,
				"type":  msgType,
				"msg":   string(msg),
			}).Info("message received")
		}
	}
}
