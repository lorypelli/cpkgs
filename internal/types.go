package internal

type JSON struct {
	Schema        string         `json:"$schema"`
	Language      string         `json:"language"`
	Compiler      string         `json:"compiler"`
	FileName      string         `json:"file_name"`
	Include       Include        `json:"include"`
	CPPExtensions *CPPExtensions `json:"cpp_extensions,omitempty"`
}

type Include struct {
	C   []string `json:"C"`
	CPP []string `json:"C++"`
	H   []string `json:"H"`
	HPP []string `json:"H++"`
}

type CPPExtensions struct {
	Code   string `json:"code,omitempty"`
	Header string `json:"header,omitempty"`
}
