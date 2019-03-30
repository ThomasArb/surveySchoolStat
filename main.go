package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"statSurvey/config"
	"statSurvey/storage"
)

func main() {
	classe := storeAClasseResults()
	fmt.Println(classe)
	jsonClasse , err := json.Marshal(classe)
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

func storeAClasseResults() storage.Classe{
	fmt.Print("Entrez le nom de la classe : ")
	reader := bufio.NewReader(os.Stdin)
	text,err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Combien d'élèves : ")
	var nbe uint8
	_, err = fmt.Scan(&nbe)
	if err != nil {
		log.Fatal(err)
	}
	classe := storage.Classe{}
	classe.Name = text[:len(text)-1]
	classe.NbStudent = nbe
	classe.Students = make([]storage.Student, classe.NbStudent)
	var i uint8
	for i = 0; i < classe.NbStudent; i++ {
		fmt.Printf("Saisir les notes pour l'élève n°%d :\n", i+1)
		classe.Students[i] = storage.Student{}
		for j := 0; j < config.NbQuestions; j++ {
			fmt.Printf("Note de la question %d : ", j+1)
			var note uint8
			_, err = fmt.Scan(&note)
			if err != nil{
				log.Fatal(err)
			} else if note > 10 {
				note = 42
			}
			classe.Students[i].Questions[j] = note
		}
	}
	return classe
}
