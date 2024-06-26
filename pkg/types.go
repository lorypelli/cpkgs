package pkg

type JSON struct {
	Compiler string  `json:"compiler"`
	FileName string  `json:"file_name"`
	Include  Include `json:"include"`
}

type Include struct {
	C []string
	H []string
}
