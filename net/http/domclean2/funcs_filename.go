package domclean2

import "fmt"

func weedoutFilename(articleId, weedoutStage int) (string, string) {
	fn := fmt.Sprintf("outp_%03v_%v.html", articleId, weedoutStage)
	prefix := fmt.Sprintf("outp_%03v", articleId)
	return fn, prefix
}
