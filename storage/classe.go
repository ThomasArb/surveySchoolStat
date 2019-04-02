package storage

import "statSurvey/statistics"

type Classe struct {
	Name      string
	School    string
	NbStudent uint
	Students  []Student
	Stats     statistics.StatClasse
}
