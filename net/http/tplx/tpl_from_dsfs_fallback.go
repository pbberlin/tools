package tplx

// When there is not template file in the google datastore
// then use THIS:
var hugoTplFallback = []byte(`<!DOCTYPE html>
<html class="no-js" lang="en-US">
<head>
    <meta charset="utf-8">
    <meta name="viewport"     content="width=device-width, initial-scale=1.0">
    <meta name="description"  content="[REPLACE_DESC]">
    <meta name="author"       content="">
    <meta name="keyword"      content="[]">
	<link rel="shortcut icon" href="data:;base64,=">
	<link rel="icon"          href="data:;base64,=">
    <title>[REPLACE_TITLE]</title>
	<link href="/css/journal.css" rel="stylesheet">
 	<link href="/css/style.css" rel="stylesheet">
</head>
<body lang="en">

	<!-- 
		=============================================
		Fallback template!
		=============================================
	-->
	<span  id="FALLBACK when ds has no templates available" ></span>

	<div class="container">
		<div class="row" style='margin-top:12px'>
			<div class="navbar navbar-default navbar-inverse" role="navigation">
				<div class="navbar-header">
					<button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-responsive-collapse">
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
					</button>
					<a class="navbar-brand" href="/">An econ perspective on technology</a> 
				</div>
				<div class="navbar-collapse collapse navbar-responsive-collapse">
					<ul class="nav navbar-nav navbar-right">
						<li><a href="/">Home</a></li>
						
						<li> &nbsp; &nbsp; &nbsp; </li>
					</ul>
				</div>
			</div>
		</div>
	</div>




	<div class="container">	
		<div class="row">
			<div class="col-md-offset-1 col-md-10">
				[REPLACE_DESC]
				
			</div>
		</div>
		<div class="row">
			<div class="col-md-offset-1 col-md-10">
				<p>[REPLACE_CONTENT]</p>

			</div>
		</div>
	</div>



	<div class="container">
	  <div class="row col-md-12">
        <footer>
 			<span id='REPLACE_FOOTER'></span>
          <div>
                <p>
                    &copy;2015 
                    &nbsp; &nbsp; &nbsp; 
                    <!--
					<a href="/tech-news">Tech-News</a>  
					-->
                </p>
            </div>
            </footer>
        </div>
    </div>


    </body>
</html>`)
