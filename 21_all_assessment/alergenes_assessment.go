package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	defer timeTrack(time.Now(), "Solving Time")
	inputFile := flag.String("f", "this is not optional!", "the input file")
	flag.Parse()

	fmt.Printf("==> Solving for Both Problems (%v):\n", *inputFile)
	dat, _ := ioutil.ReadFile(*inputFile)
	lines := strings.Split(string(dat), "\n")
	ingredientsAppearance := make(map[string]int)

	allergenePotentialCausesMap := make(map[string][]string)
	for _, line := range lines {
		ingredients, allergenes := ExtractAllergenesAndIngredients(line)
		UpdateIngredientAppearance(ingredients, &ingredientsAppearance)
		for _, allergene := range allergenes {
			if potentialCauses, exists := allergenePotentialCausesMap[allergene]; exists {
				allergenePotentialCausesMap[allergene] = Intersect(potentialCauses, ingredients)
			} else {
				allergenePotentialCausesMap[allergene] = ingredients
			}
		}
	}

	allergeneMap := make(map[string]string)

	for idedAllergene := IdentifyAllergene(allergenePotentialCausesMap, &allergeneMap); idedAllergene != ""; idedAllergene = IdentifyAllergene(allergenePotentialCausesMap, &allergeneMap) {
		for allergene, ingredients := range allergenePotentialCausesMap {
			if allergene == idedAllergene {
				continue
			}
			allergenePotentialCausesMap[allergene] = FilterItem(ingredients, allergeneMap[idedAllergene])
		}
	}

	nonAllergeneIngredientsAppearanceCount := 0
	for ingredient, count := range ingredientsAppearance {
		hasAllergen := false
		for _, ingr := range allergeneMap {
			if ingr == ingredient {
				hasAllergen = true
				break
			}
		}
		if !hasAllergen {
			nonAllergeneIngredientsAppearanceCount += count
		}
	}

	fmt.Println("Appearance Count of ingredients that cannot possibly contain any of the allergens", nonAllergeneIngredientsAppearanceCount)

	allergeneSlice := Keys(allergeneMap)
	sort.Strings(allergeneSlice)

	allergenicIngredientsSorted := make([]string, len(allergeneSlice))
	for i, allergene := range allergeneSlice {
		allergenicIngredientsSorted[i] = allergeneMap[allergene]
	}
	fmt.Println("Canonical Allergenic String is:")
	fmt.Println(strings.Join(allergenicIngredientsSorted, ","))
}

func Intersect(sliceA []string, sliceB []string) []string {
	intersection := []string{}
	for _, item := range sliceA {
		if FindIndex(sliceB, item) != -1 {
			intersection = append(intersection, item)
		}
	}
	return intersection
}
func FilterItem(ingredients []string, ref string) []string {
	index := FindIndex(ingredients, ref)
	if index == -1 {
		return ingredients
	}
	if index == 0 {
		return ingredients[1:]
	}
	if index == len(ingredients)-1 {
		return ingredients[0 : len(ingredients)-1]
	}
	return append(ingredients[:index], ingredients[index+1:]...)
}

func FindIndex(ingredients []string, ref string) int {
	for i, ingredient := range ingredients {
		if ref == ingredient {
			return i
		}
	}
	return -1
}

func IdentifyAllergene(allergenePotentialCausesMap map[string][]string, allergeneMap *map[string]string) string {
	for allergene, ingredients := range allergenePotentialCausesMap {
		if _, alreadyIdentified := (*allergeneMap)[allergene]; alreadyIdentified {
			continue
		}
		if len(ingredients) == 1 {
			(*allergeneMap)[allergene] = ingredients[0]
			return allergene
		}
	}
	return ""
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func ExtractAllergenesAndIngredients(line string) (ingredients []string, allergenes []string) {
	cleanedLine := strings.ReplaceAll(line, "(", "")
	cleanedLine = strings.ReplaceAll(cleanedLine, ")", "")
	cleanedLine = strings.ReplaceAll(cleanedLine, ",", "")
	parted := strings.Split(cleanedLine, "contains")
	ingredients = strings.Fields(parted[0])
	allergenes = strings.Fields(parted[1])
	return
}

func UpdateIngredientAppearance(ingredients []string, ingredientsAppearance *map[string]int) {
	for _, ingr := range ingredients {
		(*ingredientsAppearance)[ingr]++
	}
}

func Keys(myMap map[string]string) []string {
	keys := make([]string, len(myMap))
	i := 0
	for k, _ := range myMap {
		keys[i] = k
		i++
	}
	return keys
}
