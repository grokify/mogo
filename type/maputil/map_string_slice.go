package maputil

type MapStringSlice map[string][]string

func (mss MapStringSlice) Add(key, val string) {
	if _, ok := mss[key]; !ok {
		mss[key] = []string{}
	}
	mss[key] = append(mss[key], val)
}
