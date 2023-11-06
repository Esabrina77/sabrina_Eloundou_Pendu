package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func clearTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")

	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	clearTerminal()

	//Ouverture du fichier
	file, err := os.Open("mots.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Lecture des mots du dictionnaire
	scanner := bufio.NewScanner(file)
	var mots []string
	for scanner.Scan() {
		mots = append(mots, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Initialisation du générateur de nombres aléatoires
	rand.Seed(time.Now().UnixNano())

	//nom du joueur
	var nom string
	for {

		fmt.Println()
		fmt.Println("       **************Bienvenue dans Jeu du Pendu**********    ")
		fmt.Println()
		fmt.Println("        ************* VEUILLEZ ENTRER VOTRE NOM **********      ")
		fmt.Scanln(&nom)

		if nom == "" {
			fmt.Printf("*************  VEUILLEZ ENTRER VOTRE NOM!! ***********\n")
			clearTerminal()
			continue
		} else {
			break
		}
	}

	clearTerminal()
	fmt.Printf("Bienvenue à toi %s! \n", nom)
	fmt.Println("Choisissez un niveau de difficulté :")
	fmt.Println("1. Facile")
	fmt.Println("2. Moyen")
	fmt.Println("3. Difficile")
	fmt.Println("4. GOLD LEVEL")
	// Kaporal game- CHOIX

	var choix int
	fmt.Scanln(&choix)
	var mot string
	switch choix {
	case 1:
		mot = choisirMot(mots, 3, 6)
	case 2:
		mot = choisirMot(mots, 6, 8)
	case 3:
		mot = choisirMot(mots, 8, 15)
	case 4:
		mot = choisirMot(mots, 15, 40)
	default:
		fmt.Println("Choix invalide")
		return
	}
	clearTerminal()

	vie := 9
	Score := 0
	Punition := 2
	lettresDevinees := make([]bool, len(mot))
	LettreAleatoires(lettresDevinees)
	lettresDevinees[0] = true
	lettresProposees := make(map[string]bool)

	// OUVRIR LE FICHIER DU BONHOMME PENDU
	penduFile, err := os.Open("Hangman.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer penduFile.Close()
	var pendu []string
	scanner = bufio.NewScanner(penduFile)
	bloc := ""
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			pendu = append(pendu, bloc)
			bloc = ""
			continue
		}
		bloc += line + "\n"
	}
	var penduInverse []string
	for i := len(pendu) - 1; i >= 0; i-- {
		penduInverse = append(penduInverse, pendu[i])
	}
	pendu = penduInverse
	reader := bufio.NewReader(os.Stdin)
	for vie > 0 {
		fmt.Println()
		fmt.Println("Mot : ", afficherMot(mot, lettresDevinees))
		fmt.Println("Vies restantes : ", vie)
		fmt.Println("Entrez une lettre : ")

		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input)

		if len(input) == 0 {
			fmt.Println(" Veuillez Entrez une lettre !! ")
			continue
		}

		if len(input) == 1 {
			lettre := strings.ToLower(input)
			if lettreDejaProposee(lettre, lettresProposees) {
				fmt.Println("lettre déjà proposée. Veuillez saisir une nouvelle lettre.")
				continue
			}

			if strings.Contains(mot, lettre) {
				lettresProposees[lettre] = true
				for i, l := range mot {
					if string(l) == lettre {
						lettresDevinees[i] = true
					}
				}

				fmt.Println()
				fmt.Println("lettre correcte")
				Score++

			} else {
				fmt.Println("lettre incorrecte")
				vie--
				afficherPendu(pendu, vie)
			}

		} else if len(input) > 1 {
			if input == mot {
				lettresDevinees = make([]bool, len(mot))

				fmt.Println()
				fmt.Printf("Bien joué! joueur %s vous avez deviné le mot %s:", nom, mot)
				Score += 5
				fmt.Println()
				fmt.Println("Votre score FINAL est :", Score)
				return
			} else {
				vie -= Punition
				Score -= 2
				if Score < 0 {
					Score = 0
				}
				afficherPendu(pendu, vie)
				fmt.Println("LE MOT PRPOPOSE EST INCORRECT!!")
			}
		}

		if motComplet(lettresDevinees) {
			fmt.Println()
			fmt.Printf("Bien joué! %s vous avez deviné le mot '%s':\n ", nom, mot)

			break
		}

		if vie == 0 {
			fmt.Println()
			fmt.Printf("Désolé %s! vous avez perdu. Le mot était '%s':\n ", nom, mot)
			fmt.Println("Votre score FINAL est :", Score)
			return
		}
		fmt.Println("Votre score actuel est :", Score)
	}
	fmt.Println("Votre score FINAL est :", Score)
}

func lettreDejaProposee(lettre string, lettresproposees map[string]bool) bool {
	for k := range lettresproposees {
		if strings.ContainsRune(k, []rune(lettre)[0]) {
			return true
		}
	}
	return false
}

// Sélection  d'un mot de façon aléatoire dans le dictionnaire(en fonction de la longueur du mot)
func choisirMot(mots []string, minLen, maxLen int) string {
	var motsFiltres []string
	for _, mot := range mots {

		if len(mot) >= minLen && len(mot) <= maxLen {
			motsFiltres = append(motsFiltres, mot)
		}
	}

	if len(motsFiltres) == 0 {
		log.Fatal("Aucun mot disponible dans cette plage de longueurs")
	}

	index := rand.Intn(len(motsFiltres))
	return motsFiltres[index]
}

// choisir une lettre aléatoire à afficher dans le mot excepté la première lettre
// la premiere lettre du mot est marqué comme étant déjà deviné.
func LettreAleatoires(mot []bool) {
	mot[rand.Intn(len(mot)-2)+1] = true
}

// masquer les lettres non devinées du mot
func afficherMot(mot string, lettresDevinees []bool) string {
	var affichage string
	for i, l := range mot {
		if lettresDevinees[i] {
			affichage += string(l) + " "
		} else {
			affichage += "_"
		}
	}
	return affichage
}

// vérifier si le mot a été totalement trouvé par le joueur
func motComplet(lettresDevinees []bool) bool {
	for _, devinee := range lettresDevinees {
		if !devinee {
			return false
		}
	}
	return true
}

// Afficher le bonhomme pendu
func afficherPendu(pendu []string, vie int) {
	if vie < len(pendu) {
		fmt.Println(pendu[vie])
	} else {
		fmt.Println(pendu[len(pendu)-1])
	}
}
