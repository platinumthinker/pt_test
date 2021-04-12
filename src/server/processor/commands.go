package processor

type Operation interface {
	Eval([]int) (int, error)
}

type SumOperation struct{}

func (_ *SumOperation) Eval(args []int) (int, error) {
	sum := 0
	for _, arg := range args {
		sum += arg
	}

	return sum, nil
}

type MulOperation struct{}

func (_ *MulOperation) Eval(args []int) (int, error) {
	mul := 1
	for _, arg := range args {
		mul *= arg
	}

	return mul, nil
}
