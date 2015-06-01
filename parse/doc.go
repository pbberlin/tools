package parse

// https://code.google.com/p/go/source/browse/html/node.go?repo=net

// example use
// 	parse.Fetch(1)
// 	parse.Tokenize()
// 	parse.ParseHtmlFiles()

//
// Todo: Build a stack
// opening tag => push; stack depth ++
// closing tag => pop ; stack depth --

// on opening tag *repeating*
//		=> repTag[depth]=tagname
//      add attribute repeat="true"
// on closing tag
//     if repTag[depth]
//     	delete  repTag[depth]
//      add attribute repeat="true"
