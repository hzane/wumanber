package wumanber

// from https://code.google.com/wu-manber-chinese

type wumanber struct {
	pat_mlength int
	block_n     int
	table_size  int
	//	shift       map[int]int
	//	prefix      map[int][]prefix_pair
	shift    []int
	prefix   [][]struct{ hash, id int }
	patterns [][]rune
}

func hash(chars []rune) (h int) {
	for _, c := range chars {
		h = int(c) + h<<6 + h<<16 - h
	}
	return 0
}

func New(patterns []string) *wumanber {
	v := &wumanber{
		patterns:    make([][]rune, len(patterns)),
		pat_mlength: 0xffffffff, // large enough
		block_n:     3,
	}
	for i, pat := range patterns {
		v.patterns[i] = []rune(pat)
		len := len(v.patterns[i])
		if v.pat_mlength > len {
			v.pat_mlength = len
		}
	}
	if v.block_n > v.pat_mlength {
		v.block_n = v.pat_mlength
	}

	var primes = []int{66047, 263167, 16785407}
	v.table_size = primes[2]
	for _, prime := range primes {
		if prime > v.pat_mlength*10*len(patterns) {
			v.table_size = prime
			break
		}
	}
	v.shift = make([]int, v.table_size)
	for i, _ := range v.shift {
		v.shift[i] = v.pat_mlength + 1 - v.block_n
	}

	v.prefix = make([][]struct{ hash, id int }, v.table_size)

	for id, pat := range v.patterns {
		for i := v.pat_mlength; i >= v.block_n; i = i - 1 {
			h := hash(pat[i-v.block_n:i]) % v.table_size
			if v.shift[h] > v.pat_mlength-i {
				v.shift[h] = v.pat_mlength - i
			}
			if i == v.pat_mlength {
				ph := hash(pat[:v.block_n]) % v.table_size
				v.prefix[h] = append(v.prefix[h], struct{ hash, id int }{ph, id})
			}
		}
	}
	return v
}

func (wm *wumanber) search(text []rune) (v []int) {
	var index = wm.pat_mlength - 1
	for index < len(text) {
		h := hash(text[index-wm.block_n+1:index+1]) % wm.table_size
		shift := wm.shift[h]
		if shift > 0 {
			index = index + shift
			continue
		}
		phidx := index - wm.pat_mlength + 1
		ph := hash(text[phidx:phidx+wm.block_n]) % wm.table_size
		prefixes := wm.prefix[ph]
		for _, prefix := range prefixes {
			if prefix.hash != ph {
				continue
			}
			pat := wm.patterns[prefix.id]
			if rune_has_prefix(text[phidx:], pat) {
				v = append(v, prefix.id)
			}
		}
		index++
	}
	return
}

func (wm *wumanber) Search(text string) (v []int) {
	return wm.search([]rune(text))
}

func rune_has_prefix(r []rune, prefix []rune) bool {
	if len(r) < len(prefix) {
		return false
	}
	for idx, p := range prefix {
		if r[idx] != p {
			return false
		}
	}
	return true
}
