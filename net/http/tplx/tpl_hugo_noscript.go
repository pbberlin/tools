package tplx

// Some content, such as the gitkit widget cannot stomach jQuery.
// Thus we need this additional template.
var HugoTplNoScript = `<!DOCTYPE html>
<html class="no-js" lang="en-US">

<head>
    <meta charset="utf-8">
    <meta name="viewport"     content="width=device-width, initial-scale=1.0">
    <meta name="description"  content="{{ .HtmlDescription }}">
    <meta name="author"       content="">
    <meta name="keyword"      content="[]">
	<link rel="shortcut icon" href="data:;base64,=">
	<link rel="icon"          href="data:;base64,=">
    
    <title>{{ .HtmlTitle }}</title>
 	
	<link href="/css/journal.css" rel="stylesheet">
 	<link href="/css/style.css" rel="stylesheet">

 	{{ .HtmlHeaders }}

</head>
<body lang="en">
	<div class="container">
		<div class="row" style='margin-top:12px'>
			<div class="navbar navbar-default navbar-inverse" role="navigation">
				<div class="navbar-header">
					<button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-responsive-collapse">
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
					</button>
					<a class="navbar-brand" href="/">{{ .ClaimTitle }}</a> 
				</div>
				<div class="navbar-collapse collapse navbar-responsive-collapse">
					<ul class="nav navbar-nav navbar-right">
						<li> &nbsp; &nbsp; &nbsp; </li>
					</ul>
				</div>
			</div>
		</div>
	</div>

	<div class="container">	
		<div class="row">
			<div class="col-md-offset-1 col-md-10">
				{{ .HtmlContent }}
			</div>
		</div>
	</div>

	<script>
		// No JQuery, that's the point
		// console.log("end of page " + $.fn.jquery + " Version")
	</script> 

    </body>
</html>`
