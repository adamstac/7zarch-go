package archive

type sltCase struct {
	in   string
	want int
}

// Sample minimal -slt outputs for tests
type filesCountCase struct {
	in   string
	want int
}

var sltCases = []sltCase{
	{in: "Path = file1.txt\nSize = 10\n\nPath = dir/file2.bin\nSize = 20\n", want: 2},
	{in: "Some header\n\nPath = one\n\nPath = two\n\nPath = three\n", want: 3},
	{in: "No paths here", want: 0},
}

var fileCountCases = []filesCountCase{
	{in: "Files: 15", want: 15},
	{in: "15 files, 2048576 bytes", want: 15},
	{in: "files", want: 0},
}

func Example_countPathsInSlt() {
	for _, c := range sltCases {
		_ = countPathsInSlt(c.in)
	}
	// Output:
}
