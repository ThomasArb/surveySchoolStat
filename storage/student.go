package storage

import "statSurvey/config"

type Student struct {
	Questions [config.NbQuestions]uint
}
