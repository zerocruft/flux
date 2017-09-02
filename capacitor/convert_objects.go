package capacitor

import (
	"bytes"
	"encoding/base64"
	"strings"
)

func fluxConnectionRequestToBytes() []byte {
	msg := FLUX_TYPE_CONNECTION_REQUEST + ":{NO-TOKEN}:{NO-DATA}:{NO-PAYLOAD}"
	msg = strings.Replace(msg, "{NO-TOKEN}", "0", 1)
	msg = strings.Replace(msg, "{NO-DATA}", "0", 1)
	msg = strings.Replace(msg, "{NO-PAYLOAD}", tobase64("0"), 1)
	return []byte(msg)
}

//TODO note: I don't want this exported, so I may have to maintain 2 identical snippets
func FluxConnectionResponseToBytes(clientToken string) []byte {
	msg := FLUX_TYPE_CONNECTION_RESPONSE + ":{TOKEN}:{NO-DATA}:{NO-PAYLOAD}"
	msg = strings.Replace(msg, "{TOKEN}", clientToken, 1)
	msg = strings.Replace(msg, "{NO-DATA}", "0", 1)
	msg = strings.Replace(msg, "{NO-PAYLOAD}", tobase64("0"), 1)
	return []byte(msg)
}

func fluxTopicSubscriptionRequestToBytes(clientToken, topic string) []byte {
	msg := FLUX_TYPE_TOPIC_SUBSCRIPTION + ":{TOKEN}:{TOPIC}:{NO-PAYLOAD}"
	msg = strings.Replace(msg, "{TOKEN}", clientToken, 1)
	msg = strings.Replace(msg, "{TOPIC}", topic, 1)
	msg = strings.Replace(msg, "{NO-PAYLOAD}", tobase64("0"), 1)
	return []byte(msg)
}

func fluxMessageToBytes(clientToken string, flxMsg FluxMessage) []byte {
	msg := FLUX_TYPE_STANDARD_MESSAGE + ":{TOKEN}:{TOPIC}:{PAYLOAD}"
	msg = strings.Replace(msg, "{TOKEN}", clientToken, -1)
	msg = strings.Replace(msg, "{TOPIC}", flxMsg.Topic, -1)
	msg = strings.Replace(msg, "{PAYLOAD}", tobase64(string(flxMsg.Payload)), -1)

	return []byte(msg)
}

func bytesToFluxMessage(msgBytes []byte) FluxMessage {
	sections := bytes.Split(msgBytes, []byte(":"))
	if len(sections) != 4 {
		// TODO throw an error or notify downstream somehow
		return FluxMessage{}
	}
	fluxMessage := FluxMessage{
		Topic:   string(sections[2]),
		Payload: frombase64(sections[3]),
	}
	return fluxMessage
}

func bytesToFluxObject(object []byte) (RawFluxObject, bool) {
	sections := bytes.Split(object, []byte(":"))
	if len(sections) != 4 {
		// TODO throw an error or notify downstream somehow
		return RawFluxObject{}, false
	}

	flxO := RawFluxObject{
		_type:        sections[0],
		_clientToken: sections[1],
		_data:        sections[2],
		_payload:     frombase64(sections[3]),
	}

	return flxO, true
}

func tobase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func frombase64(value []byte) (rv []byte) {
	rv = make([]byte, len(value))
	_, err := base64.StdEncoding.Decode(rv, value)
	if err != nil {
		return []byte{}
	}
	return
}
