package ast

// chunk ::= block
// type Chunk *Block

// block ::= {stat} [retstat]
// retstat ::= return [explist] [‘;’]
// explist ::= exp {‘,’ exp}
type Block struct {
	LastLine int // todo
	Stats    []Stat
	RetStat  *RetStat
}

type RetStat struct {
	Line     int // ?
	LastLine int // ?
	ExpList  []Exp
}
