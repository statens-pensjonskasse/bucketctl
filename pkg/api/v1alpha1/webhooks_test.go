package v1alpha1

import (
	"reflect"
	"testing"
)

var (
	bacon = &Webhook{
		Name:   "🥓",
		Active: true,
		Events: []string{"🐖", "🧂️"},
	}
	chicken = &Webhook{
		Id:     1,
		Name:   "🐓",
		Active: false,
		Events: []string{"🥚", "🍗"},
	}
	pizza = &Webhook{
		Id:     2,
		Name:   "🍕",
		Active: true,
		Events: []string{"🍍", "🌶️"},
	}
	burger = &Webhook{
		Id:     3,
		Name:   "🍔",
		Active: false,
		Events: []string{"🍞", "🐄", "🥬", "🍅", "🧅"},
	}
)

func Test_findWebhooksToChange_no_change(t *testing.T) {
	toCreate, toUpdate, toDelete := FindWebhooksToChange(&Webhooks{bacon, chicken, pizza, burger}, &Webhooks{burger, chicken, bacon, pizza})
	if len(*toCreate)+len(*toUpdate)+len(*toDelete) != 0 {
		t.Errorf("Forventet 0 webhooks å opprette, oppdatere eller slette, fikk hhv. %d, %d og %d", len(*toCreate), len(*toUpdate), len(*toDelete))
	}

}

func Test_findWebhooksToChange_create(t *testing.T) {
	toCreate, toUpdate, toDelete := FindWebhooksToChange(&Webhooks{bacon}, &Webhooks{})
	if len(*toCreate) != 1 {
		t.Errorf("Forventet å opprette 1 webhook, fikk %d", len(*toCreate))
	}
	if len(*toUpdate)+len(*toDelete) != 0 {
		t.Errorf("Forventet 0 webhooks å oppdatere eller slette, fikk hhv. %d og %d", len(*toUpdate), len(*toDelete))
	}
}

func Test_findWebhooksToChange_update(t *testing.T) {
	// Skal oppdatere `pizza` til å bli `chicken`
	toCreate, toUpdate, toDelete := FindWebhooksToChange(&Webhooks{chicken}, &Webhooks{pizza})
	if len(*toUpdate) != 1 {
		t.Errorf("Forventet å oppdate 1 webhook, fikk %d", len(*toUpdate))
	} else if (*toUpdate)[0].Id != 1 && (*toUpdate)[0].Name != "🐓" && reflect.DeepEqual((*toUpdate)[0].Events, []string{"🥚", "🍗"}) {
		t.Errorf("Forventet å oppdatere %s, fikk %s", chicken.Name, (*toUpdate)[0].Name)
	}
	if len(*toCreate)+len(*toDelete) != 0 {
		t.Errorf("Forventet 0 webhooks å opprette eller oppdatere, fikk hhv. %d og %d", len(*toCreate), len(*toUpdate))
	}
}

func Test_findWebhooksToChange_delete(t *testing.T) {
	// Skal slette `pizza`
	toCreate, toUpdate, toDelete := FindWebhooksToChange(&Webhooks{}, &Webhooks{pizza})
	if len(*toDelete) != 1 {
		t.Errorf("Forventet å slette 1 webhook, fikk %d", len(*toDelete))
	}

	if len(*toCreate)+len(*toUpdate) != 0 {
		t.Errorf("Forventet 0 webhooks å opprette eller oppdatere, fikk hhv. %d og %d", len(*toCreate), len(*toUpdate))
	}
}

func Test_findWebhooksToChange_create_and_update(t *testing.T) {
	// Skal opprette `bacon`, oppdatere `burger` til `chicken` og la `pizza` være uendret.
	toCreate, toUpdate, toDelete := FindWebhooksToChange(&Webhooks{bacon, chicken, pizza}, &Webhooks{pizza, burger})
	if len(*toCreate) != 1 {
		t.Errorf("Vi skal lage én webhooks, fikk %d", len(*toCreate))
	} else if (*toCreate)[0].Name != "🥓" {
		t.Errorf("Forventet å lage %s, fikk %s", bacon.Name, (*toCreate)[0].Name)
	}
	if len(*toUpdate) != 1 {
		t.Errorf("Forventet å oppdatere 1 webhook, fikk %d", len(*toUpdate))
	} else if (*toUpdate)[0].Name != "🐓" {
		t.Errorf("Forventet å oppdatere %s, fikk %s", chicken.Name, (*toUpdate)[0].Name)
	}
	if len(*toDelete) != 0 {
		t.Errorf("Vi skal ikke slette noen webhooks, fikk %d", len(*toDelete))
	}
}

func Test_findWebhooksToChange_keep_and_delete(t *testing.T) {
	toCreate, toUpdate, toDelete := FindWebhooksToChange(&Webhooks{pizza, burger}, &Webhooks{burger, bacon, pizza})
	if len(*toCreate) != 0 {
		t.Errorf("Forventet å lage 0 webhooks, fikk %d", len(*toCreate))
	}
	if len(*toUpdate) != 0 {
		t.Errorf("Forventet å oppdatere 0 webhooks, fikk %d", len(*toUpdate))
	}
	if len(*toDelete) != 1 {
		t.Errorf("Forventet å slette 1 webhook, fikk %d", len(*toDelete))
	} else if (*toDelete)[0].Name != "🥓" {
		t.Errorf("Forventet å slette %s, fikk %s", bacon.Name, (*toDelete)[0].Name)
	}
}

func Test_rateSimilarWebhooks(t *testing.T) {
	ratedWebhooks := rateCandidateWebhooksSimilarity([]*Webhook{pizza, burger}, []**Webhook{&pizza, &bacon, &chicken})

	if len(ratedWebhooks) != 2 {
		t.Errorf("Forventet 2 rangerte webhooks, fikk %d", len(ratedWebhooks))
	}

	if ratedWebhooks[0].webhook.Name != "🍕" {
		t.Errorf("Forventet å finne %s først i listen, fikk %s", pizza.Name, ratedWebhooks[0].webhook.Name)
	} else if len(ratedWebhooks[0].candidates) != 3 {
		t.Errorf("Forventet å få 3 kandidater, fikk %d", len(ratedWebhooks[0].candidates))
	} else if !(ratedWebhooks[0].candidates[0].similarity >= ratedWebhooks[0].candidates[1].similarity && ratedWebhooks[0].candidates[1].similarity >= ratedWebhooks[0].candidates[2].similarity) {
		t.Errorf("Forventet at kandidatene er sorter i synkende rekkefølge av likhet, fikk %f, %f, %f", ratedWebhooks[0].candidates[0].similarity, ratedWebhooks[0].candidates[1].similarity, ratedWebhooks[0].candidates[2].similarity)
	} else if (*ratedWebhooks[0].candidates[0].webhook).Name != "🍕" {
		t.Errorf("Forventet å finne %s som beste kandidat, fikk %s med likhet %f", pizza.Name, (*ratedWebhooks[0].candidates[0].webhook).Name, ratedWebhooks[0].candidates[0].similarity)
	} else if (*ratedWebhooks[0].candidates[1].webhook).Name != "🥓" {
		t.Errorf("Forventet å finne %s som nest beste kandidat, fikk %s med likhet %f", bacon.Name, (*ratedWebhooks[0].candidates[1].webhook).Name, ratedWebhooks[0].candidates[1].similarity)
	} else if (*ratedWebhooks[0].candidates[2].webhook).Name != "🐓" {
		t.Errorf("Forventet å finne %s som tredje beste kandidat, fikk %s med likhet %f", chicken.Name, (*ratedWebhooks[0].candidates[2].webhook).Name, ratedWebhooks[0].candidates[2].similarity)
	}

	if ratedWebhooks[1].webhook.Name != "🍔" {
		t.Errorf("Forventet å finne %s på andreplass i listen, fikk %s", burger.Name, ratedWebhooks[1].webhook.Name)
	}
}

func Test_sortByBestAvailableCandidate(t *testing.T) {
	sortedRatedWebhooks := rateCandidateWebhooksSimilarity([]*Webhook{pizza, chicken}, []**Webhook{&pizza, &bacon, &burger})
	sortWebhooksByBestAvailableCandidate(sortedRatedWebhooks)

	if sortedRatedWebhooks[0].webhook.Name != "🍕" {
		t.Errorf("Forventet å finne %s først i listen, fikk %s", pizza.Name, sortedRatedWebhooks[0].webhook.Name)
	}

	// Utilgjengeliggjør kandidater som passer best til `pizza`. Står da igjen med én kandidat som passer best til `chicken`
	*(sortedRatedWebhooks[0].candidates[0]).webhook = nil
	*(sortedRatedWebhooks[0].candidates[1]).webhook = nil
	sortWebhooksByBestAvailableCandidate(sortedRatedWebhooks)

	if sortedRatedWebhooks[0].webhook.Name != "🐓" {
		t.Errorf("Forventet å finne %s først i listen, fikk %s", chicken.Name, sortedRatedWebhooks[0].webhook.Name)
	}
}

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
	if !webhookA.Equals(webhookACopy) {
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
