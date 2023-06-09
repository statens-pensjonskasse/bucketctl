package webhook

import (
	"bucketctl/pkg/types"
	"github.com/pterm/pterm"
	"testing"
)

func Test_findWebhooksToChange(t *testing.T) {
	motorcycle := &types.Webhook{
		Name:   "🏍️",
		Events: []string{"🔧", "⚙️", "🔨"},
	}
	chicken := &types.Webhook{
		Id:     1,
		Name:   "🐓",
		Events: []string{"🥚", "🍗"},
	}
	pizza := &types.Webhook{
		Id:     2,
		Name:   "🍕",
		Events: []string{"🍍", "🌶️"},
	}
	burger := &types.Webhook{
		Id:     3,
		Name:   "🍔",
		Events: []string{"🍞", "🐄", "🥬", "🍅", "🧅"},
	}

	//{
	//	// Skal opprette A og B. Har ingen å endre
	//	toCreate, toUpdate, toDelete := findWebhooksToChange([]*types.Webhook{motorcycle}, []*types.Webhook{})

	//	if len(toCreate) != 1 {
	//		t.Errorf("Vi skal opprette to webhooks")
	//	}
	//	if len(toUpdate)+len(toDelete) != 0 {
	//		t.Errorf("Vi skal hverken oppdatere eller slette webhooks")
	//	}
	//}
	//{
	//	// Skal slette C
	//	toCreate, toUpdate, toDelete := findWebhooksToChange([]*types.Webhook{}, []*types.Webhook{pizza})
	//	if len(toDelete) != 1 {
	//		t.Errorf("Vi skal slette en webhook")
	//	}

	//	if len(toCreate)+len(toUpdate) != 0 {
	//		t.Errorf("Vi skal hverken opprette eller oppdatere webhooks")
	//	}

	//}
	//{
	//	// Skal oppdatere C til å bli B
	//	toCreate, toUpdate, toDelete := findWebhooksToChange([]*types.Webhook{chicken}, []*types.Webhook{pizza})
	//	if len(toUpdate) != 1 {
	//		t.Errorf("Vi skal oppdatere en webhook")
	//	} else if toUpdate[0].Id != 1 && toUpdate[0].Name != "🐓" && reflect.DeepEqual(toUpdate[0].Events, []string{"🥚", "🍗"}) {
	//		t.Errorf("Forventer å være lik B")
	//	}
	//	if len(toCreate)+len(toDelete) != 0 {
	//		t.Errorf("Vi skal hverken opprette eller slette webhhoks")
	//	}
	//}
	{
		pterm.Error.Println("AAAAAAAA")
		toCreate, toUpdate, toDelete := findWebhooksToChange([]*types.Webhook{motorcycle, chicken, pizza}, []*types.Webhook{pizza, burger})
		if len(toCreate) != 2 {
			t.Errorf("Vi skal lage 2 webhooks, fikk %d", len(toCreate))
		}
		pterm.Info.Println(toCreate[0])
		pterm.Info.Println(toUpdate[0])
		if len(toUpdate) != 1 {
			t.Errorf("Vi skal oppdatere én webhook, fikk %d", len(toUpdate))
		} else if toUpdate[0].Name != "🍕" {
			t.Errorf("Forventet å oppdatere %s, fikk %s", pizza.Name, toUpdate[0].Name)
		}
		if len(toDelete) != 0 {
			t.Errorf("Vi skal ikke slette noen webhooks, fikk %d", len(toDelete))
		}

	}

}
