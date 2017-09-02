package internal

import (
	"bytes"
	"github.com/zerocruft/flux/capacitor"
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

func propogateMsg(token string, msg []byte) {
	msgSections := bytes.Split(msg, []byte(":"))
	if len(msgSections) == 4 {
		if string(msgSections[0]) == capacitor.FLUX_TYPE_TOPIC_SUBSCRIPTION {
			if string(msgSections[1]) != "0" && string(msgSections[2]) != "0" {
				addSubToTopic(token, string(msgSections[2]))
			}
			debug.Log("Topic subscription [" + string(msgSections[2]) + "] from: " + string(msgSections[1]))
			return
		}

		if string(msgSections[0]) == capacitor.FLUX_TYPE_STANDARD_MESSAGE {
			if string(msgSections[2]) != "" {
				subs := getCopyOfSubsForTopic(string(msgSections[2]))
				for _, s := range subs {
					debug.Log("Topic Distribute - Topic[" + string(msgSections[2]) + "] - Sub[" + s + "]")
					go sendMsgToClient(s, msg)
				}
			}
			return
		}

		if string(msgSections[0]) == capacitor.FLUX_TYPE_SYSTEM_PING {
			//TODO system ping
			return
		}
	}
}

func parseFluxMsg(msgBytes []byte) (FluxMsg, bool) {
	msgSections := bytes.Split(msgBytes, []byte(":"))
	if len(msgSections) != 4 {
		return FlxMsg{}, false
	}
}

type FluxMsg struct {

}
