package internal

import (
	"bytes"
	"github.com/zerocruft/flux/debug"
)

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

func propogateMsg(token string, msgBytes []byte) {

	msg, ok := parseFluxMsg(msgBytes)
	if !ok {
		debug.Log("err in parsing message from: "+token)
	}

	switch msg.msgType {

	case FLUX_TOPIC_SUBSCRIBE:
		if msg.token != "0" && msg.glance != "0" {
			addSubToTopic(token, msg.glance)
		}
		debug.Log("Topic subscription [" + msg.glance + "] from: " + token)
		return

	case FLUX_MESSAGE_TEXT:
		if msg.glance != "" {
			subscribers := getCopyOfSubsForTopic(msg.glance)
			for _, subscriber := range subscribers {
				debug.Log("Topic Distribute - Topic[" + msg.glance + "] - Sub[" + subscriber + "]")
				go sendMsgToClient(subscriber, msgBytes)
			}
		}
		return
	default:
		debug.Log("Invalid Flux msg type: " + msg.msgType)
		return
	}
}

func parseFluxMsg(msgBytes []byte) (fluxMsg, bool) {
	msgSections := bytes.Split(msgBytes, []byte(":"))
	if len(msgSections) != 4 {
		return fluxMsg{}, false
	}
	msgType := string(msgSections[0])
	msg := fluxMsg{
		msgType: msgType,
		token:   string(msgSections[1]),
		glance:  string(msgSections[2]),
		payload: msgSections[3],
	}

	return msg, true
}

type fluxMsg struct {
	msgType string
	token   string
	glance  string
	payload []byte
}
