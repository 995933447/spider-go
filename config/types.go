package config

type(
	Engine interface {
		Run(...Request)
	}

	Request struct {
		Url string
		Fetch func(string) (content []byte, err error)
		Parse func([]byte) ParseResult
		ItemChan chan []Item
		Distinctor Distinctor
	}

	Item interface {

	}

	ParseResult struct {
		Requests []Request
		Items []Item
	}

	Schedule interface {
		Submit(request Request)
		WorkerReady(chan Request)
		Run()
	}

	Distinctor interface {
		RecordFetched(string)
		CheckIsFetched(string) (bool, error)
	}
)

func NilFetch(string) ([]byte, error) {
	return []byte{}, nil
}

func NilParse([]byte) ParseResult {
	return ParseResult{}
}

func NilItemWork() chan []Item {
	itemSaverChan := make(chan []Item)
	go func() {
		for {
			<- itemSaverChan
		}
	}()
	return itemSaverChan
}

type NilDistinctor struct {

}

func (NilDistinctor) RecordFetched(url string) {

}

func (NilDistinctor) CheckIsFetched(url string) (bool, error) {
	return false, nil
}