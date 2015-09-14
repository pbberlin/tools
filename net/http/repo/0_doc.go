// Package repo takes http JSON commands;
// downloading html files in parallel from the designated source;
// making them available via quasi-static http fileserver.
//
// Further inspiration could be taken from github.com/lox/httpcache
//
//  Todo:
//		Re-integrate RSS feeds into the crawling
//		Exlude redirect-responses
//
//
package repo

import (
	"fmt"
	"log"
)

/*

[{ 	'Host':           'www.handelsblatt.com',
 	'SearchPrefix':   '/politik/international',
 	'RssXMLURI':      '/contentexport/feed/schlagzeilen',
}]



curl -X POST -d "[{ \"Host\": \"www.handelsblatt.com\",  \"SearchPrefix\":  \"/politik/deutschland\"         }]"  localhost:8085/fetch/command-receive
curl -X POST -d "[{ \"Host\": \"www.welt.de\"         ,  \"SearchPrefix\":  \"/wirtschaft/deutschland\"      }]"  localhost:8085/fetch/command-receive
curl -X POST -d "[{ \"Host\": \"www.economist.com\"   ,  \"SearchPrefix\":  \"/news/business-and-finance\"   }]"  localhost:8085/fetch/command-receive

curl -X POST -d "[{ \"Host\": \"test.economist.com\"  ,  \"SearchPrefix\":  \"/news/business-and-finance\"   }]"  localhost:8085/fetch/command-receive
curl -X POST -d "[{ \"Host\": \"test.economist.com\"  ,  \"SearchPrefix\":  \"/\"                            }]"  localhost:8085/fetch/command-receive

curl -X POST -d "[{ \"Host\": \"www.welt.de\",           \"SearchPrefix\": \"/wirtschaft/deutschland\" ,  \"RssXMLURI\": \"/wirtschaft/?service=Rss\" }]" localhost:8085/fetch/command-receive


curl localhost:8085/fetch/similar?uri-x=www.welt.de/politik/ausland/article146154432/Tuerkische-Bodentruppen-marschieren-im-Nordirak-ein.html

curl --data url-x=a.com  localhost:8085/fetch/similar
curl --data url-x=https://www.welt.de/politik/ausland/article146154432/Tuerkische-Bodentruppen-marschieren-im-Nordirak-ein.html  localhost:8085/fetch/similar
curl --data url-x=http://www.economist.com/news/britain/21663648-hard-times-hard-hats-making-britain-make-things-again-proving-difficult  localhost:8085/fetch/similar
curl --data url-x=http://www.economist.com/news/americas/21661804-gender-equality-good-economic-growth-girl-power  localhost:8085/fetch/similar

curl --data "cnt=1&url-x=http://www.economist.com/news/americas/21661804-gender-equality-good-economic-growth-girl-power"  localhost:8085/fetch/similar

*/

var pf = fmt.Printf
var pfRestore = fmt.Printf

var spf = fmt.Sprintf
var wpf = fmt.Fprintf

var lpf = log.Printf
