package fetch

var TestData = map[string][]byte{
	"test.economist.com/someurl": []byte("requesting test.economist.com/someurl will yield this content"),

	"test.economist.com/sections/business-finance/rss.xml": []byte(`<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0" xml:base="http://test.economist.com/sections/business-finance/rss.xml"  xmlns:dc="http://purl.org/dc/elements/1.1/">
<channel>
 <title>Business and finance</title>
 <link>http://test.economist.com/sections/business-finance/rss.xml</link>
 <description></description>
 <language>en</language>
 <image> <url>http://cdn.static-economist.com/sites/default/files/economist_logo_gn.png</url>
 <title>The Economist online</title>
 <link>http://test.economist.com</link>
</image>
<item>
 <title><![CDATA[Tinkering around the edges]]></title>
 <link>http://test.economist.com/news/finance-and-economics/21662584-ukraines-deal-its-creditors-less-impressive-it-appears-tinkering?fsrc=rss</link>
 <description><![CDATA[<div class="esi-generated pullquote"><p>&nbsp;</p><p>&nbsp;</p><p>&nbsp;</p><p>&nbsp;</p></div><p><div class="content-image-full"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/full-width/images/2015/08/articles/main/20150829_fnp501.jpg" alt="" title=""  width="595" height="335" /></div></p><p>AT FIRST sight, it was a triumph. After months of negotiations Ukraine and a committee of its creditors (which include Franklin Templeton, an American investment house and BTG Pactual, a Brazilian one) reached a deal this week to restructure the country’s international bonds, as well as a smattering of other sorts of debt, worth about $18 billion. Payments have been pushed back, meaning that the government will not need to cough up any principal or interest on the debts in question until 2019. The principal on the bonds will also be cut by 20% on average.</p><p><div class="content-image-float-290 retina-290"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/original-size/images/print-edition/20150829_FNC579.png" alt="" title=""  width="580" height="634" /></div></p><p>This is a better deal for Ukraine than many were expecting. It is rare for a country to get a haircut on its debts without also defaulting (one exception is Greece). When the negotiations began, the creditors had refused even to consider writing off any of the country’s...<a href="http://test.economist.com/news/finance-and-economics/21662584-ukraines-deal-its-creditors-less-impressive-it-appears-tinkering?fsrc=rss">Continue reading</a>]]></description>
 <author>The Economist online</author>
 <comments>http://test.economist.com/news/finance-and-economics/21662584-ukraines-deal-its-creditors-less-impressive-it-appears-tinkering#comments</comments>
 <category domain="http://test.economist.com/topics/approved">Approved</category>
 <category domain="http://test.economist.com/category/print-sections/finance-and-economics">Finance and economics</category>
 <category domain="http://test.economist.com/category/web-sections/finance">FINANCE</category>
 <pubDate>Thu, 27 Aug 2015 10:49:10 +0000</pubDate>
 <guid isPermaLink="false">21662584 at http://test.economist.com</guid>
</item>
<item>
 <title><![CDATA[Controversial hybrids]]></title>
 <link>http://test.economist.com/news/business-and-finance/21662355-monsanto-reportedly-raises-its-bid-syngenta-antitrust-regulators-may-take-some?fsrc=rss</link>
 <description><![CDATA[<p><div class="content-image-full"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/full-width/images/2015/08/articles/main/20150829_wbp003.jpg" alt="" title=""  width="595" height="335" /><span class="caption">Roundup Ready—but so are some weeds</span></div></p><p>SIX companies dominate the business of farm supplies. The interest of Monsanto, the world’s biggest seed producer, in buying Syngenta, the largest agrochemicals firm, had threatened to whittle them down to five. That raised worries about whether the reduction in competition would mean less innovation—and thus slower improvements in crop yields—as well as higher costs for farmers.</p><p>However, Syngenta played hard to get. It rebuffed a bid of $45 billion in June. And another, made on August 18th, worth around $47 billion. So, on August 26th, Monsanto walked away. But consolidation of the industry may be in prospect anyway. The takeover battle stimulated the interest of other big agricultural suppliers: BASF, another of the big six, had reportedly sought financing to make a rival offer for Syngenta. And Monsanto itself may not be done. Next year the firm may set its sights on another target, reckons John Klein, an analyst at Berenberg, a bank.</p><p></p><p>Two decades ago the industry was far more fragmented. In 1994 the top four companies in the worldwide market for seeds and crop...<a href="http://test.economist.com/news/business-and-finance/21662355-monsanto-reportedly-raises-its-bid-syngenta-antitrust-regulators-may-take-some?fsrc=rss">Continue reading</a>]]></description>
 <author>The Economist online</author>
 <comments>http://test.economist.com/news/business-and-finance/21662355-monsanto-reportedly-raises-its-bid-syngenta-antitrust-regulators-may-take-some#comments</comments>
 <category domain="http://test.economist.com/topics/approved">Approved</category>
 <category domain="http://test.economist.com/category/web-sections/business">Business</category>
 <category domain="http://test.economist.com/category/print-sections/business">Business</category>
 <pubDate>Tue, 25 Aug 2015 15:30:35 +0000</pubDate>
 <guid isPermaLink="false">21662355 at http://test.economist.com</guid>
</item>
<item>
 <title><![CDATA[What the super-rich expect on holiday]]></title>
 <link>http://test.economist.com/blogs/oliver/2015/08/luxury-travel?fsrc=rss</link>
 <description><![CDATA[<p><div class="content-image-full"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/full-width/images/2015/08/blogs/oliver/20150829_wbp503.jpg" alt="" title=""  width="595" height="335" /></div></p><p>GULLIVER interviewed Kevin Johnson as research for a <a href="http://test.economist.com/news/international/21662558-way-rich-travel-changing-place-lay-your-bread">piece on luxury travel</a> in the print edition this week. Mr Johnson worked as a chief-of-staff and a palace manager for various heads of state, Russian billionaires and royal families. Part of his job was to oversee the travel arrangements of his super-wealthy employers. There wasn’t space to include many of his anecdotes in the print article, but it was a fun chat.</p><p>What was fascinating was, for people for whom money really is no object, just how prosaic some of their needs were. It was not unusual for them to send an employee ahead to some remote island to install their own IT infrastructure, for example, at fantastic cost. Mr Johnson said that this was often done under the guise of cyber-security. Sometimes it was: even on beach holidays they would often bring in their own guy to work with hotel staff to make e-mail and the like completely secure. But often, he suspected, it was more to do with a fear of missing something on the box. “You could be on the most remote island possible...<a href="http://test.economist.com/blogs/oliver/2015/08/luxury-travel?fsrc=rss">Continue reading</a>]]></description>
 <author>The Economist online</author>
 <comments>http://test.economist.com/blogs/oliver/2015/08/luxury-travel#comments</comments>
 <pubDate>Fri, 28 Aug 2015 13:19:41 +0000</pubDate>
 <guid isPermaLink="false">21662616 at http://test.economist.com</guid>
</item>
<item>
 <title><![CDATA[Was the crash that big?]]></title>
 <link>http://test.economist.com/blogs/freeexchange/2015/08/chinas-stockmarket-0?fsrc=rss</link>
 <description><![CDATA[<p><div class="content-image-full"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/full-width/images/2015/08/blogs/free-exchange/20150829_fnp505.jpg" alt="" title=""  width="595" height="335" /></div></p><p><div class="content-image-float-290 retina-290"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/original-size/images/2015/08/blogs/free-exchange/20150829_woc487_0.png" alt="" title=""  width="580" height="1592" /></div></p><p>HOW big was China’s stockmarket crash? Dubbed “Black Monday”, August 24th ended with Chinese equities down 8.5%, wiping out hundreds of billions of dollars in market capitalisation. Like many things about China, this sounds massive. But is a one-day drop of this magnitude that unusual?&nbsp;</p><p>Answering this question requires an understanding of China’s historical market volatility. In the last 25 years, the Shanghai Composite, China’s benchmark stock index, has closed within one percentage point of the previous day’s close on just 56% of all trading days, with an average movement of 0.09% (see chart). The standard deviation of the Shanghai Composite’s daily returns—a measure of market volatility—<span>has been&nbsp;</span>2.6 percentage points.</p><p></p><p>The stockmarkets of rich countries during this period have been considerably less volatile. In the last 25 years, the S&amp;P...<a href="http://test.economist.com/blogs/freeexchange/2015/08/chinas-stockmarket-0?fsrc=rss">Continue reading</a>]]></description>
 <author>The Economist online</author>
 <comments>http://test.economist.com/blogs/freeexchange/2015/08/chinas-stockmarket-0#comments</comments>
 <pubDate>Fri, 28 Aug 2015 12:41:28 +0000</pubDate>
 <guid isPermaLink="false">21662633 at http://test.economist.com</guid>
</item>
<item>
 <title><![CDATA[The Bank of England will follow the Federal Reserve]]></title>
 <link>http://test.economist.com/blogs/freeexchange/2015/08/interest-rates?fsrc=rss</link>
 <description><![CDATA[<p><div class="content-image-full"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/full-width/images/2015/08/blogs/free-exchange/20150829_fnp504.jpg" alt="" title=""  width="595" height="335" /></div></p><p>THE UPWARD revision to American GDP on August 27th provided a shot in the arm to global stockmarkets, which have endured their most volatile week of trading in years. The American economy is now thought to have grown by 3.7% at an annual rate in the second quarter, a much higher estimate than the 2.3% that was given in an initial evaluation of the quarter. After so much bad news about the Chinese economy, which has been a factor behind the meltdown of Chinese share indices, the US figure came as a relief to markets. Commodities also rallied: Brent crude soared by 10%, its biggest one-day gain in over six years, to settle above $47.50 a barrel. On August 28th Britain's statistics office left unchanged its estimate of second-quarter GDP growth in the UK economy at 2.6%.</p><p>With the American economy racing ahead, investors wondered whether the Federal Reserve might indeed raise interest rates in September, a move that the central bank had been preparing markets for but which looked increasingly unlikely given the turmoil of the past week or so. The case for a rate increase next month remains <a...<a href="http://test.economist.com/blogs/freeexchange/2015/08/interest-rates?fsrc=rss">Continue reading</a>]]></description>
 <author>The Economist online</author>
 <comments>http://test.economist.com/blogs/freeexchange/2015/08/interest-rates#comments</comments>
 <pubDate>Fri, 28 Aug 2015 12:22:23 +0000</pubDate>
 <guid isPermaLink="false">21662631 at http://test.economist.com</guid>
</item>
<item>
 <title><![CDATA[Core concern]]></title>
 <link>http://test.economist.com/news/finance-and-economics/21662606-life-keeps-getting-harder-japans-central-bankers-core-concern?fsrc=rss</link>
 <description><![CDATA[<p><div class="content-image-float-290 retina-290"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/original-size/images/print-edition/20150829_FNC538.png" alt="" title=""  width="580" height="598" /></div></p><p>AFTER two years of remission, Japan seems likely to sink back into the “chronic disease” of deflation, as Haruhiko Kuroda, the governor of the Bank of Japan (BoJ), calls it. New data are expected to show on August 28th that core CPI, the central bank’s preferred indicator of inflation, turned negative in July for the first time since the bank launched a big programme of quantitative easing (printing money to buy bonds) in April 2013 (see chart). At the time, it pledged to lift inflation to 2% in two years.</p><p>The news will heap further pressure on the BoJ to ease monetary policy yet more this year, as will worries about Chinese growth. The fact that Japan’s economy shrank by 1.6% in the second quarter on an annualised basis adds to the concerns. The central bank is currently buying about ¥80 trillion ($670 billion) of long-term Japanese government bonds (JGBs) a year, or twice the annual issuance. It now holds over ¥300 trillion of JGBs, or nearly a third of all outstanding bonds.</p><p></p><p>Mr Kuroda’s excuse for deflation’s apparent return is that the falling oil price has pushed down core CPI, which excludes fresh food but...<a href="http://test.economist.com/news/finance-and-economics/21662606-life-keeps-getting-harder-japans-central-bankers-core-concern?fsrc=rss">Continue reading</a>]]></description>
 <author>The Economist online</author>
 <comments>http://test.economist.com/news/finance-and-economics/21662606-life-keeps-getting-harder-japans-central-bankers-core-concern#comments</comments>
 <category domain="http://test.economist.com/topics/approved">Approved</category>
 <category domain="http://test.economist.com/category/print-sections/finance-and-economics">Finance and economics</category>
 <category domain="http://test.economist.com/category/web-sections/finance">FINANCE</category>
 <pubDate>Thu, 27 Aug 2015 14:43:29 +0000</pubDate>
 <guid isPermaLink="false">21662606 at http://test.economist.com</guid>
</item>
<item>
 <title><![CDATA[Non-profit paradise]]></title>
 <link>http://test.economist.com/news/finance-and-economics/21662604-low-prices-attract-new-type-customer-non-profit-paradise?fsrc=rss</link>
 <description><![CDATA[<p><div class="content-image-float-290"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/290-width/images/print-edition/20150829_FNP003_0.jpg" alt="" title=""  width="290" height="435" /><span class="caption">Frond management</span></div></p><p>JUST off the west coast of Florida lies the sun-drenched island of Little Bokeelia. It is blessed with cascading waterfalls, tennis courts, pools and a Spanish-style villa. Despite such enticing features, the island languished on the market for three years, before selling in July for a mere $14.5m—half the original asking price.</p><p>Little Bokeelia is not the only island that is proving hard to shift. In the Bahamas, where prices per acre are among the world’s highest, hundreds of atolls lie unbought. The price of undeveloped islands, which make up around 80% of the market, has dropped roughly by half since the financial crisis, says Farhad Vladi, a private-island broker.</p><p></p><p>In the early 2000s private islands were the trophy of choice for millionaires but the recession sapped demand. Building on an island is much pricier than on a mainland plot, and there are many potential pitfalls. It is not for the faint-hearted, says Edward Childs of Smiths Gore, an estate agent in the British Virgin Islands. Mega-yachts and private jets are seen as more predictable investments. As a result private islands can...<a href="http://test.economist.com/news/finance-and-economics/21662604-low-prices-attract-new-type-customer-non-profit-paradise?fsrc=rss">Continue reading</a>]]></description>
 <author>The Economist online</author>
 <comments>http://test.economist.com/news/finance-and-economics/21662604-low-prices-attract-new-type-customer-non-profit-paradise#comments</comments>
 <category domain="http://test.economist.com/topics/approved">Approved</category>
 <category domain="http://test.economist.com/category/print-sections/finance-and-economics">Finance and economics</category>
 <category domain="http://test.economist.com/category/web-sections/finance">FINANCE</category>
 <pubDate>Thu, 27 Aug 2015 14:43:29 +0000</pubDate>
 <guid isPermaLink="false">21662604 at http://test.economist.com</guid>
</item>
<item>
 <title><![CDATA[Wheelbarrows to the rescue]]></title>
 <link>http://test.economist.com/news/finance-and-economics/21662603-nigerias-central-bank-pioneers-new-method-shore-up-local?fsrc=rss</link>
 <description><![CDATA[<p><div class="content-image-full"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/full-width/images/print-edition/20150829_FNP004_1.jpg" alt="" title=""  width="595" height="335" /></div></p><p>YOU might think Godwin Emefiele, the governor of Nigeria’s central bank, had problems enough. The collapsing oil price has slashed Nigeria’s export earnings. Foreign reserves have fallen from more than $40 billion early last year to just over $30 billion now. In response Mr Emefiele (pictured) devalued the local currency, the naira, in November and again in February. The devaluations are stoking inflation. Like many other central bankers in commodity-exporting countries, he is faced with the unenviable choice of raising rates despite the damage to an already faltering economy, or leaving them be despite rising inflation and a swooning currency. Unlike other central bankers, however, Mr Emefiele has decided to compound the awkwardness of his position by getting involved in industrial policy as well.</p><p>In June the central bank said it would not provide foreign exchange for 41 categories of imports, ranging from wheelbarrows to private jets. The idea, Mr Emefiele says, is both to conserve dollars and to stimulate local manufacturing. “Central banks in developing countries like ours cannot sit idly by and concentrate only on price and monetary...<a href="http://test.economist.com/news/finance-and-economics/21662603-nigerias-central-bank-pioneers-new-method-shore-up-local?fsrc=rss">Continue reading</a>]]></description>
 <author>The Economist online</author>
 <comments>http://test.economist.com/news/finance-and-economics/21662603-nigerias-central-bank-pioneers-new-method-shore-up-local#comments</comments>
 <category domain="http://test.economist.com/topics/approved">Approved</category>
 <category domain="http://test.economist.com/category/print-sections/finance-and-economics">Finance and economics</category>
 <category domain="http://test.economist.com/category/web-sections/finance">FINANCE</category>
 <pubDate>Thu, 27 Aug 2015 14:43:29 +0000</pubDate>
 <guid isPermaLink="false">21662603 at http://test.economist.com</guid>
</item>
<item>
 <title><![CDATA[The kindness of neighbours]]></title>
 <link>http://test.economist.com/news/finance-and-economics/21662596-new-paper-asks-why-some-tanzanian-farmers-are-more-selfish-others?fsrc=rss</link>
 <description><![CDATA[<p><div class="content-image-full"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/full-width/images/print-edition/20150829_FNP002_0.jpg" alt="" title=""  width="595" height="335" /><span class="caption">Turning his back</span></div></p><p>IN THE semi-arid lowlands of Mufindi, in southern Tanzania, water is hard to come by. Villagers rely on irrigation to grow maize, potatoes and spinach. Informal and often woolly codes govern how much water each farmer diverts to their own fields, and how much they leave for their neighbours downstream. Some farmers, naturally, turn out to be more grasping than others. Economists typically see such decisions as irreducible: there is no accounting for individuals’ values and preferences. But a new study* investigates why there is such variation in generosity among Mufindi’s farmers.</p><p>The researchers asked other villagers to rank each farmer’s social status on a scale of one to four. Then they invited the farmers to take part in a game in which participants had to decide how much water they would take under different scenarios. Participants were paid small sums, which varied according to how well they did in the game. They received more money if they reaped a bigger harvest by taking more than their share of water, but less if the other villagers fined them for violating water-sharing...<a href="http://test.economist.com/news/finance-and-economics/21662596-new-paper-asks-why-some-tanzanian-farmers-are-more-selfish-others?fsrc=rss">Continue reading</a>]]></description>
 <author>The Economist online</author>
 <comments>http://test.economist.com/news/finance-and-economics/21662596-new-paper-asks-why-some-tanzanian-farmers-are-more-selfish-others#comments</comments>
 <category domain="http://test.economist.com/topics/approved">Approved</category>
 <category domain="http://test.economist.com/category/print-sections/finance-and-economics">Finance and economics</category>
 <category domain="http://test.economist.com/category/web-sections/finance">FINANCE</category>
 <pubDate>Thu, 27 Aug 2015 14:43:29 +0000</pubDate>
 <guid isPermaLink="false">21662596 at http://test.economist.com</guid>
</item>
<item>
 <title><![CDATA[Off the block]]></title>
 <link>http://test.economist.com/news/finance-and-economics/21662595-economists-may-idolise-auctions-most-people-do-not-block?fsrc=rss</link>
 <description><![CDATA[<p><div class="content-image-float-290 retina-290"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/original-size/images/print-edition/20150829_FNC545.png" alt="" title=""  width="580" height="634" /></div></p><p>THE first item sold on eBay, an online marketplace, was a broken laser pointer, which was snapped up for $14.83 in September 1995. By 2002 eBay had hosted nearly $15 billion of transactions and had more registered users than Britain had people. Yet the fad for online auctions faded almost as quickly as it appeared. Only 20% of sales on eBay, which turns 20 on September 3rd, now involve auctions.</p><p>At eBay’s inception, users could sell things only by auction. This was tremendously exciting for economists, who love the things for their ability to magic prices out of thin air and to allocate goods efficiently by determining who values them most highly. The main obstacle to holding auctions is the cost of bringing together enough interested buyers and sellers. But eBay made connecting buyers and sellers cheap. Without it, that broken laser pointer may well have languished unsold.</p><p></p><p>EBay also benefited from a first-mover advantage. Buyers want to go where there are lots of competing sellers, and sellers will flock to wherever they can find the most eager customers. The size of eBay’s network was its own, self-perpetuating engine of...<a href="http://test.economist.com/news/finance-and-economics/21662595-economists-may-idolise-auctions-most-people-do-not-block?fsrc=rss">Continue reading</a>]]></description>
 <author>The Economist online</author>
 <comments>http://test.economist.com/news/finance-and-economics/21662595-economists-may-idolise-auctions-most-people-do-not-block#comments</comments>
 <category domain="http://test.economist.com/topics/approved">Approved</category>
 <category domain="http://test.economist.com/category/print-sections/finance-and-economics">Finance and economics</category>
 <category domain="http://test.economist.com/category/web-sections/finance">FINANCE</category>
 <pubDate>Thu, 27 Aug 2015 14:43:29 +0000</pubDate>
 <guid isPermaLink="false">21662595 at http://test.economist.com</guid>
</item>
<item>
 <title><![CDATA[Building works]]></title>
 <link>http://test.economist.com/news/finance-and-economics/21662593-historic-opportunity-improve-infrastructure-cheap-danger?fsrc=rss</link>
 <description><![CDATA[<p><div class="content-image-full"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/full-width/images/print-edition/20150829_FNP001_0.jpg" alt="" title=""  width="595" height="335" /></div></p><p>IT IS hard to exaggerate the decrepitude of infrastructure in much of the rich world. One in three railway bridges in Germany is over 100 years old, as are half of London’s water mains. In America the average bridge is 42 years old and the average dam 52. The American Society of Civil Engineers rates around 14,000 of the country’s dams as “high hazard” and 151,238 of its bridges as “deficient”. This crumbling infrastructure is both dangerous and expensive: traffic jams on urban highways cost America over $100 billion in wasted time and fuel each year; congestion at airports costs $22 billion and another $150 billion is lost to power outages.</p><p>The B20, the business arm of the G20, a club of big economies, estimates that the global backlog of spending needed to bring infrastructure up to scratch will reach $15 trillion-20 trillion by 2030. McKinsey, a consultancy, reckons that in 2007-12 investment in infrastructure in rich countries was about 2.5% of GDP a year when it should have been 3.5%. If anything, the problem is becoming more acute as some governments whose finances have been racked by the crisis cut back. In 2013 in the euro zone, general...<a href="http://test.economist.com/news/finance-and-economics/21662593-historic-opportunity-improve-infrastructure-cheap-danger?fsrc=rss">Continue reading</a>]]></description>
 <author>The Economist online</author>
 <comments>http://test.economist.com/news/finance-and-economics/21662593-historic-opportunity-improve-infrastructure-cheap-danger#comments</comments>
 <category domain="http://test.economist.com/topics/approved">Approved</category>
 <category domain="http://test.economist.com/category/print-sections/finance-and-economics">Finance and economics</category>
 <category domain="http://test.economist.com/category/web-sections/finance">FINANCE</category>
 <pubDate>Thu, 27 Aug 2015 14:43:28 +0000</pubDate>
 <guid isPermaLink="false">21662593 at http://test.economist.com</guid>
</item>
<item>
 <title><![CDATA[Patently problematic]]></title>
 <link>http://test.economist.com/news/business/21662552-proposals-consistent-global-rules-company-tax-cause-worries-all-round-patently?fsrc=rss</link>
 <description><![CDATA[<p><div class="content-image-float-290 retina-290"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/original-size/images/print-edition/20150829_WBC521.png" alt="" title=""  width="580" height="816" /></div></p><p>CLARITY or chaos? Supporters of the Base Erosion and Profit-Shifting (BEPS) project, being worked on by the OECD, argue that it will bind multinationals to a consistent set of global tax rules, providing them with less licence than they now have to short-change governments through artful use of loopholes in national laws. Sceptics worry that it could only lead to chaos if countries adopt the new guidelines to differing degrees, or if some governments conclude they are too soft, and take unilateral action to stop tax revenue on profits being siphoned abroad.</p><p>With two months to go before the club of rich and middle-income countries presents its plan to the G20 for approval, much of its detail remains unclear. But no one doubts that, with so many clashing national interests at stake, there are limits to what can be achieved. Discord has been evident lately, even among allies: in June a US Treasury official accused Britain and Australia of undermining international agreement by “going their own way”. One beef is over Britain’s new “diverted profits” tax, which imposes a levy on profits routed to tax havens through “contrived...<a href="http://test.economist.com/news/business/21662552-proposals-consistent-global-rules-company-tax-cause-worries-all-round-patently?fsrc=rss">Continue reading</a>]]></description>
 <author>The Economist online</author>
 <comments>http://test.economist.com/news/business/21662552-proposals-consistent-global-rules-company-tax-cause-worries-all-round-patently#comments</comments>
 <category domain="http://test.economist.com/topics/approved">Approved</category>
 <category domain="http://test.economist.com/category/web-sections/business">Business</category>
 <category domain="http://test.economist.com/category/print-sections/business">Business</category>
 <pubDate>Thu, 27 Aug 2015 14:43:23 +0000</pubDate>
 <guid isPermaLink="false">21662552 at http://test.economist.com</guid>
</item>
</channel>
</rss>`),
	"test.economist.com": []byte(`<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml"
      xml:lang="en" lang="en" dir="ltr"
      xmlns:og="http://ogp.me/ns#"
      xmlns:fb="https://www.facebook.com/2008/fbml">

<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<title>The Economist - World News, Politics, Economics, Business & Finance</title>
<link rel="shortcut icon" href="http://cdn.static-economist.com/sites/default/files/econfinal_favicon.ico" type="image/x-icon" />
<link rel="canonical" href="http://test.economist.com/" />
<meta name="description" content="The Economist offers authoritative insight and opinion on international news, politics, business, finance, science, technology and the connections between them." />
<meta name="pubdate" content="20120516" />
<meta name="revisit-after" content="1 day" />
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:site" content="@TheEconomist" />
<meta property="fb:app_id" content="173277756049645" />
<meta property="og:image" content="http://cdn.static-economist.com/sites/default/files/the-economist-logo.gif" />
<meta property="og:site_name" content="The Economist" />
<link rel="apple-touch-icon" sizes="144x144" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-144x144.png" /><link rel="apple-touch-icon" sizes="120x120" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-120x120.png" /><link rel="apple-touch-icon" sizes="114x114" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-114x114.png" /><link rel="apple-touch-icon" sizes="72x72" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-72x72.png" /><link rel="apple-touch-icon" href="http://cdn.static-economist.com/sites/default/files/images/icons/touch-icon-iphone.png" /><link rel="apple-touch-icon-precomposed" sizes="144x144" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-144x144-precomposed.png"><link rel="apple-touch-icon-precomposed" sizes="120x120" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-120x120-precomposed.png"><link rel="apple-touch-icon-precomposed" sizes="114x114" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-114x114-precomposed.png"><link rel="apple-touch-icon-precomposed" sizes="72x72" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-72x72-precomposed.png"><link rel="apple-touch-icon-precomposed" href="http://cdn.static-economist.com/sites/default/files/images/icons/touch-icon-iphone-precomposed.png">  <link type="text/css" rel="stylesheet" media="all" href="http://cdn.static-economist.com/sites/default/files/css/css_5b7f88bc4d30779c79a733cfcf63a4c8.css" />
<link rel="publisher" href="https://plus.google.com/100470681032489535736" />
</head>

<body class="front not-logged-in page-node node-type-homepage one-sidebar sidebar-right path-node-21555491 world-menu business-menu economics-menu printedition-menu science-technology-menu culture-menu">
      
        <div id="fb-root"></div>
<div id="leaderboard" class="clearfix">
    <div id="block-ec_ads-leaderboard_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div id="leaderboard-ad"><!-- Site: Web.  Zone: Home |  --> <div id="gpt_leaderboard_ad" data-cb-ad-id="Leaderboard ad">
    </div></div>  </div>
</div>

<div id="block-ec_ads-subscription_ad" class="block block-ec_ads ec-ads-gpt">
    <div class="content clearfix">
    <div id="subslug-ad"><!-- Site: Web.  Zone: Home |  --> <div id="gpt_subscription_ad" data-cb-ad-id="Subscription ad">

</div></div>  </div>
</div>

<div id="block-ec_ads-pencil_slug_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_pencil_slug_ad" data-cb-ad-id="Pencil slug ad">

</div>  </div>
</div>

  </div>
<header class="mh">
  <!-- 
  To use as a standalone component please wrap it up into a <header class="mh"></header> 
-->
<div class="mh-stripe">
  <div class="mh-stripe-wrap">
    <ul class="mh-user-menu"><li class="first"><span>More from The Economist</span><ul><li class="first"><a href="/digital">The Economist digital editions</a></li>
<li><a href="/newsletters">Newsletters</a></li>
<li><a href="/events">Events</a></li>
<li><a href="http://jobs.economist.com">Jobs.Economist.com</a></li>
<li><a href="http://store.economist.com/">The Economist Store</a></li>
<li class="last"><a href="/bookmarks" data-ec-bookmark-click="20|9408299|bookmark &gt; more from the economist &gt; BM saved" data-ec-omniture="masthead|act_prod|bookmarks">Timekeeper reading list</a></li>
</ul></li>
<li><span>My Subscription</span><ul><li class="first"><a href="/products/subscribe">Subscribe to The Economist</a></li>
<li><a href="/activate">Activate my digital subscription</a></li>
<li><a href="/user">Manage my subscription</a></li>
<li class="last"><a href="/products/renew">Renew</a></li>
</ul></li>
<li class="masthead-user"><a href="/user/login?destination=node%2F21555491" class="show-login">Log in or register</a></li>
<li class="masthead-subscribe even last"><a href="https://subscriptions.economist.com/GLB/MAST/T1" class="show-subscribe">Subscribe</a></li>
</ul>
          <div class="mh-search">
        <form action="http://google.com/search"  accept-charset="UTF-8" method="GET" id="search-theme-form">
<div><div id="search" class="container-inline">
  <div class="form-item clearfix" id="edit-search-theme-form-1-wrapper">
 <label for="edit-search-theme-form-1">Search this site:</label>
<input type="text" maxlength="128" name="query" id="edit-search-theme-form-1" size="15" value="" title="Enter the terms you wish to search for." autocorrect="off" class="form-text search-field" />
</div>
<input type="submit" name="op" id="edit-submit" value="Search"  class="form-submit" />
<input type="hidden" name="form_id" id="edit-search-theme-form" value="search_theme_form"  />
<input type="hidden" name="sitesearch" id="edit-sitesearch" value="economist.com"  />
</div>

</div></form>
      </div>
      </div> <!-- /.mh-stripe-wrap -->
</div> <!-- /.mh-stripe -->  <div class="mh-nav mh-big">
    <div class="mh-nav-wrap">
              <h1 class="svg-logo"><a href="/" class="active"><img class="mh-logo" width="170" height="85" src="//cdn.static-economist.com/sites/all/themes/econfinal/images/svg/logo.svg" alt="The Economist" /></a></h1>
                    <nav>
          <ul class="mh-nav-links"><li class="first"><a href="/content/politics-this-week" title="" class="sub-menu-link">World politics</a><ul class="mh-subnav"><li class="first"><a href="/content/politics-this-week">Politics this week</a></li>
<li><a href="/sections/united-states">United States</a></li>
<li><a href="/sections/britain">Britain</a></li>
<li><a href="/sections/europe">Europe</a></li>
<li><a href="/sections/china">China</a></li>
<li><a href="/sections/asia">Asia</a></li>
<li><a href="/sections/newcontinent">xxx</a></li>
<li><a href="/sections/americas">Americas</a></li>
<li><a href="/sections/middle-east-africa">Middle East &amp; Africa</a></li>
<li class="last"><a href="/sections/international">International</a></li>
</ul></li>
<li><a href="/sections/business-finance" class="sub-menu-link">Business &amp; finance</a><ul class="mh-subnav"><li class="first"><a href="/sections/business-finance">All Business &amp; finance</a></li>
<li class="even last"><a href="/whichmba">Which MBA?</a></li>
</ul></li>
<li class=""><a href="/sections/economics" class="sub-menu-link">Economics</a><ul class="mh-subnav"><li class="first"><a href="/sections/economics">All Economics</a></li>
<li><a href="/economics-a-to-z">Economics A-Z</a></li>
<li><a href="/markets-data">Markets &amp; data</a></li>
<li class="even last"><a href="/indicators">Indicators</a></li>
</ul></li>
<li><a href="/sections/science-technology" class="sub-menu-link">Science &amp; technology</a><ul class="mh-subnav"><li class="first"><a href="/sections/science-technology">All Science &amp; technology</a></li>
<li class="even last"><a href="/technology-quarterly" title="Technology Quarterly">Technology Quarterly</a></li>
</ul></li>
<li class=""><a href="/sections/culture" class="sub-menu-link">Culture</a><ul class="mh-subnav"><li class="first"><a href="/sections/culture">All Culture</a></li>
<li><a href="http://moreintelligentlife.com/">More Intelligent Life</a></li>
<li><a href="/styleguide/introduction">Style guide</a></li>
<li class="even last"><a href="/economist-quiz">The Economist Quiz</a></li>
</ul></li>
<li><a href="/blogs" class="sub-menu-link">Blogs</a><ul class="mh-subnav"><li class="first"><a href="/blogs">Latest updates</a></li>
<li><a href="/blogs/buttonwood" title="Financial markets">Buttonwood&#039;s notebook</a></li>
<li><a href="/blogs/democracyinamerica" title="American politics">Democracy in America</a></li>
<li><a href="/blogs/erasmus">Erasmus</a></li>
<li><a href="/blogs/freeexchange" title="Economics">Free exchange</a></li>
<li><a href="/blogs/gametheory" title="Sports">Game theory</a></li>
<li><a href="/blogs/graphicdetail" title="Charts, maps and infographics">Graphic detail</a></li>
<li><a href="/blogs/gulliver" title="Business travel">Gulliver</a></li>
<li><a href="/blogs/prospero" title="Books, arts and culture">Prospero</a></li>
<li class="even last"><a href="/blogs/economist-explains">The Economist explains</a></li>
</ul></li>
<li class=""><a href="http://test.economist.com/debate" class="sub-menu-link">Debate</a><ul class="mh-subnav"><li class="first"><a href="http://test.economist.com/debate">Economist debates</a></li>
<li class="even last"><a href="/content/letters-to-the-editor" title="">Letters to the editor</a></li>
</ul></li>
<li><a href="/multimedia" class="sub-menu-link">Multimedia</a><ul class="mh-subnav"><li class="first"><a href="/films">Economist Films</a></li>
<li><a href="http://radio.economist.com">Economist Radio</a></li>
<li><a href="/multimedia">Multimedia library</a></li>
<li class="even last"><a href="/audio-edition">The Economist in audio</a></li>
</ul></li>
<li class="last"><a href="/printedition" title="" class="sub-menu-link">Print edition</a><ul class="mh-subnav"><li class="first"><a href="/printedition/">Current issue</a></li>
<li><a href="/printedition/covers">Previous issues</a></li>
<li><a href="/printedition/specialreports">Special reports</a></li>
<li><a href="/content/politics-this-week">Politics this week</a></li>
<li><a href="/content/business-this-week">Business this week</a></li>
<li><a href="/sections/leaders">Leaders</a></li>
<li><a href="/printedition/kallery">KAL&#039;s cartoon</a></li>
<li class="even last"><a href="/sections/obituary">Obituaries</a></li>
</ul></li>
</ul>                  </nav>
          </div>
  </div> <!-- /.mh-nav -->
</header> <!-- /header -->
  <div id="page" class="container">
    <a name="top" id="navigation-top"></a>

            
    <div id="columns" class="clearfix">
                        <div id="leadspot" class="grid-16 clearfix">
        <div id="block-ec_homepage-ec_homepage_superhero" class="block block-ec_homepage 
">
    <div class="content clearfix">
        <div id="superhero" class="clearfix">
      <div class="hero-superhero"><ul id="hero" class="hero-multiple"><li class="selected"><div class="hero-item hero-item-1"><div class="hero-comment"><a href="/node/21662544/comments#comments" title="Comments" class="comment-icon"><span>5</span></a></div><a href="/news/leaders/21662544-fear-about-chinas-economy-can-be-overdone-investors-are-right-be-nervous-great-fall" class="hero-tab"><h2 class="fly-title">Financial markets</h2><p class="headline">The Great Fall of China</p></a></div><div class="hero-media"><a href="/news/leaders/21662544-fear-about-chinas-economy-can-be-overdone-investors-are-right-be-nervous-great-fall" class="hero-image" rel="nofollow"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/superhero/20150829_LDP001_473.jpg" alt="The Great Fall of China" title="The Great Fall of China"  class="imagecache imagecache-superhero" width="473" height="266" /></a></div></li><li><div class="hero-item hero-item-2"><a href="/news/finance-and-economics/21662584-ukraines-deal-its-creditors-less-impressive-it-appears-tinkering" class="hero-tab"><h2 class="fly-title">Ukraine’s debt restructuring</h2><p class="headline">Tinkering at the edges</p></a></div><div class="hero-media" style="visibility:hidden;opacity:0"><a href="/news/finance-and-economics/21662584-ukraines-deal-its-creditors-less-impressive-it-appears-tinkering" class="hero-image" rel="nofollow"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/superhero/20150829_FNP503_473.jpg" alt="Tinkering at the edges" title="Tinkering at the edges"  class="imagecache imagecache-superhero" width="473" height="266" /></a></div></li><li><div class="hero-item hero-item-3"><a href="/blogs/graphicdetail/2015/08/american-and-british-flight-safety-airlines-v-light-aircraft" class="hero-tab"><h2 class="fly-title">Flight safety</h2><p class="headline">Perils of private planes</p></a></div><div class="hero-media" style="visibility:hidden;opacity:0"><a href="/blogs/graphicdetail/2015/08/american-and-british-flight-safety-airlines-v-light-aircraft" class="hero-image" rel="nofollow"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/original-size/20150829_WOC584_473.png" alt="Perils of private planes" title="Perils of private planes"  class="imagecache imagecache-original-size" width="946" height="532" /></a></div></li><li><div class="hero-item hero-item-4"><div class="hero-comment"><a href="/blogs/economist-explains/2015/08/economist-explains-20#comments" title="Comments" class="comment-icon"><span>6</span></a></div><a href="/blogs/economist-explains/2015/08/economist-explains-20" class="hero-tab"><h2 class="fly-title">Black-hole theory</h2><p class="headline">Hawking a new idea</p></a></div><div class="hero-media" style="visibility:hidden;opacity:0"><a href="/blogs/economist-explains/2015/08/economist-explains-20" class="hero-image" rel="nofollow"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/superhero/20150829_BLP509_473.jpg" alt="Hawking a new idea" title="Hawking a new idea"  class="imagecache imagecache-superhero" width="473" height="266" /></a></div></li></ul></div>
      <div class="cover-image-container"><a href="/printedition" class="cover-image" data-ec-omniture="home|touts|issue|cover"></a><ul><li><a href="/printedition" data-ec-omniture="home|touts|issue|printedition">Full contents</a></li><li class="last"><a href="/products/subscribe" data-ec-omniture="home|touts|issue|subscribe">Subscribe</a></li></ul></div>
    </div>  </div>
</div>

      </div>
      
      <div id="column-content" class="grid-10 grid-first clearfix">
                                <!-- Create left column on search pages -->
                                                  <!-- DoubleClick Floodlight Tag: Please do not remove -->
<!-- Homepage node -->
<div class="grid-7 grid-first push-3">
  <div id="homepage-center-inner">
    <section class="news-package typog-package">

  <h1 class="fly-title">Shia, not shale</h1>

    <article>
  <a href="/news/finance-economics/21662570-kingdom-can-stand-more-pain-it-will-take-much-cheaper-oil-saudi-arabia-take-action" >
          <div>
          <h2 class="headline">It will take much cheaper oil for Saudi Arabia to take action</h2>
      <p class="rubric">
        The kingdom can stand more pain        <span data-href-redirect="/node/21662570/comments#comments" class="comment-icon"><span>0</span></span>      </p>
    </div>
      </a>
</article>
  
  
    <div class="package-more"><a href="/world/middle-east" title="Middle East &amp;amp; Africa">More in Middle East &amp; Africa &raquo;</a></div>
  </section><section class="news-package typog-package">

  <h1 class="fly-title">Murder and politics in Northern Ireland</h1>

    <article>
  <a href="/news/britain/21662500-murder-former-ira-man-causes-political-tremors-consequences-killing" >
          <div>
          <h2 class="headline">The consequences of a killing</h2>
      <p class="rubric">
        The murder of a former IRA man causes political tremors        <span data-href-redirect="/node/21662500/comments#comments" class="comment-icon"><span>1</span></span>      </p>
    </div>
      </a>
</article>
  
  
    <div class="package-more"><a href="/sections/britain" title="Britain">More in Britain &raquo;</a></div>
  </section><section class="ec-homepage-player news-package typog-package">
      <h1 class="fly-title">Latest audio and video</h1>
    <div style="height:422px;"><object id="ec-homepage-video" class="BrightcoveExperience"><param name="bgcolor" value="#FFFFFF" /><param name="isUI" value="true" /><param name="isVid" value="true" /><param name="dynamicStreaming" value="true" /><param name="autoStart" value="false" /><param name="wmode" value="opaque" /><param name="includeAPI" value="true" /><param name="linkBaseURL" value="http://test.economist.com/multimedia" /><param name="playerID" value="1545427201001" /><param name="playerKey" value="AQ~~,AAABDH-R__E~,dB4S9tmhdOrAcjB6eqWZCo1XXp-OU2vB" /><param name="width" value="402" /><param name="height" value="422" /><param name="templateLoadHandler" value="ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.omnitureBCTemplateLoaded" /><param name="@videoPlayer" value="" /><param name="labels" value="http://cdn.static-economist.com/sites/all/modules/custom/ec_brightcove/EcBcLables.xml" /></object></div>
  <p class="package-more">
    <a href="/multimedia">More in Multimedia &raquo;</a>  </p>
</section><section class="news-package typog-package">

  <h1 class="fly-title">Economics</h1>

    <article>
  <a href="/blogs/freeexchange/2015/08/american-economy" >
          <div>
          <h2 class="headline">The American economy</h2>
      <p class="rubric">
        How exposed are American households to the stock market?        <span data-href-redirect="/blogs/freeexchange/2015/08/american-economy#comments" class="comment-icon"><span>20</span></span>      </p>
    </div>
      </a>
</article>
  
  
    <div class="package-more"><a href="/sections/business-finance" title="Business &amp;amp; finance">More in Business &amp; finance &raquo;</a></div>
  </section><section class="news-package typog-package">

  <h1 class="fly-title">European mobile telecoms</h1>

    <article>
  <a href="/news/business/21661660-eus-new-competition-chief-will-have-rule-wave-mergers-together-we-stand" >
          <div>
          <h2 class="headline">Together we stand</h2>
      <p class="rubric">
        The EU’s new competition chief will have to rule on a wave of mergers        <span data-href-redirect="/node/21661660/comments#comments" class="comment-icon"><span>6</span></span>      </p>
    </div>
      </a>
</article>
  
    <ul><li class="first last"><div class="">
  <a href="https://espresso.economist.com/385877ed6e1207dc4a965ffe024e7862" class="headline">Shedding some light on Europe’s most cut-throat mobile market</a>  </div>
</li>
</ul>  
    <div class="package-more"><a href="/sections/business-finance" title="Business &amp;amp; finance">More in Business &amp; finance &raquo;</a></div>
  </section><section class="news-package typog-package">

  <h1 class="fly-title">Singapore’s democracy</h1>

    <article>
  <a href="/news/asia/21662410-fifty-years-singapores-ruling-party-looks-secure-unequal-contest" >
          <div>
          <h2 class="headline">The first election since the death of Lee Kuan Yew</h2>
      <p class="rubric">
        Fifty years on, Singapore’s ruling party looks secure        <span data-href-redirect="/node/21662410/comments#comments" class="comment-icon"><span>23</span></span>      </p>
    </div>
      </a>
</article>
  
  
    <div class="package-more"><a href="/sections/asia" title="Asia">More in Asia &raquo;</a></div>
  </section><section class="news-package typog-package">

  <h1 class="fly-title">China&#039;s economy</h1>

    <article>
  <a href="/blogs/freeexchange/2015/08/chinas-stockmarket" >
          <div>
          <h2 class="headline">The government giveth and taketh away</h2>
      <p class="rubric">
        China reduces its direct interventions in wobbly markets, but the central bank aims a monetary boost at the broader economy        <span data-href-redirect="/blogs/freeexchange/2015/08/chinas-stockmarket#comments" class="comment-icon"><span>66</span></span>      </p>
    </div>
      </a>
</article>
  
    <ul><li class="first"><div class="">
  <a href="https://espresso.economist.com/376187c6d721478403b2d6a3f4aedd05?utm_content=buffer355f4&amp;utm_medium=social&amp;utm_source=twitter.com&amp;utm_campaign=buffer" class="headline">Spreading the pain: where a Chinese slowdown really hurts</a>  </div>
</li>
<li class="even last"><div class="">
  <a href="/blogs/graphicdetail/2015/08/daily-chart-9" class="headline">Daily chart: The gravity of China’s great fall</a>  <a href="/blogs/graphicdetail/2015/08/daily-chart-9#comments" title="Comments" class="comment-icon"><span>4</span></a></div>
</li>
</ul>  
    <div class="package-more"><a href="/sections/china" title="China">More in China &raquo;</a></div>
  </section><section class="news-package typog-package">

  <h1 class="fly-title">American oil</h1>

    <article>
  <a href="/news/finance-and-economics/21661673-long-overdue-easing-protectionist-export-ban-nafta-naphtha" >
          <div>
          <h2 class="headline">Nafta naphtha</h2>
      <p class="rubric">
        A long-overdue easing of a protectionist export ban        <span data-href-redirect="/node/21661673/comments#comments" class="comment-icon"><span>7</span></span>      </p>
    </div>
      </a>
</article>
  
  
    <div class="package-more"><a href="/sections/business-finance" title="Business &amp;amp; finance">More in Business &amp; finance &raquo;</a></div>
  </section><section class="news-package typog-package">

  <h1 class="fly-title">Africa&#039;s worst war</h1>

    <article>
  <a href="/news/middle-east-and-africa/21662478-disagreeing-agree-south-sudan-agrees-peace-deal-unlikely-last" >
          <div>
          <h2 class="headline">Disagreeing to agree</h2>
      <p class="rubric">
        South Sudan signs a peace deal that is unlikely to last        <span data-href-redirect="/node/21662478/comments#comments" class="comment-icon"><span>1</span></span>      </p>
    </div>
      </a>
</article>
  
  
    <div class="package-more"><a href="/world/middle-east" title="Middle East &amp;amp; Africa">More in Middle East &amp; Africa &raquo;</a></div>
  </section>      </div>
</div>

<div class="side-box multiple-box grey-palette typog-highlights grid-3 pull-7">
  <div id="homepage-highlight-1"><article >
  <a href="/news/britain/21662591-net-migration-britain-has-never-been-higher-immigration-breaks-record">
    <div>
                    <h1 class="fly-title">Immigration breaks a record</h1>
                          <p class="rubric">David Cameron promised to reduce net migration “from the hundreds of thousands to the tens of thousands”. Thankfully, he failed.
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150829_brc577.png" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="136" />      </a>
</article></div><div id="homepage-highlight-2"><article >
  <a href="/blogs/prospero/2015/08/british-theatre">
    <div>
                    <h1 class="fly-title">Adrian Noble</h1>
                          <p class="rubric">A chat about Oscar Wilde and British theatre
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150829_BKP503_168.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="95" />      </a>
</article></div><div id="homepage-highlight-3"><article >
  <a href="/blogs/erasmus/2015/08/europes-religious-war">
    <div>
                    <h1 class="fly-title">Religion and debt</h1>
                          <p class="rubric">A French minister is the latest to explain the euro-crisis in religious terms
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150829_BLP508_168.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="95" />      </a>
</article></div><div id="homepage-highlight-4"><article >
  <a href="/blogs/democracyinamerica/2015/08/down-syndrome">
    <div>
                    <h1 class="fly-title">Down syndrome</h1>
                          <p class="rubric">A controversial abortion bill puts moderate Republican presidential candidate John Kasich in an awkward spot
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150829_USP501_168.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="95" />      </a>
</article></div><div id="homepage-highlight-5"><article >
  <a href="/news/science-and-technology/21662365-scientists-are-developing-jab-might-only-need-be-given-once-lifetime-why-universal">
    <div>
                    <h1 class="fly-title">The &quot;universal&quot; flu vaccine</h1>
                          <p class="rubric">Scientists are developing a jab that might only need to be given once in a lifetime
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150829_BLP507_168.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="95" />      </a>
</article></div><div id="homepage-highlight-6"><article >
  <a href="/news/united-states/21661815-obscure-dispute-about-rubbish-could-reshape-agency-working-whos-boss">
    <div>
                    <h1 class="fly-title">Companies and employment</h1>
                          <p class="rubric">An obscure dispute about rubbish could reshape agency working
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150822_USD001_168.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="95" />      </a>
</article></div><div id="homepage-highlight-7"><article >
  <a href="http://test.economist.com/sciencebriefs">
    <div>
                    <h1 class="fly-title">Science briefs</h1>
                          <p class="rubric">The missing 95% of the universe
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/sci_320.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="119" />      </a>
</article></div><div id="homepage-highlight-8"><article >
  <a href="/news/middle-east-and-africa/21661826-costly-valuable-lessons-guerrilla-army-once-fought">
    <div>
                    <h1 class="fly-title">Hizbullah’s learning curve</h1>
                          <p class="rubric">Costly but valuable lessons for a guerrilla army that once fought in the shadows
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150822_MAP003_168.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="95" />      </a>
</article></div></div>
<!-- End homepage node -->              </div> <!-- /#main-area -->

              <div id="column-right" class="grid-6 clearfix">
          <div id="homepage-touts">
  <div class="grid-3 grid-first">
    <div class="tout tout-1">
          </div>

    <div class="tout tout-2">
          </div>
  </div> <!-- /#homepage-touts -->

  <div class="grid-3">
    <div class="tout tout-issue">
      <div class="cover-image-container">
              </div>
    </div>
  </div> <!-- /#homepage-issue -->
</div>
          <div id="block-ec_ads-mpu_no_wrapper_top_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div class="ec-ads ec-ads-remove-if-empty ec-ads-remove-wrapper"><!-- Site: Web.  Zone: Home |  --> <div id="gpt_mpu_no_wrapper_top_ad" data-cb-ad-id="Mpu no wrapper top ad">



</div></div>  </div>
</div>

<div id="block-ec_ads-top_mpu_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div class="ec-ads ec-ads-remove-if-empty"><p class="ec-ads-label">Advertisement</p><!-- Site: Web.  Zone: Home |  --> <div id="gpt_top_mpu_ad" data-cb-ad-id="Top mpu ad">

    </div></div>  </div>
</div>

<div id="block-ec_ads-ribbon_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div class="ec_topic_ribbon ec-ads-remove-if-empty"><!-- Site: Web.  Zone: Home |  --> <div id="gpt_ribbon_ad" data-cb-ad-id="Ribbon ad">
    </div></div>  </div>
</div>

<div id="block-ec_social-right_rail_social_share_buttons" class="block block-ec_social 
">
    <div class="content clearfix">
    <div id="social-share-buttons-block">
  <div class="title">Follow <cite>The Economist</cite></div>
  <div class="social-share-buttons">
    <ul class="clearfix">
      <li class="facebook">
        <a data-ec-omniture="rightrail|social_share|facebook" href="http://www.facebook.com/TheEconomist" title="Facebook" target="_blank">Facebook</a>
      </li>
      <li class="twitter">
        <a data-ec-omniture="rightrail|social_share|twitter" href="http://twitter.com/TheEconomist" title="Twitter" target="_blank">Twitter</a>
      </li>
      <li class="linked-in">
        <a data-ec-omniture="rightrail|social_share|linked-in" href="http://www.linkedin.com/groups/Economist-official-group-Economist-newspaper-3056216" title="Linked in" target="_blank">Linked in</a>
      </li>
      <li class="google-plus">
        <a data-ec-omniture="rightrail|social_share|plusone" href="https://plus.google.com/100470681032489535736/posts" title="Google plus" target="_blank">Google plus</a>
      </li>
      <li class="tumblr">
        <a data-ec-omniture="rightrail|social_share|tumblr" href="http://theeconomist.tumblr.com/" title="Tumblr" target="_blank">Tumblr</a>
      </li>
      <li class="instagram">
        <a data-ec-omniture="rightrail|social_share|instagram" href="http://instagram.com/theeconomist/" title="Instagram" target="_blank">Instagram</a>
      </li>
      <li class="youtube">
        <a data-ec-omniture="rightrail|social_share|youtube" href="http://www.youtube.com/user/economistmagazine" title="YouTube" target="_blank">YouTube</a>
      </li>
      <li class="rss">
        <a data-ec-omniture="rightrail|social_share|rss" href="/rss" title="RSS" target="_blank">RSS</a>
      </li>
      <li class="newsletters">
        <a data-ec-omniture="rightrail|social_share|newsletters" href="/newsletters" title="Newsletters" target="_blank">Newsletters</a>
      </li>
    </ul>
  </div>
</div>  </div>
</div>

<div id="block-ec_blogs-ec_blogs_block_recent" class="block block-ec_blogs 
">
    <div class="content clearfix">
    <div class="title">
                 <h6><a href="/latest-updates">Latest updates &raquo;</a></h6>
                 </div><div id="latest-updates"><article id="node-21662570" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/news/finance-economics/21662570-kingdom-can-stand-more-pain-it-will-take-much-cheaper-oil-saudi-arabia-take-action">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/08/articles/main/20150829_blp511.jpg" alt="Shia, not shale" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/news/finance-economics/21662570-kingdom-can-stand-more-pain-it-will-take-much-cheaper-oil-saudi-arabia-take-action"><span class="latest-updates-fly-title">Shia, not shale</span>: It will take much cheaper oil for Saudi Arabia to take...</a></p>
  <p class="dateline">
    <span class="section">Finance & Economics</span>
    <span class="timestamp" title="2015-08-27T13:21:14+00:00"> 37 mins ago</span>
  </p>
</article>

<article id="node-21662583" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/blogs/graphicdetail/2015/08/american-and-british-flight-safety-airlines-v-light-aircraft">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/08/blogs/graphic-detail/20150829_woc583.png" alt="American and British flight safety: airlines v light aircraft" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/blogs/graphicdetail/2015/08/american-and-british-flight-safety-airlines-v-light-aircraft"><span class="latest-updates-fly-title">American and British flight safety: airlines v light aircraft</span>: The perils...</a></p>
  <p class="dateline">
    <span class="section">Graphic detail</span>
    <span class="timestamp" title="2015-08-27T12:46:21+00:00"> 1 hrs 22 mins ago</span>
  </p>
</article>

<article id="node-21657486" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/blogs/graphicdetail/2015/08/daily-dispatches-0">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/08/blogs/graphic-detail/china.png" alt="Daily dispatches" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/blogs/graphicdetail/2015/08/daily-dispatches-0"><span class="latest-updates-fly-title">Daily dispatches</span>: China crisis</a></p>
  <p class="dateline">
    <span class="section">Graphic detail</span>
    <span class="timestamp" title="2015-08-27T12:12:48+00:00"> 1 hrs 55 mins ago</span>
  </p>
</article>

<article id="node-21662580" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/news/markets-and-data/21662580-retail-sales-producer-prices-wages-and-exchange-rates">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/08/articles/main/20150829_int600.png" alt="" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/news/markets-and-data/21662580-retail-sales-producer-prices-wages-and-exchange-rates"><span class="latest-updates-fly-title"></span>: Retail sales, producer prices, wages and exchange rates</a></p>
  <p class="dateline">
    <span class="section">Markets and data</span>
    <span class="timestamp" title="2015-08-27T10:15:56+00:00"> 3 hrs 47 mins ago</span>
  </p>
</article>

<article id="node-21662579" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/news/markets-and-data/21662579-foreign-reserves">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/08/articles/main/20150829_int500.png" alt="" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/news/markets-and-data/21662579-foreign-reserves"><span class="latest-updates-fly-title"></span>: Foreign reserves</a></p>
  <p class="dateline">
    <span class="section">Markets and data</span>
    <span class="timestamp" title="2015-08-27T10:15:09+00:00"> 3 hrs 34 mins ago</span>
  </p>
</article>

<article id="node-21624322" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/blogs/graphicdetail/2015/08/ebola-graphics">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/07/blogs/graphic-detail/20150829_wom999.png" alt="Ebola in graphics" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/blogs/graphicdetail/2015/08/ebola-graphics"><span class="latest-updates-fly-title">Ebola in graphics</span>: The toll of a tragedy</a></p>
  <p class="dateline">
    <span class="section">Graphic detail</span>
    <span class="timestamp" title="2015-08-27T09:23:55+00:00"> August 27th, 9:23</span>
  </p>
</article>

<article id="node-21662502" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/blogs/graphicdetail/2015/08/daily-chart-14">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/08/blogs/graphic-detail/20150829_woc578_0.png" alt="Daily chart" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/blogs/graphicdetail/2015/08/daily-chart-14"><span class="latest-updates-fly-title">Daily chart</span>: Who wants to live forever?</a></p>
  <p class="dateline">
    <span class="section">Graphic detail</span>
    <span class="timestamp" title="2015-08-27T00:01:10+00:00"> August 27th, 0:01</span>
  </p>
</article>

</div><div class="more-latest-updates"><a href="/latest-updates" class="more">More latest updates &raquo;</a></div>  </div>
</div>

<div id="block-ec_mostx-mostpopular" class="block block-ec_mostx 
">
    <div class="content clearfix">
    <div id="most-lists" class="block">
  <div id="block-title"><p>Most commented</p></div>
  <div class="list-wrapper">
          <ul id="commented-list" class="show">
        <li class="mostx-first"><a href="/news/europe/21661941-wanting-burden-shared-germany-eu-country-which-takes-most-asylum-seekers-straining?spc=scode&amp;spv=xm&amp;ah=9d7f7ab945510a56fa6d37c30b6f1709" id="track-commented-21661941"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/mostx_block/images/2015/08/articles/main/20150822_eup503.jpg" alt="Germany, the EU country which takes the most asylum seekers, is straining" title=""  class="mostx-image-first" width="168" height="95" /><span class="mostx-list">1</span><span class="mostx-text"><span class="mostx-fly-title">Refugees in Europe</span>Germany, the EU country which takes the most asylum seekers, is straining</span></a></li><li><a href="/news/europe/21662019-it-may-have-been-frances-latest-islamist-attack-time-no-one-was-killed-attempted-murder?spc=scode&amp;spv=xm&amp;ah=9d7f7ab945510a56fa6d37c30b6f1709" id="track-commented-21662019"><span class="mostx-list">2</span><span class="mostx-text"><span class="mostx-fly-title">Terrorism in France</span>: Attempted murder on the Paris express</span></a></li><li><a href="/news/europe/21661810-journey-capital-hinterland-shows-how-grim-life-has-become-and-how-russians?spc=scode&amp;spv=xm&amp;ah=9d7f7ab945510a56fa6d37c30b6f1709" id="track-commented-21661810"><span class="mostx-list">3</span><span class="mostx-text"><span class="mostx-fly-title">Russia’s economy</span>: The path to penury</span></a></li><li><a href="/news/international/21661812-islamic-states-revival-slavery-extreme-though-it-finds-disquieting-echoes-across?spc=scode&amp;spv=xm&amp;ah=9d7f7ab945510a56fa6d37c30b6f1709" id="track-commented-21661812"><span class="mostx-list">4</span><span class="mostx-text"><span class="mostx-fly-title">Islam and slavery</span>: The persistence of history</span></a></li><li><a href="/news/business-and-finance/21662092-china-sneezing-rest-world-rightly-nervous-causes-and-consequences-chinas?spc=scode&amp;spv=xm&amp;ah=9d7f7ab945510a56fa6d37c30b6f1709" id="track-commented-21662092"><span class="mostx-list">5</span><span class="mostx-text"><span class="mostx-fly-title">Market turmoil</span>: The causes and consequences of China's market crash</span></a></li>      </ul>
    
  </div>
</div>
  </div>
</div>

<div id="block-ec_ads-bottom_mpu_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div class="ec-ads ec-ads-remove-if-empty"><p class="ec-ads-label">Advertisement</p><!-- Site: Web.  Zone: Home |  --> <div id="gpt_bottom_mpu_ad" data-cb-ad-id="Bottom mpu ad">
    </div></div>  </div>
</div>

<div id="block-ec_pixel_tracking_onscroll-bottom_onscroll" class="block block-ec_pixel_tracking_onscroll 
">
    <div class="content clearfix">
    <div id="onscroll-ad-holder-mpu2"></div>  </div>
</div>

<div id="block-block-1" class="block block-block 
">
    <div class="content clearfix">
    <div id="product-events">
<div class="title">Products and events</div>
<div class="product-first linked">
<p class="products-events-section"><strong><a  class="social-link" href="/economist-quiz" target="_blank">Test your EQ</a></strong><br />
Take our weekly news quiz to stay on top of the headlines
</p>
</div>
<div class="product-section product-section-last linked">
<p class="products-events-section"><strong><a  class="social-link" href="http://econ.st/R7pQMy" target="_blank">Want more from <i>The Economist?</i></a></strong><br />
Visit The Economist e-store and you’ll find a range of carefully selected products for business and pleasure, Economist books and diaries, and much more</p>
</div>
</div>
  </div>
</div>

<div id="block-ec_ads-bottom_right_mpu_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div class="ec-ads ec-ads-remove-if-empty"><p class="ec-ads-label">Advertisement</p><!-- Site: Web.  Zone: Home |  --> <div id="gpt_bottom_right_mpu_ad" data-cb-ad-id="Bottom right mpu ad">

    </div></div>  </div>
</div>

<div id="block-ec_pixel_tracking_onscroll-bottom_right_onscroll" class="block block-ec_pixel_tracking_onscroll 
">
    <div class="content clearfix">
    <div id="onscroll-ad-holder"></div>  </div>
</div>

        </div> <!-- /#side-bar -->
      
      
    </div> <!-- /#columns -->
          <div id="footer-classifieds" class="clearfix">
        <div class="title">Classified ads</div>
        <div id="block-ec_ads-button_ad_1" class="block block-ec_ads 
ec-ads-gpt ec-classified-box ec-ads-classified first">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_button_ad_1" data-cb-ad-id="Button ad 1">

    </div>  </div>
</div>

<div id="block-ec_ads-button_ad_2" class="block block-ec_ads 
ec-ads-gpt ec-classified-box ec-ads-classified">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_button_ad_2" data-cb-ad-id="Button ad 2">

   </div>  </div>
</div>

<div id="block-ec_ads-button_ad_3" class="block block-ec_ads 
ec-ads-gpt ec-classified-box ec-ads-classified">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_button_ad_3" data-cb-ad-id="Button ad 3">
    </div>  </div>
</div>

<div id="block-ec_ads-button_ad_4" class="block block-ec_ads 
ec-ads-gpt ec-classified-box ec-ads-classified">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_button_ad_4" data-cb-ad-id="Button ad 4">
    </div>  </div>
</div>

<div id="block-ec_ads-button_ad_5" class="block block-ec_ads 
ec-ads-gpt ec-classified-box ec-ads-classified">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_button_ad_5" data-cb-ad-id="Button ad 5">

    </div>  </div>
</div>

<div id="block-ec_ads-button_ad_6" class="block block-ec_ads 
ec-ads-gpt ec-classified-box ec-ads-classified last">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_button_ad_6" data-cb-ad-id="Button ad 6"></div>  </div>
</div>

      </div>
        <aside class="site-index">
      <div>
  <div class="svg-logo"><img class="mh-logo" width="170" height="85" src="//cdn.static-economist.com/sites/all/themes/econfinal/images/svg/logo.svg" alt="The Economist" /></div>
      <ul class="site-index-1">
              <li><a href="/contact-info" data-ec-omniture="mini_map_home|contact_us">Contact us</a></li>
              <li><a href="/help/home" data-ec-omniture="mini_map_home|help">Help</a></li>
              <li><a href="/user" data-ec-omniture="mini_map_home|my_account">My account</a></li>
              <li><a href="/products/subscribe" data-ec-omniture="mini_map_home|subscribe">Subscribe</a></li>
              <li><a href="/printedition" data-ec-omniture="mini_map_home|print_edition">Print edition</a></li>
              <li><a href="/digital" data-ec-omniture="mini_map_home|digital_editions">Digital editions</a></li>
              <li><a href="/events-conferences" data-ec-omniture="mini_map_home|events">Events</a></li>
              <li><a href="http://jobs.economist.com/" data-ec-omniture="mini_map_home|jobs_economist_com">Jobs.Economist.com</a></li>
              <li><a href="/bookmarks" data-ec-omniture="mini_map_home|timekeeper_saved_articles">Timekeeper saved articles</a></li>
          </ul>
  </div>
<div>
  <h6>Sections</h6>
  <ul class="footer-index-site-2-content">
                <li><a href="/sections/united-states" data-ec-omniture="mini_map_home|united_states">United States</a></li>
                <li><a href="/sections/britain" data-ec-omniture="mini_map_home|britain">Britain</a></li>
                <li><a href="/sections/europe" data-ec-omniture="mini_map_home|europe">Europe</a></li>
                <li><a href="/sections/china" data-ec-omniture="mini_map_home|china">China</a></li>
                <li><a href="/sections/asia" data-ec-omniture="mini_map_home|asia">Asia</a></li>
                <li><a href="/sections/americas" data-ec-omniture="mini_map_home|americas">Americas</a></li>
                <li><a href="/sections/middle-east-africa" data-ec-omniture="mini_map_home|middle_east_africa">Middle East &amp; Africa</a></li>
                <li><a href="/sections/international" data-ec-omniture="mini_map_home|international">International</a></li>
                <li><a href="/sections/business-finance" data-ec-omniture="mini_map_home|business_finance">Business &amp; finance</a></li>
                <li><a href="/sections/economics" data-ec-omniture="mini_map_home|economics">Economics</a></li>
                <li><a href="/markets-data" data-ec-omniture="mini_map_home|markets_data">Markets &amp; data</a></li>
                <li><a href="/sections/science-technology" data-ec-omniture="mini_map_home|science_technology">Science &amp; technology</a></li>
                <li><a href="http://test.economist.com/specialreports" data-ec-omniture="mini_map_home|special_reports">Special reports</a></li>
                <li><a href="/sections/culture" data-ec-omniture="mini_map_home|culture">Culture</a></li>
                <li><a href="/multimedia" data-ec-omniture="mini_map_home|multimedia_library">Multimedia library</a></li>
            </ul>
      <h6>Debate and discussion</h6>
    <ul>
            <li><a href="http://test.economist.com/debate" data-ec-omniture="mini_map_home|the_economist_debates">The Economist debates</a></li>
                <li><a href="http://test.economist.com/letters" data-ec-omniture="mini_map_home|letters_to_the_editor">Letters to the editor</a></li>
                <li><a href="/economist-quiz" data-ec-omniture="mini_map_home|the_economist_quiz">The Economist Quiz</a></li>
      </ul>
</div>
  <div>
    <h6>Blogs</h6>
    <ul>

              <li><a href="/blogs/buttonwood" data-ec-omniture="mini_map_home|buttonwoods_notebook">Buttonwood&#039;s notebook</a></li>
              <li><a href="/blogs/democracyinamerica" data-ec-omniture="mini_map_home|democracy_in_america">Democracy in America</a></li>
              <li><a href="/blogs/erasmus" data-ec-omniture="mini_map_home|erasmus">Erasmus</a></li>
              <li><a href="/blogs/freeexchange" data-ec-omniture="mini_map_home|free_exchange">Free exchange</a></li>
              <li><a href="/blogs/gametheory" data-ec-omniture="mini_map_home|game_theory">Game theory</a></li>
              <li><a href="/blogs/graphicdetail" data-ec-omniture="mini_map_home|graphic_detail">Graphic detail</a></li>
              <li><a href="/blogs/gulliver" data-ec-omniture="mini_map_home|gulliver">Gulliver</a></li>
              <li><a href="/blogs/prospero" data-ec-omniture="mini_map_home|prospero">Prospero</a></li>
              <li><a href="/blogs/economist-explains" data-ec-omniture="mini_map_home|the_economist_explains">The Economist explains</a></li>
          </ul>
  </div>
  <div>
    <h6>Research and insights</h6>
    <ul>
                              <li><a href="/topics" data-ec-omniture="mini_map_home|topics">Topics</a></li>
                                        <li><a href="/economics-a-to-z" data-ec-omniture="mini_map_home|economics_a_z">Economics A-Z</a></li>
                                        <li><a href="/styleguide/introduction" data-ec-omniture="mini_map_home|style_guide">Style guide</a></li>
                                        <li><a href="http://www.theworldin.com/" data-ec-omniture="mini_map_home|the_world_in_2015">The World in 2015</a></li>
                                        <li><a href="/whichmba" data-ec-omniture="mini_map_home|which_mba_">Which MBA?</a></li>
                                        <li><a href="https://success.economist.com/?fsrc=econfooter" data-ec-omniture="mini_map_home|mba_services">MBA Services</a></li>
                                        <li><a href="https://gmat.economist.com/?gsrc=economist_footer&amp;c3ch=Economist&amp;c3nid=site footer" data-ec-omniture="mini_map_home|the_economist_gmat_tutor">The Economist GMAT Tutor</a></li>
                                        <li><a href="https://execed.economist.com/?fsrc=econ-foot" data-ec-omniture="mini_map_home|executive_education_navigator">Executive Education Navigator</a></li>
                                        <li><a href="http://test.economist.com/rights" data-ec-omniture="mini_map_home|reprints_and_permissions">Reprints and permissions</a></li>
                                        </ul>
              <h6><a href="http://www.economistgroup.com" data-ec-omniture="mini_map_home|the_economist_group">The Economist Group &raquo;</a></h6>
            <ul>
                                        <li><a href="http://www.eiu.com" data-ec-omniture="mini_map_home|the_economist_intelligence_unit">The Economist Intelligence Unit</a></li>
                                        <li><a href="http://store.eiu.com" data-ec-omniture="mini_map_home|the_economist_intelligence_unit_store">The Economist Intelligence Unit Store</a></li>
                                        <li><a href="http://www.corporatenetwork.com" data-ec-omniture="mini_map_home|the_economist_corporate_network">The Economist Corporate Network</a></li>
                                        <li><a href="http://ideaspeoplemedia.com/" data-ec-omniture="mini_map_home|ideas_people_media">Ideas People Media</a></li>
                                        <li><a href="http://www.moreintelligentlife.com" data-ec-omniture="mini_map_home|intelligent_life">Intelligent Life</a></li>
                                        <li><a href="http://www.rollcall.com/?t=0506EC&amp;p=econ&amp;s=econ" data-ec-omniture="mini_map_home|roll_call">Roll Call</a></li>
                                        <li><a href="http://www.cq.com/news.do" data-ec-omniture="mini_map_home|cq">CQ</a></li>
                                        <li><a href="http://www.eurofinance.com" data-ec-omniture="mini_map_home|eurofinance">EuroFinance</a></li>
                                        <li><a href="http://store.economist.com" data-ec-omniture="mini_map_home|the_economist_store">The Economist Store</a></li>
                                        </ul>
              <h6 class="minimap-site-index"><a href="/content/site-index" data-ec-omniture="mini_map_home|view_complete_site_index_">View complete site index »</a></h6>
            <ul>
                      </ul>
  </div>


    </aside>
    <div id="block-ec_ads-slider_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_slider_ad" data-cb-ad-id="Slider ad">

    </div>  </div>
</div>

<div id="block-ec_ads-adcast" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div class="ec-ads-remove-if-empty"><!-- Site: Web.  Zone: Home |  --> <div id="gpt_adcast" data-cb-ad-id="Adcast">
    </div></div>  </div>
</div>


    <footer>
        <div class="footer-stripe-top">
    <ul>
                          <li><a href="/contact-info" data-ec-omniture="footer_home|contact_us">Contact us</a></li>
                          <li><a href="http://test.economist.com/help" data-ec-omniture="footer_home|help">Help</a></li>
                          <li><a href="http://test.economist.com/help/about-us#About_Economistcom" data-ec-omniture="footer_home|about_us">About us</a></li>
                          <li><a href="http://www.economistgroupmedia.com" data-ec-omniture="footer_home|advertise_with_us">Advertise with us</a></li>
                          <li><a href="/mediadirectory" data-ec-omniture="footer_home|editorial_staff">Editorial Staff</a></li>
                          <li><a href="/mediadirectory/books" data-ec-omniture="footer_home|staff_books">Staff Books</a></li>
                          <li><a href="http://www.economistgroupcareers.com" data-ec-omniture="footer_home|careers">Careers</a></li>
                          <li><a href="/content/site-index" data-ec-omniture="footer_home|site_index">Site index</a></li>
                          </ul>
  </div>

<div class="footer-stripe-bottom">
  <ul>
    <li>Copyright &copy; The Economist Newspaper Limited 2015. All rights reserved.</li>
                        <li><a href="/help/accessibilitypolicy" data-ec-omniture="footer_home|accessibility">Accessibility</a></li>
                  <li><a href="http://www.economistgroup.com/results_and_governance/governance/privacy" data-ec-omniture="footer_home|privacy_policy">Privacy policy</a></li>
                  <li><a href="/cookies-info" data-ec-omniture="footer_home|cookies_info">Cookies info</a></li>
                  <li><a href="/legal/terms-of-use" data-ec-omniture="footer_home|terms_of_use">Terms of use</a></li>
                            <li class="cookie-pref"><div id="teconsent">

                            
                            </div></li>
        </ul>
</div>
    </footer> <!-- /footer -->
  </div>

 

    </body>
</html>`),
}
