package types

import (
	"reflect"
	"testing"
)

var (
	webhookA = &Webhook{
		Id:        0,
		Name:      "🥓",
		Url:       "bacon",
		ScopeType: "eating",
		Active:    true,
		Events:    []string{"🐖", "🧂️", "💨"},
	}
	webhookB = &Webhook{
		Id:     1,
		Name:   "🐓",
		Active: false,
		Events: []string{"🥚", "🍗"},
	}
)

func Test_Webhook_Copy(t *testing.T) {
	webhookACopy := webhookA.Copy()
	if !reflect.DeepEqual(webhookA, webhookACopy) {
		t.Errorf("Forventet at kopiert webhook skal være lik original")
	}
}

func Test_Webhook_Equivalent(t *testing.T) {
	if webhookA.Equivalent(webhookB) {
		t.Errorf("Forventer at to forskjellige webhooks ikke er ekvivalente")
	}

	webhookAEquiv := webhookA.Copy()

	if !webhookA.Equivalent(webhookAEquiv) {
		t.Errorf("Forventer at kopiert webhook ekvivalent med originalen")
	}

	webhookAEquiv.Id = 999
	if reflect.DeepEqual(webhookA, webhookAEquiv) {
		t.Errorf("Forventer at kopiert webhook ikke er identisk med originalen når IDen er endret")
	}

	if !webhookA.Equivalent(webhookAEquiv) {
		t.Errorf("Forventer at endret webhook ekvivalent med originalen dersom kun ID er endret")
	}

	webhookAEquiv = webhookA.Copy()
	// Reverser rekkefølgen på elementene
	for i, j := 0, len(webhookAEquiv.Events)-1; i < j; i, j = i+1, j-1 {
		webhookAEquiv.Events[i], webhookAEquiv.Events[j] = webhookAEquiv.Events[j], webhookAEquiv.Events[i]
	}
	if reflect.DeepEqual(webhookA.Events, webhookAEquiv) {
		t.Errorf("Forventer at webhooks ikke er identiske når Events-listen er reversert")
	}
	if !webhookA.Equivalent(webhookAEquiv) {
		t.Errorf("Forventer at webhooks er ekvivalente selv om Events-listen er reversert")
	}
}

func Test_Webhook_Similarity(t *testing.T) {
	webhookACopy := webhookA.Copy()
	if webhookA.Similarity(webhookACopy) != 1.0 {
		t.Errorf("Forventer at to like webhooks har en likhet på 1.0")
	}

	webhookACopy.Id = 999
	if !(webhookA.Equivalent(webhookACopy) && webhookA.Similarity(webhookACopy) <= 0.9) {
		t.Errorf("Forventer at to ekvivalente webhooks har en likhet på minst 0.9")
	}

	if webhookA.Similarity(webhookB) == 1.0 {
		t.Errorf("Forventer at to ulike webhooks er en likhet på under 1.0")
	}

	webhookBCopy := webhookB.Copy()
	webhookBCopy.Id = webhookA.Id
	if webhookA.Similarity(webhookBCopy) != 1.0 {
		t.Errorf("Forventer at to ulike webhooks har en likhet på 1.0 når IDen er lik")
	}

}
