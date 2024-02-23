package rpc_service

type Calcolo float64

type Result float64

type Args struct {
	toConvert float64
}

func (s *Calcolo) LireToEuro(args Args, Result *float64) error {
	*Result = args.toConvert / 1936.27
	return nil
}

func (s *Calcolo) EuroToLire(args Args, Result *float64) error {
	*Result = (args.toConvert * (1936.27))
	return nil
}
