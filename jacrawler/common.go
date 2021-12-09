package jacrawler

func getParams(emb IEmbassyData, text string) (paramsMap map[string]string) {
	match := emb.GetRegex().FindStringSubmatch(text)

	paramsMap = make(map[string]string)
	for i, name := range emb.GetRegex().SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}
