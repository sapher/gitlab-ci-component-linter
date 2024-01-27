package linter

type Linter struct {
	Workdir string
	Output  LinterOutput
	Rules   []LinterRuleFunc
}

func New(workDir string, output LinterOutput) *Linter {
	return &Linter{
		Workdir: workDir,
		Output:  output,
		Rules:   ruleFuncs,
	}
}

func (l *Linter) Execute() (LinterResults, error) {
	ruleResults := LinterResults{}
	for _, rule := range l.Rules {
		result, err := rule(l)
		if err != nil {
			return ruleResults, err
		}
		ruleResults = append(ruleResults, result)
	}
	return ruleResults, nil
}
