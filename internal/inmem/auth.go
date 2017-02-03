package inmem

func (i *InMem) isValidAuth(auth uint64) bool {
	for k := range i.auths {
		if k == auth {
			return true
		}
	}

	return false
}
