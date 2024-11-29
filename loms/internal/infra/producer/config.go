package producer

import (
	"github.com/IBM/sarama"
)

func PrepareConfig(opts ...Option) *sarama.Config {
	c := sarama.NewConfig()

	{

		c.Producer.Partitioner = sarama.NewHashPartitioner
	}

	// acks параметр
	{
		c.Producer.RequiredAcks = sarama.WaitForAll
	}

	{
		// Уменьшаем пропускную способность, тем самым гарантируем строгий порядок отправки сообщений/батчей
		c.Net.MaxOpenRequests = 1
	}

	{
		/*
			Если эта конфигурация используется для создания `SyncProducer`, оба параметра должны быть установлены
			в значение true, и вы не не должны читать данные из каналов, поскольку это уже делает продьюсер под капотом.
		*/
		c.Producer.Return.Successes = true
		//c.Producer.Return.Errors = true
	}

	for _, opt := range opts {
		_ = opt.Apply(c)
	}

	return c
}
