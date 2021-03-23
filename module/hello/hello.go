package hello

import (
	//"rsc.io/quote"
	quoteV3 "rsc.io/quote/v3"
)

func Hello() string {
    //return "Hello, world."
	//return quote.Hello()
	return quoteV3.HelloV3()
}

func Proverb() string {
	return quoteV3.Concurrency()
}

