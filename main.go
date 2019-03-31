package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"statSurvey/config"
	"statSurvey/statistics"
	"statSurvey/storage"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	loop := true
	for loop {
		fmt.Println("Choisir une option :")
		fmt.Println("\t- new : créer une nouvelle classe")
		fmt.Println("\t- export : génére le tableau des statistiques")
		fmt.Println("\t- stop : stop le programme")
		input, _ := reader.ReadString('\n')
		switch input {
		case "new\n":
			classe := storeAClasseResults()
			createAllStatsForAClass(&classe)
			saveInJSON(&classe)
		case "export\n":
			classes := loadAllClasses()
			fmt.Println(classes)
		case "stop\n":
			loop = false
		}
	}

	/*
		classe := storeAClasseResults()
		fmt.Println(classe)
		createAllStatsForAClass(&classe)
		saveInJSON(&classe)
		classeLoad := loadAClass(classe.Name + ".json")
		fmt.Println(classeLoad)
	*/
}

func loadAllClasses() []storage.Classe {
	var classes = make([]storage.Classe, 0)
	dirname := "."
	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".json") {
			fmt.Println(file.Name())
			classes = append(classes, loadAClass(file.Name()))
		}
	}
	return classes
}

func loadAClass(fileName string) storage.Classe {
	classe := storage.Classe{}
	// Open our jsonFile
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &classe)
	defer jsonFile.Close()
	return classe
}

func createAllStatsForAClass(classe *storage.Classe) {
	classeStat := statistics.StatClasse{}
	createQuestionsStats(classe, &classeStat)
	createStudentsStats(classe, &classeStat)
	classe.Stats = classeStat
}

func createQuestionsStats(classe *storage.Classe, stats *statistics.StatClasse) {
	var i uint
	stats.StatQuestions = make([]statistics.StatQuestion, config.NbQuestions)
	for i = 0; i < config.NbQuestions; i++ {
		// Generate stat for a question
		var notHappy uint
		var happy uint
		var sum uint
		var j uint
		for j = 0; j < classe.NbStudent; j++ {
			note := classe.Students[j].Questions[i]
			sum += note
			if note >= 4 {
				happy++
			} else {
				notHappy++
			}
		}
		stats.StatQuestions[i].Average = float32(sum) / float32(classe.NbStudent)
		stats.StatQuestions[i].PercentageHigh = float32(happy) * 100.0 / float32(classe.NbStudent)
		stats.StatQuestions[i].PercentageLow = float32(notHappy) * 100.0 / float32(classe.NbStudent)
	}

}

func createStudentsStats(classe *storage.Classe, stats *statistics.StatClasse) {
	var i uint
	stats.StatStudents = make([]statistics.StatStudent, classe.NbStudent)
	for i = 0; i < classe.NbStudent; i++ {
		//Generate stat for a student
		student := classe.Students[i]
		var j uint
		for j = 0; j < config.NbQuestions; j++ {
			if j < 9 {
				stats.StatStudents[i].Sum1to9 += student.Questions[j]
			} else {
				stats.StatStudents[i].Sum10to19 += student.Questions[j]
			}
			stats.StatStudents[i].SumTotal += student.Questions[j]
		}
	}
}

func saveInCSV(classe *storage.Classe) {
}

func saveInJSON(classe *storage.Classe) {
	jsonClasse, err := json.Marshal(classe)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(classe.Name + ".json")
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(jsonClasse)
	if err != nil {
		log.Fatal(err)
	}
	f.Sync()
	f.Close()
}

func storeAClasseResults() storage.Classe {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Entrez le nom de la classe : ")
	className, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Entrez le nom de l'école : ")
	schoolName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Combien d'élèves : ")
	var nbe uint
	_, err = fmt.Scan(&nbe)
	if err != nil {
		log.Fatal(err)
	}
	classe := storage.Classe{}
	classe.Name = className[:len(className)-1]
	classe.School = schoolName[:len(schoolName)-1]
	classe.NbStudent = nbe
	classe.Students = make([]storage.Student, classe.NbStudent)
	var i uint
	for i = 0; i < classe.NbStudent; i++ {
		fmt.Printf("Saisir les notes pour l'élève n°%d :\n", i+1)
		classe.Students[i] = storage.Student{}
		for j := 0; j < config.NbQuestions; j++ {
			fmt.Printf("Note de la question %d : ", j+1)
			var note uint
			_, err = fmt.Scan(&note)
			if err != nil {
				log.Fatal(err)
			} else if note > config.MaxNote {
				note = 42
			}
			classe.Students[i].Questions[j] = note
		}
	}
	return classe
}
