package webhook

import (
	"bucketctl/pkg/types"
	"reflect"
	"testing"
)

func Test_findWebhooksToChange(t *testing.T) {
	webhookA := &types.Webhook{
		Id:     0,
		Name:   "🏍️",
		Events: []string{"🔧", "⚙️", "🔨"},
	}
	webhookB := &types.Webhook{
		Id:     1,
		Name:   "🐓",
		Events: []string{"🥚", "🍗"},
	}
	webhookC := &types.Webhook{
		Id:     1,
		Name:   "🍕",
		Events: []string{"🍍", "🌶️"},
	}

	{
		// Skal opprette A og B. Har ingen å endre
		toCreate, toUpdate, toDelete := findWebhooksToChange([]*types.Webhook{webhookA}, []*types.Webhook{})

		if len(toCreate) != 1 {
			t.Errorf("Vi skal opprette to webhooks")
		}
		if len(toUpdate)+len(toDelete) != 0 {
			t.Errorf("Vi skal hverken oppdatere eller slette webhooks")
		}
	}
	{
		// Skal slette C
		toCreate, toUpdate, toDelete := findWebhooksToChange([]*types.Webhook{}, []*types.Webhook{webhookC})
		if len(toDelete) != 1 {
			t.Errorf("Vi skal slette en webhook")
		}

		if len(toCreate)+len(toUpdate) != 0 {
			t.Errorf("Vi skal hverken opprette eller oppdatere webhooks")
		}

	}
	{
		// Skal oppdatere C til å bli B
		toCreate, toUpdate, toDelete := findWebhooksToChange([]*types.Webhook{webhookB}, []*types.Webhook{webhookC})

		if len(toUpdate) != 1 {
			t.Errorf("Vi skal oppdatere en webhook")
		}
		if toUpdate[0].Id != 1 && toUpdate[0].Name != "🐓" && reflect.DeepEqual(toUpdate[0].Events, []string{"🥚", "🍗"}) {
			t.Errorf("Forventer å være lik B")
		}
		if len(toCreate)+len(toDelete) != 0 {
			t.Errorf("Vi skal hverkan opprette eller slette webhhoks")
		}
	}
	{
		//toCreate, toUpdate, toDelete := findWebhooksToChange([]*types.Webhook{candidate}, []*types.Webhook{webhookC})

	}

}
