package utils

type Killable interface {
	Kill(msg string)
}

func Assert(assertion bool, msg string, k Killable) {
	if !assertion {
		if k == nil {
			panic(msg)
		}
		k.Kill(msg)
	}
}
