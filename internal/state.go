package internal

import (
	"bytes"
	"github.com/zerocruft/capacitor"
	"github.com/zerocruft/flux/debug"
)

func NumberOfConnections() int {
	return len(fccs)
}

func persistClient(fcc *fluxClientConnection) {
	fccMutex.Lock()
	fccs[fcc.token] = fcc
	fccMutex.Unlock()
}

func killClient(token string) {
	fccMutex.Lock()
	defer fccMutex.Unlock()

	delete(fccs, token)
	deleteSubFromAllTopics(token)
}

func sendMsgToClient(token string, msg []byte) {
	fccMutex.Lock()
	client, exists := fccs[token]
	if !exists {
		fccMutex.Unlock()
		return
	}
	fccMutex.Unlock()
	client.sendToClient <- msg
}

//---------
// Topics

func addSubToTopic(sub, topic string) {
	topicsMutex.Lock()
	defer topicsMutex.Unlock()

	topicSubs, exists := topics[topic]
	if !exists {
		newTopicSubs := []string{sub}
		topics[topic] = newTopicSubs
		return
	}

	// Iterate over subscribers making sure to not duplicate subscriber for topic
	for _, s := range topicSubs {
		if s == sub {
			debug.Log("subsriber[" + sub + "] already subscribed to topic[" + topic + "]")
			return
		}
	}

	topicSubs = append(topicSubs, sub)
	topics[topic] = topicSubs
}

func getCopyOfSubsForTopic(topic string) []string {
	topicsMutex.Lock()
	defer topicsMutex.Unlock()

	return topics[topic]
}

func deleteSubFromAllTopics(sub string) {
	for t, _ := range topics {
		deleteSubFromTopic(t, sub)
	}
}
func deleteSubFromTopic(topic, sub string) {
	topicsMutex.Lock()
	defer topicsMutex.Unlock()

	subs := topics[topic]
	newSubs := []string{}
	for _, s := range subs {
		if s != sub {
			newSubs = append(newSubs, s)
		}
	}
	if len(newSubs) == 0 {
		debug.Log("Removing Topic: " + topic)
		delete(topics, topic)
		return
	}

	topics[topic] = newSubs
	return
}

//-----------
// Msg

func PropogateMsg(token string, msgBytes []byte) {

	msg, ok := parseFluxMsg(msgBytes)
	if !ok {
		debug.Log("err in parsing message from: " + token)
	}

	switch msg.Control {

	case CONTROL_TOPIC_SUBSCRIBE:
		if msg.Topic != "0" {
			addSubToTopic(token, msg.Topic)
		}
		debug.Log("Topic subscription [" + msg.Topic + "] from: " + token)
		return

	case CONTROL_MESSAGE_TEXT:
		distributeMsgToSubscribers(msg, msgBytes)
		return
	case CONTROL_PEER_CHAT:
		distributeMsgToSubscribers(msg, msgBytes)
		return
	default:
		debug.Log("Invalid Flux msg type: " + msg.Control)
		return
	}
}

func distributeMsgToSubscribers(msg capacitor.FluxMessage, msgBytes []byte) {
	if msg.Topic != "0" {
		subscribers := getCopyOfSubsForTopic(msg.Topic)
		for _, subscriber := range subscribers {
			debug.Log("Topic Distribute - Topic[" + msg.Topic + "] - Sub[" + subscriber + "]")
			go sendMsgToClient(subscriber, msgBytes)
		}
	}
}

func parseFluxMsg(msgBytes []byte) (capacitor.FluxMessage, bool) {
	msgSections := bytes.Split(msgBytes, []byte("::"))
	if len(msgSections) != 3 {
		return capacitor.FluxMessage{}, false
	}
	msgType := string(msgSections[0])
	msg := capacitor.FluxMessage{
		Control: msgType,
		Topic:   string(msgSections[1]),
		Payload: msgSections[2],
	}

	return msg, true
}
