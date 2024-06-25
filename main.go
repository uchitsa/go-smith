package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	minStatements = 5
	maxStatements = 20
)

var (
	operators = []string{"+", "-", "*", "/"}
	types     = []string{"int", "float64", "string"}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		program := generateRandomProgram()
		programFile := fmt.Sprintf("random_program_%d.go", i)
		err := os.WriteFile(programFile, []byte(program), 0644)
		if err != nil {
			fmt.Println("Error writing program to file:", err)
			continue
		}

		output1, err1 := compileAndRun(programFile, "go")
		output2, err2 := compileAndRun(programFile, "gccgo")

		if err1 != nil || err2 != nil {
			fmt.Println("Compilation or execution error:", err1, err2)
			continue
		}

		if string(output1) != string(output2) {
			fmt.Println("Differential testing detected a discrepancy!")
			fmt.Println("Program:", program)
			fmt.Println("Go output:", string(output1))
			fmt.Println("Gccgo output:", string(output2))
		} else {
			fmt.Println("No discrepancy detected.")
		}

		os.Remove(programFile)
	}
}

func generateRandomProgram() string {
	statements := rand.Intn(maxStatements-minStatements+1) + minStatements
	var buf strings.Builder

	buf.WriteString("package main\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\"fmt\"\n")
	buf.WriteString(")\n\n")
	buf.WriteString("func main() {\n")

	for i := 0; i < statements; i++ {
		buf.WriteString(generateRandomStatement())
	}

	buf.WriteString("}\n")
	return buf.String()
}

func generateRandomStatement() string {
	var buf strings.Builder
	typ := types[rand.Intn(len(types))]
	op := operators[rand.Intn(len(operators))]

	switch typ {
	case "int":
		buf.WriteString(fmt.Sprintf("a%d := %d %s %d\n", rand.Intn(100), rand.Intn(100), op, rand.Intn(100)))
	case "float64":
		buf.WriteString(fmt.Sprintf("a%d := %f %s %f\n", rand.Intn(100), rand.Float64()*100, op, rand.Float64()*100))
	case "string":
		buf.WriteString(fmt.Sprintf("a%d := \"%s\" + \"%s\"\n", rand.Intn(100), randomString(5), randomString(5)))
	}

	buf.WriteString(fmt.Sprintf("fmt.Println(a%d)\n", rand.Intn(100)))
	return buf.String()
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func compileAndRun(programFile, compiler string) ([]byte, error) {
	var cmd *exec.Cmd
	if compiler == "go" {
		cmd = exec.Command("go", "run", programFile)
	} else if compiler == "gccgo" {
		cmd = exec.Command("gccgo", programFile, "-o", "output", "-static-libgo")
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
		cmd = exec.Command("./output")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
