package types

type BaseType struct {
	Name string
	Meta string
}

type CredType struct {
	BaseType
	Username string
	Password string
}

type TextType struct {
	BaseType
	Data string
}

type ByteType struct {
	BaseType
	Data []byte
}

type CardType struct {
	BaseType
	Data int64
}
