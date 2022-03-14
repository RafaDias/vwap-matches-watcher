package watchmatches

import (
	"github.com/rafadias/crypto-watcher/internal/infra/providers/exchange/memory"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestWatchMatcherUseCase_Execute(t *testing.T) {
	subscriptions := []string{"BTC-USD", "ETH-USD", "ETH-BTC"}
	channels := []string{"matches"}
	inmemoryExchange := memory.New(subscriptions, channels)
	log := log.New(os.Stdout, "Running tests: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	svc := New(log, inmemoryExchange, 1)
	svc.Execute()

	vwaps := svc.GetVWAP()
	log.Println(vwaps)

	assert.Equal(t,  10.0, vwaps["BTC-USD"])
	assert.Equal(t,  10.0, vwaps["ETH-BTC"])
	assert.Equal(t,  20.0, vwaps["ETH-USD"])

}
