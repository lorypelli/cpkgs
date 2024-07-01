package pkg

type JSON struct {
	Schema               string               `json:"$schema"`
	Language             string               `json:"language"`
	Compiler             string               `json:"compiler"`
	FileName             string               `json:"file_name"`
	Include              Include              `json:"include"`
	AdditionalCPPOptions AdditionalCPPOptions `json:"additional_cpp_options"`
}

type Include struct {
	C   []string `json:"C"`
	CPP []string `json:"C++"`
	H   []string `json:"H"`
	HPP []string `json:"H++"`
}

type AdditionalCPPOptions struct {
	Extensions Extensions `json:"extensions"`
}

type Extensions struct {
	Code   string `json:"code"`
	Header string `json:"header"`
}
