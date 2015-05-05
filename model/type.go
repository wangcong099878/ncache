package model

//定义子map表结构    一行 key=>val   key可以用于精确检索  val可以用于正则或条件检索
type Mii map[int]int //json不支持     map[int]int 转换
type Mss map[string]string
type Mis map[int]string
type Msi map[string]int

//db
type MiiDB map[string]Mii
type MssDB map[string]Mss
type MisDB map[string]Mis
type MsiDB map[string]Msi
