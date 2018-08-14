package main

import ()

type Page struct {
	Content string
}

func pageBuilder() string {
	pageReturn := `
	<html lang="en">
	<head>
	  <meta charset="UTF_8">
	  <meta http-equiv="refresh" content="5">
	  <title>EPS-CONDUIT</title>
	  <link rel="stylesheet" href="resources/css/lcars.css">
	  <style>

	  html, body { background: black }
	  p, h1, h2, h3 {
		margin-top: 1em;
	  }

	  pre {
		  display: inline;
		    white-space: pre-line;
		    word-wrap: break-word;
	  }
	  </style>
	</head>
  <body>
    <div class="lcars-app-container">
	<!-- HEADER================================== -->
	<div id="header" class="lcars-row header">
	  <!-- ELBOW -->
	  <div class="lcars-elbow left-bottom lcars-tan-bg"></div>

	  <!-- BAR WITH TITLE -->
	  <div class="lcars-bar horizontal">
	    <div class="lcars-title right">eps-conduit ` + VERSION + `</div>
	  </div>

	  <!-- ROUNDED EDGE DECORATED -->
	  <div class="lcars-bar horizontal right-end decorated"></div>
	</div>

	<!-- SIDE MENU ================== -->

	<div id="left-menu" class="lcars-column start-space lcars-u-1">

	  <!-- FILLER -->
	  <div class="lcars-bar lcars-u-1"></div>
	</div>

	<!-- FOOTER ============================ -->

	<div id="footer" class="lcars-row ">
	  <!-- ELBOW -->
	  <div class="lcars-elbow left-top lcars-tan-bg"></div>
	  <!-- BAR -->
	  <div class="lcars-bar horizontal both-divider bottom"></div>
	  <!-- ROUNDED EDGE -->
	  <div class="lcars-bar horizontal right-end left-divider bottom"></div>
	</div>

	<!-- MAIN CONTAINER -->
	<div id="container">
	  <!-- COLUMN LAYOUT -->
	  <div class="lcars-column lcars-u-5">

	    {{ .Content }}

	  </div>
	</div>
  </div>

  </body>
  </html>
	`
	return pageReturn
}
