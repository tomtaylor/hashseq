# Hasheq

Automatically map your primary key ints around hashids in go.

This is a helper package around [go-hashid](github.com/speps/go-hashids) which
implements database interfaces to automatically conver your table's PKs into
Hashids.

It also implements JSON interfaces to automatically marshal IDs into Hashids for
your API. 
