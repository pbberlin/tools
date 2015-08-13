package fetch_rss

var testDocs = make([]string, 2)

func init() {

	testDocs[0] = `<!DOCTYPE html><html><head>
		<script type="text/javascript" src="./article01_files/empty.js"></script>
		<link href="./article01_files/vendor.css" rel="stylesheet" type="text/css"/>
		</head><body><p>Links:
				<span>span01</span>
				<span>span02-line1<br>span02-line2</span>
				<span>span03</span>
			</p>
			<style> p {font-size:17px}</style>
			<ul>
				<li id='332' ><a   href="/some/first/page.html">Linktext1 <span>inside</span></a>
				<li><a   href="/snd/page" title="wi-title">LinkT2</a>
			</ul>
			<div>
				<div>div-1-content</div>
				<div>div-2-content</div>
				<p>pararaph in between</p>
				<div>div-3-content with iimmage<img alt="alt-cnt" title='title-cnt' 
				href='some-long-href-some-long-href-some-long-href-some-long-href'>after img</div>
			</div>


				<menu type="popup" typeFirefox="context"  id="myMenu">
				    <menuitem label="Item 1" title="First menu item"> </menuitem>
				    <menuitem label="Item 2" title="Second menu item" disabled> </menuitem>
				</menu>

				<div contextmenu="myMenu">Right-Click Me</div>
				<a contextmenu="myMenu" href='#' >A link with custom context menu</a>

				<script>
					// changing type attr for firefox
					var menus = document.querySelectorAll('menu');
					Array.prototype.forEach.call(menus, function(menu) {
					if (menu.getAttribute('type') == 'popup' && menu.type == 'list') {
						menu.type = 'context';
					}
					});
				</script>

			</body></html>`

	testDocs[1] = `	<p>
				Ja so sans<br/>
				Ja die sans.
			</p>
			<ul>
				<li>Gepriesen sei der Mann, der wenn er nichts</li>
				<li>zu sagen hat</li>
			</ul>`

	// write out to http doc root
	// for i := 0; i < len(testDocs); i++ {
	// 	fn := fmt.Sprintf(docRoot+"/%v/art%02v.html", hosts[0], i)
	// 	osutilpb.Bytes2File(fn, []byte(testDocs[i]))
	// }

}
