package kong

import (
	"reflect"
	"sort"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestConsumersService(T *testing.T) {
	assert := assert.New(T)

	client, err := NewClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	consumer := &Consumer{
		Username: String("foo"),
	}

	createdConsumer, err := client.Consumers.Create(defaultCtx, consumer)
	assert.Nil(err)
	assert.NotNil(createdConsumer)

	consumer, err = client.Consumers.Get(defaultCtx, createdConsumer.ID)
	assert.Nil(err)
	assert.NotNil(consumer)

	consumer.Username = String("bar")
	consumer, err = client.Consumers.Update(defaultCtx, consumer)
	assert.Nil(err)
	assert.NotNil(consumer)
	assert.Equal("bar", *consumer.Username)

	err = client.Consumers.Delete(defaultCtx, createdConsumer.ID)
	assert.Nil(err)

	// ID can be specified
	id := uuid.NewV4().String()
	consumer = &Consumer{
		Username: String("foo"),
		ID:       String(id),
	}

	createdConsumer, err = client.Consumers.Create(defaultCtx, consumer)
	assert.Nil(err)
	assert.NotNil(createdConsumer)
	assert.Equal(id, *createdConsumer.ID)

	err = client.Consumers.Delete(defaultCtx, createdConsumer.ID)
	assert.Nil(err)
}

func TestConsumerListEndpoint(T *testing.T) {
	assert := assert.New(T)

	client, err := NewClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	// fixtures
	consumers := []*Consumer{
		&Consumer{
			Username: String("foo1"),
		},
		&Consumer{
			Username: String("foo2"),
		},
		&Consumer{
			Username: String("foo3"),
		},
	}

	// create fixturs
	for i := 0; i < len(consumers); i++ {
		consumer, err := client.Consumers.Create(defaultCtx, consumers[i])
		assert.Nil(err)
		assert.NotNil(consumer)
		consumers[i] = consumer
	}

	consumersFromKong, next, err := client.Consumers.List(defaultCtx, nil)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(consumersFromKong)
	assert.Equal(3, len(consumersFromKong))

	// check if we see all consumers
	assert.True(compareConsumers(consumers, consumersFromKong))

	// Test pagination
	consumersFromKong = []*Consumer{}

	// first page
	page1, next, err := client.Consumers.List(defaultCtx, &ListOpt{Size: 1})
	assert.Nil(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	consumersFromKong = append(consumersFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.Consumers.List(defaultCtx, next)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(2, len(page2))
	consumersFromKong = append(consumersFromKong, page2...)

	assert.True(compareConsumers(consumers, consumersFromKong))

	for i := 0; i < len(consumers); i++ {
		assert.Nil(client.Consumers.Delete(defaultCtx, consumers[i].ID))
	}
}

func compareConsumers(expected, actual []*Consumer) bool {
	var expectedUsernames, actualUsernames []string
	for _, consumer := range expected {
		expectedUsernames = append(expectedUsernames, *consumer.Username)
	}

	for _, consumer := range actual {
		actualUsernames = append(actualUsernames, *consumer.Username)
	}

	return (compareSlices(expectedUsernames, actualUsernames))
}

func compareSlices(expected, actual []string) bool {
	sort.Strings(expected)
	sort.Strings(actual)
	return (reflect.DeepEqual(expected, actual))
}