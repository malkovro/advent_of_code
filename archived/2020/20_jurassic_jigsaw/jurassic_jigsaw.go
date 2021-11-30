package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Image struct {
	Id      string
	Content [][]rune
	Size    int
	Sides   []int
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func SideToInt(side string) int {
	binarySideRep := strings.ReplaceAll(side, ".", "0")
	binarySideRep = strings.ReplaceAll(binarySideRep, "#", "1")
	intRep, err := strconv.ParseInt(binarySideRep, 2, 64)
	check(err)
	return int(intRep)
}

func main() {
	defer timeTrack(time.Now(), "Solving Time")
	inputFile := flag.String("f", "this is not optional!", "the input file")
	problem2Ptr := flag.Bool("problem2", false, "to solve the problem 2nd part")
	flag.Parse()

	problemNumber := 1
	if *problem2Ptr {
		problemNumber = 2
	}
	fmt.Printf("==> Solving for Problem %v (%v):\n", problemNumber, *inputFile)
	dat, _ := ioutil.ReadFile(*inputFile)
	lines := strings.Split(string(dat), "\n")

	images := []Image{}
	sidesMap := make(map[int][]string)

	image := Image{}
	var index int
	var leftSide []rune
	var rightSide []rune

	for _, line := range lines {
		if line == "" {
			index = 0
			continue
		}
		if strings.HasPrefix(line, "Tile") {
			id := strings.Split(strings.Split(line, " ")[1], ":")[0]
			image.Id = id
			continue
		}
		if index == 0 {
			image.Size = len(line)
			leftSide = []rune{[]rune(line)[0]}
			rightSide = []rune{[]rune(line)[len(line)-1]}
			image.Sides = []int{SideToInt(line)}
			image.Content = make([][]rune, len(line)-2)
			index++
			continue
		}

		leftSide = append(leftSide, rune(line[0]))
		rightSide = append(rightSide, rune(line[len(line)-1]))
		if index == image.Size-1 {
			image.Sides = append(image.Sides, SideToInt(string(rightSide)))
			image.Sides = append(image.Sides, reversedSideSignature(SideToInt(line), image.Size))
			image.Sides = append(image.Sides, reversedSideSignature(SideToInt(string(leftSide)), image.Size))
			images = append(images, image)
			for _, side := range image.Sides {

				ids, exists := sidesMap[side]
				if exists {
					sidesMap[side] = append(ids, image.Id)
				} else {
					sidesMap[side] = []string{image.Id}
				}

				reversedSide := reversedSideSignature(side, image.Size)
				ids, exists = sidesMap[reversedSide]
				if exists {
					sidesMap[reversedSide] = append(ids, image.Id)
				} else {
					sidesMap[reversedSide] = []string{image.Id}
				}
			}
		} else {
			image.Content[index-1] = []rune(line[1 : len(line)-1])
		}
		index++
	}

	corners := findCorners(images, sidesMap)
	if len(corners) != 4 {
		panic("Not the right amount of corners :(")
	}

	cornersMultiplied := 1
	for _, corner := range corners {
		id, _ := strconv.Atoi(corner.Id)
		cornersMultiplied *= id
	}
	fmt.Println("Corner images ids multiplied", cornersMultiplied)

	seaMapSize := int(math.Sqrt(float64(len(images))))
	seaMap := make([][]Image, seaMapSize)
	for j, _ := range seaMap {
		seaMap[j] = make([]Image, seaMapSize)
	}
	imageSet := ImageSet(images)
	tileMap := imageSet.ToMap()

	rotationCount := 0
	verticallyFlipped := false
	horizontallyFlipped := false
	var corner Image

	for _, candidate := range tileMap {
		for corner.Id == "" {
			if candidate.IsTopLeftCorner(sidesMap) {
				corner = candidate
			}
			err := candidate.Shuffle(&rotationCount, &verticallyFlipped, &horizontallyFlipped)
			if err != nil {
				break
			}
		}
	}
	seaMap[0][0] = corner

	cornerIndex := findIndex(images, corner)
	imageSet.RemoveAtIndex(cornerIndex)

	_, seaMap = RecursiveMapResolver(imageSet, seaMap, tileMap)

	fullSeaMap := ImageLite(SeaMap(seaMap).ToMap())

	seaMonstersCount := fullSeaMap.CountSeaMonsters()
	rotationCount = 0
	verticallyFlipped = false
	horizontallyFlipped = false
	var err error
	for seaMonstersCount == 0 && err == nil {
		err = fullSeaMap.Shuffle(&rotationCount, &verticallyFlipped, &horizontallyFlipped)
		seaMonstersCount = fullSeaMap.CountSeaMonsters()
	}
	// fullSeaMap.Print()
	if seaMonstersCount == 0 {
		check(err)
	}

	fmt.Printf("We removed %v sea monsters \n", seaMonstersCount)
	// fmt.Println(fullSeaMapImage.ToString())

	count := 0
	for _, line := range fullSeaMap {
		count += strings.Count(string(line), "#")
	}
	fmt.Printf("Roughness of the sea is: %v \n", count-15*seaMonstersCount)

}

func (img ImageLite) Print() {
	for _, line := range img {
		fmt.Println(string(line))
	}
}

var firstLineMonster = regexp.MustCompile("[.#]{18}#")
var secondLineMonster = regexp.MustCompile("#[.#]{4}##[.#]{4}##[.#]{4}###")
var thirdLineMonster = regexp.MustCompile("[.#](#[.#]{2}){6}[.#]")

func (image ImageLite) CountSeaMonsters() (count int) {
	for i, line := range image {
		if i == 0 || i == len(image)-1 {
			continue
		}
		if locs := secondLineMonster.FindAllStringIndex(string(line), -1); locs != nil {
			for _, loc := range locs {
				firstLineLoc := firstLineMonster.FindStringIndex(string(image[i-1][loc[0]:]))
				thirdLineLoc := thirdLineMonster.FindStringIndex(string(image[i+1][loc[0]:]))
				if firstLineLoc != nil && firstLineLoc[0] == 0 && thirdLineLoc != nil && thirdLineLoc[0] == 0 {
					count++
				}
			}
		}
	}
	return
}

type ImageSet []Image

func (s ImageSet) ToMap() map[string]Image {
	imageMap := make(map[string]Image)
	for _, image := range s {
		imageMap[image.Id] = image
	}
	return imageMap
}

func (s *ImageSet) RemoveAtIndex(index int) {
	if index == 0 {
		if len([]Image(*s)) == 1 {
			*s = ImageSet([]Image{})
		} else {
			*s = []Image(*s)[1:]
		}
		return
	}

	if index == len([]Image(*s))-1 {
		*s = []Image(*s)[:index]
		return
	} else {
		*s = append([]Image(*s)[0:index], []Image(*s)[index+1:len([]Image(*s))]...)
		return
	}
}

type SeaMap [][]Image

func (s SeaMap) ToMap() [][]rune {
	seaSize := len(s)
	contentLength := len(s[0][0].Content)
	fullMap := make([][]rune, seaSize*contentLength)
	for j := 0; j < seaSize; j++ {
		for i := 0; i < seaSize; i++ {
			for index, line := range s[j][i].Content {
				fullMap[j*contentLength+index] = append(fullMap[j*contentLength+index], line...)
			}
		}
	}
	return fullMap
}

func (sm SeaMap) ToString() string {
	lines := []string{}
	for _, tiles := range sm {
		tileIds := make([]string, len(tiles))
		for i, image := range tiles {
			tileIds[i] = image.Id
		}
		lines = append(lines, (strings.Join(tileIds, ", ")))
	}
	return strings.Join(lines, "\n")
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func (im Image) ToString() string {
	lines := make([]string, im.Size+3)
	lines[0] = fmt.Sprintf("%10v(%3v / %3v)", "", im.Sides[0], reversedSideSignature(im.Sides[0], im.Size))
	lines[1] = fmt.Sprintf("%10vTile: %v", "", im.Id)

	for i, line := range im.Content {
		if (im.Size-2)/2 == i {
			lines[i+2] = fmt.Sprintf("(%3v / %3v) %v (%3v / %3v)", im.Sides[3], reversedSideSignature(im.Sides[3], im.Size), string(line), im.Sides[1], reversedSideSignature(im.Sides[1], im.Size))
		} else {
			lines[i+2] = fmt.Sprintf("%12v%v", "", string(line))
		}
	}
	lines[len(lines)-1] = fmt.Sprintf("%10v(%3v / %3v)", "", im.Sides[2], reversedSideSignature(im.Sides[2], im.Size))
	return strings.Join(lines, "\n")
}

func reverse(s []rune, length int) []rune {
	a := make([]rune, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	for i := len(s); i < length; i++ {
		a = append(a, '0')
	}

	return a
}

func reversedSideSignature(sideSignature int, length int) int {
	reversedBinary := reverse([]rune(strconv.FormatInt(int64(sideSignature), 2)), length)
	return SideToInt(string(reversedBinary))
}

func findCorners(images []Image, sidesMap map[int][]string) (corners []Image) {
	indexMap := make(map[string]int)
	for i, image := range images {
		indexMap[image.Id] = i
	}
	lonelySides := make(map[string]int)
	for k, vs := range sidesMap {
		reversed, _ := sidesMap[reversedSideSignature(k, images[0].Size)]
		if len(vs) == 1 && len(reversed) == 1 {
			lonelySides[vs[0]]++
		}
	}
	for k, lonelySidesCount := range lonelySides {
		if lonelySidesCount == 4 {
			corners = append(corners, images[indexMap[k]])
		}
	}
	return
}

func getLastDiscoveredTileAndMatchDirection(seaMap [][]Image, tileMap map[string]Image) (Image, string, int, int) {
	var lastSetTile Image
	var matchDirection string
	var seaSize = len(seaMap)

	for j := 0; j < seaSize; j++ {
		for index := 0; index < seaSize; index++ {
			i := index
			if j%2 == 1 {
				i = seaSize - 1 - index
			}
			if seaMap[j][i].Id != "" {
				lastSetTile = seaMap[j][i]
			} else {
				// We reached the tile we want to resolve
				if i == 0 && j%2 == 0 || i == seaSize-1 && j%2 == 1 {
					matchDirection = "bottom"
				} else if j%2 == 0 {
					matchDirection = "right"
				} else {
					matchDirection = "left"
				}
				return lastSetTile, matchDirection, i, j
			}
		}
	}
	panic("Something went wrong trying to find the next tile to  fill...")
}

func cloneSeaMap(seaMap [][]Image) [][]Image {
	cloned := make([][]Image, len(seaMap))
	for j, lines := range seaMap {
		cloned[j] = make([]Image, len(lines))
		for i, col := range lines {
			cloned[j][i] = col
		}
	}
	return cloned
}

type SeaMapStruct struct {
	tileSet ImageSet
	seaMap  [][]Image
}

func RecursiveMapResolver(tileSet ImageSet, seaMap [][]Image, tileMap map[string]Image) (bool, [][]Image) {
	if len(tileSet) == 0 {
		return true, seaMap
	}
	// seaSize := len(seaMap)
	currentTile, matchDirection, i, j := getLastDiscoveredTileAndMatchDirection(seaMap, tileMap)
	// fmt.Printf("We will try to find a tile for coord (%v, %v) trying to match on the %v \n", i, j, matchDirection)
	seaMapPotentials := []SeaMapStruct{}

	for index, candidate := range tileSet {
		var matcher func(Image) (bool, []Image)
		switch matchDirection {
		case "right":
			matcher = currentTile.MatchRight
		case "bottom":
			matcher = currentTile.MatchBottom
		case "left":
			matcher = currentTile.MatchLeft
		}
		_, orientedCandidates := matcher(candidate)

		for _, orientedCandidate := range orientedCandidates {
			clonedSeaMap := cloneSeaMap(seaMap)
			clonedSeaMap[j][i] = orientedCandidate
			clonedTileSet := ImageSet(append([]Image{}, tileSet...))
			clonedTileSet.RemoveAtIndex(index)
			seaMapPotentials = append(seaMapPotentials, SeaMapStruct{clonedTileSet, clonedSeaMap})
		}
	}
	if len(seaMapPotentials) == 0 {
		fmt.Println("Did not find any")
		return false, seaMap
	}

	for _, seaMapPotential := range seaMapPotentials {
		res, seaMapResolved := RecursiveMapResolver(seaMapPotential.tileSet, seaMapPotential.seaMap, tileMap)
		if res {
			return res, seaMapResolved
		}
	}
	return false, seaMap
}

func (s ImageSet) keys() []string {
	ids := make([]string, len(s))
	for i, image := range s {
		ids[i] = image.Id
	}
	return ids
}

func (candidate *Image) Shuffle(rotationCount *int, verticallyFlipped *bool, horizontallyFlipped *bool) error {
	if *rotationCount < 3 {
		candidate.Rotate()
		*rotationCount++
	} else if !*verticallyFlipped {
		*verticallyFlipped = true
		candidate.Rotate()
		candidate.VerticalFlip()
		*rotationCount = 0
	} else if !*horizontallyFlipped {
		*horizontallyFlipped = true
		candidate.Rotate()
		candidate.VerticalFlip()
		*verticallyFlipped = false
		candidate.HorizontalFlip()
		*rotationCount = 0
	} else {
		return errors.New("Could not shuffle enough :/")
	}

	return nil
}

func (candidate *ImageLite) Shuffle(rotationCount *int, verticallyFlipped *bool, horizontallyFlipped *bool) error {
	if *rotationCount < 3 {
		candidate.Rotate()
		*rotationCount++
	} else if !*verticallyFlipped {
		*verticallyFlipped = true
		candidate.Rotate()
		candidate.VerticalFlip()
		*rotationCount = 0
	} else if !*horizontallyFlipped {
		*horizontallyFlipped = true
		candidate.Rotate()
		candidate.VerticalFlip()
		*verticallyFlipped = false
		candidate.HorizontalFlip()
		*rotationCount = 0
	} else {
		return errors.New("Could not shuffle enough :/")
	}

	return nil
}

func (img Image) MatchRight(candidate Image) (bool, []Image) {
	imagesOrientations := []Image{}
	rightSignaturePlug := reversedSideSignature(img.Sides[1], img.Size)

	rotationCount := 0
	verticallyFlipped := false
	horizontallyFlipped := false
	var err error

	for err == nil {
		if rightSignaturePlug == candidate.Sides[3] {
			alreadyConsidered := false
			for _, imageOrientation := range imagesOrientations {
				if Equal(imageOrientation.Sides, candidate.Sides) {
					alreadyConsidered = true
				}
			}
			if !alreadyConsidered {
				cloneImg := Image{candidate.Id, candidate.Content, candidate.Size, candidate.Sides}
				imagesOrientations = append(imagesOrientations, cloneImg)
			}
		}
		err = candidate.Shuffle(&rotationCount, &verticallyFlipped, &horizontallyFlipped)
	}
	return len(imagesOrientations) != 0, imagesOrientations
}

func (img Image) MatchLeft(candidate Image) (bool, []Image) {
	imagesOrientations := []Image{}
	leftSignaturePlug := reversedSideSignature(img.Sides[3], img.Size)
	rotationCount := 0
	verticallyFlipped := false
	horizontallyFlipped := false

	var err error
	for err == nil {
		if leftSignaturePlug == candidate.Sides[1] {
			alreadyConsidered := false
			for _, imageOrientation := range imagesOrientations {
				if Equal(imageOrientation.Sides, candidate.Sides) {
					alreadyConsidered = true
				}
			}
			if !alreadyConsidered {
				cloneImg := Image{candidate.Id, candidate.Content, candidate.Size, candidate.Sides}
				imagesOrientations = append(imagesOrientations, cloneImg)
			}
		}
		err = candidate.Shuffle(&rotationCount, &verticallyFlipped, &horizontallyFlipped)
	}
	return len(imagesOrientations) != 0, imagesOrientations
}

func (img Image) MatchBottom(candidate Image) (bool, []Image) {
	imagesOrientations := []Image{}
	bottomSignaturePlug := reversedSideSignature(img.Sides[2], img.Size)

	rotationCount := 0
	verticallyFlipped := false
	horizontallyFlipped := false
	var err error

	for err == nil {
		if bottomSignaturePlug == candidate.Sides[0] {
			alreadyConsidered := false
			for _, imageOrientation := range imagesOrientations {
				if Equal(imageOrientation.Sides, candidate.Sides) {
					alreadyConsidered = true
				}
			}
			if !alreadyConsidered {
				cloneImg := Image{candidate.Id, candidate.Content, candidate.Size, candidate.Sides}
				imagesOrientations = append(imagesOrientations, cloneImg)
			}
		}
		err = candidate.Shuffle(&rotationCount, &verticallyFlipped, &horizontallyFlipped)
	}
	return len(imagesOrientations) != 0, imagesOrientations
}

func (img Image) IsTopLeftCorner(sidesMap map[int][]string) bool {
	topSide := img.Sides[0]
	leftSide := img.Sides[3]

	potentialTopMatchs := len(sidesMap[topSide]) + len(sidesMap[reversedSideSignature(topSide, img.Size)])
	potentialLeftMatchs := len(sidesMap[leftSide]) + len(sidesMap[reversedSideSignature(leftSide, img.Size)])
	return potentialLeftMatchs == 2 && potentialTopMatchs == 2
}

func findIndex(slice []Image, item Image) int {
	for i, sliceItem := range slice {
		if sliceItem.Id == item.Id {
			return i
		}
	}
	return -1
}

func findIndexById(slice []Image, id string) int {
	for i, sliceItem := range slice {
		if sliceItem.Id == id {
			return i
		}
	}
	return -1
}

func (img *Image) VerticalFlip() {
	for j := 0; j < img.Size-2; j++ {
		img.Content[j] = reverse(img.Content[j], img.Size-2)
	}
	img.Sides[0], img.Sides[2] = reversedSideSignature(img.Sides[0], img.Size), reversedSideSignature(img.Sides[2], img.Size)
	img.Sides[1], img.Sides[3] = reversedSideSignature(img.Sides[3], img.Size), reversedSideSignature(img.Sides[1], img.Size)
}

func (img *Image) HorizontalFlip() {
	for j := 0; j < (img.Size-2)/2; j++ {
		img.Content[j], img.Content[img.Size-2-j-1] = img.Content[img.Size-2-j-1], img.Content[j]
	}
	img.Sides[0], img.Sides[2] = reversedSideSignature(img.Sides[2], img.Size), reversedSideSignature(img.Sides[0], img.Size)
	img.Sides[1], img.Sides[3] = reversedSideSignature(img.Sides[1], img.Size), reversedSideSignature(img.Sides[3], img.Size)
}

type Shufflable interface {
	Rotate()
	VerticalFlip()
	HorizontalFlip()
}
type ImageLite [][]rune

func (img *ImageLite) Rotate() {
	size := len(*img)
	newContent := make([][]rune, size)

	for j := 0; j < size; j++ {
		newContent[j] = make([]rune, size)
		for i := 0; i < size; i++ {
			newContent[j][i] = (*img)[size-i-1][j]
		}
	}
	imgLite := ImageLite(newContent)
	copy(*img, imgLite)
}
func (img *ImageLite) VerticalFlip() {
	size := len(*img)
	for j := 0; j < size; j++ {
		(*img)[j] = reverse((*img)[j], size)
	}
}
func (img *ImageLite) HorizontalFlip() {
	size := len(*img)
	for j := 0; j < (size)/2; j++ {
		(*img)[j], (*img)[size-j-1] = (*img)[size-j-1], (*img)[j]
	}
}

func (img *Image) Rotate() {
	newContent := make([][]rune, img.Size-2)

	for j := 0; j < img.Size-2; j++ {
		newContent[j] = make([]rune, img.Size-2)
		for i := 0; i < img.Size-2; i++ {
			newContent[j][i] = img.Content[img.Size-2-i-1][j]
		}
	}
	img.Sides = append(img.Sides[3:4], img.Sides[0:3]...)
	img.Content = newContent
}

func Equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
