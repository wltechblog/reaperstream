package main

import (
	"log"
	"sync"
)

type Cars struct {
	Lock        sync.Mutex
	Subscribers []chan Car
	Chan        chan Car
}

var subs []chan Car

var cars Cars

func startCarChan() chan Car {
	cchan := make(chan Car)
	go carChan(cchan)
	return cchan
}

func Subscribe() chan Car {
	cc := make(chan Car)
	cars.Lock.Lock()
	if c.debug {
		log.Println("Adding subscriber")
	}
	cars.Subscribers = append(cars.Subscribers, cc)
	cars.Lock.Unlock()
	return cc
}

func Unsubscribe(me chan Car) {
	cars.Lock.Lock()
	defer cars.Lock.Unlock()
	var newsubs []chan Car
	if len(cars.Subscribers) == 1 {
		cars.Subscribers = nil
		return
	}
	for _, v := range cars.Subscribers {
		if v != me {
			if c.debug {
				log.Println("Removing subscriber")
			}
			newsubs = append(newsubs, v)
		}
		cars.Subscribers = newsubs
	}
}

func carChan(cchan chan Car) {
	for {
		car := <-cchan
		if c.debug {
			log.Println("Got a new car")
		}
		cars.Lock.Lock()
		for x := range cars.Subscribers {
			if c.debug {
				log.Println("Sending car to subscriber")
			}
			cars.Subscribers[x] <- car
		}
		cars.Lock.Unlock()
	}
}
