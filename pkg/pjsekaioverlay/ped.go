package pjsekaioverlay

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/TootieJin/pjsekai-overlay-MULTI/pkg/sonolus"
)

type PedFrame struct {
	Time       float64
	TotalScore float64
	Score      float64
	Combo      []int
}

type BpmChange struct {
	Beat float64
	Bpm  float64
}

var WEIGHT_MAP = map[string]float64{
	"#BPM_CHANGE":    0,
	"Initialization": 0,
	"InputManager":   0,
	"Stage":          0,

	"NormalTapNote":   1,
	"CriticalTapNote": 2,

	"NormalFlickNote":   1,
	"CriticalFlickNote": 3,

	"NormalSlideStartNote":   1,
	"CriticalSlideStartNote": 2,

	"NormalSlideEndNote":   1,
	"CriticalSlideEndNote": 2,

	"NormalSlideEndFlickNote":   1,
	"CriticalSlideEndFlickNote": 3,

	"HiddenSlideTickNote":   0,
	"NormalSlideTickNote":   0.1,
	"CriticalSlideTickNote": 0.2,

	"IgnoredSlideTickNote":          0.1,
	"NormalAttachedSlideTickNote":   0.1,
	"CriticalAttachedSlideTickNote": 0.2,

	"NormalSlideConnector":   0,
	"CriticalSlideConnector": 0,

	"SimLine": 0,

	"NormalSlotEffect":       0,
	"SlideSlotEffect":        0,
	"FlickSlotEffect":        0,
	"CriticalSlotEffect":     0,
	"NormalSlotGlowEffect":   0,
	"SlideSlotGlowEffect":    0,
	"FlickSlotGlowEffect":    0,
	"CriticalSlotGlowEffect": 0,

	"NormalTraceNote":   0.1,
	"CriticalTraceNote": 0.2,

	"NormalTraceSlotEffect":     0,
	"NormalTraceSlotGlowEffect": 0,

	"DamageNote":           0.1,
	"DamageSlotEffect":     0,
	"DamageSlotGlowEffect": 0,

	"NormalTraceFlickNote":         1,
	"CriticalTraceFlickNote":       3,
	"NonDirectionalTraceFlickNote": 1,

	"NormalTraceSlideStartNote":   0.1,
	"NormalTraceSlideEndNote":     0.1,
	"CriticalTraceSlideStartNote": 0.2,
	"CriticalTraceSlideEndNote":   0.2,

	"TimeScaleGroup":  0,
	"TimeScaleChange": 0,
}

func getValueFromData(data []sonolus.LevelDataEntityValue, name string) (float64, error) {
	for _, value := range data {
		if value.Name == name {
			return value.Value, nil
		}
	}
	return 0, fmt.Errorf("value not found: %s", name)
}

func getTimeFromBpmChanges(bpmChanges []BpmChange, beat float64) float64 {
	ret := 0.0
	for i, bpmChange := range bpmChanges {
		if i == len(bpmChanges)-1 {
			ret += (beat - bpmChange.Beat) * (60 / bpmChange.Bpm)
			break
		}
		nextBpmChange := bpmChanges[i+1]
		if beat >= bpmChange.Beat && beat < nextBpmChange.Beat {
			ret += (beat - bpmChange.Beat) * (60 / bpmChange.Bpm)
			break
		} else if beat >= nextBpmChange.Beat {
			ret += (nextBpmChange.Beat - bpmChange.Beat) * (60 / bpmChange.Bpm)
		} else {
			break
		}
	}
	return ret
}

func CalculateScore(levelInfo sonolus.LevelInfo, levelData sonolus.LevelData, power float64, powerWeight []float64, missRate []float64, playerPos int) []PedFrame {
	rating := levelInfo.Rating
	var weightedNotesCount float64 = 0
	for _, entity := range levelData.Entities {
		weight := WEIGHT_MAP[entity.Archetype]
		if weight == 0 {
			continue
		}
		weightedNotesCount += weight
	}

	frames := make([]PedFrame, 0, int(weightedNotesCount)+1)
	frames = append(frames, PedFrame{Time: 0, Score: 0})
	bpmChanges := ([]BpmChange{})
	levelFax := float64(rating-5)*0.005 + 1
	comboFax := []float64{1.0, 1.0, 1.0, 1.0, 1.0}

	score := []float64{0.0, 0.0, 0.0, 0.0, 0.0}
	entityCounter := []int{0, 0, 0, 0, 0}
	noteEntities := ([]sonolus.LevelDataEntity{})

	for _, entity := range levelData.Entities {
		weight := WEIGHT_MAP[entity.Archetype]
		if weight > 0.0 && len(entity.Data) > 0 {
			noteEntities = append(noteEntities, entity)
		} else if entity.Archetype == "#BPM_CHANGE" {
			beat, err := getValueFromData(entity.Data, "#BEAT")
			if err != nil {
				continue
			}
			bpm, err := getValueFromData(entity.Data, "#BPM")
			if err != nil {
				continue
			}
			bpmChanges = append(bpmChanges, BpmChange{
				Beat: beat,
				Bpm:  bpm,
			})
		}
	}
	sort.SliceStable(noteEntities, func(i, j int) bool {
		return noteEntities[i].Data[0].Value < noteEntities[j].Data[0].Value
	})
	sort.SliceStable(bpmChanges, func(i, j int) bool {
		return bpmChanges[i].Beat < bpmChanges[j].Beat
	})
	for _, entity := range noteEntities {
		weight := WEIGHT_MAP[entity.Archetype]

		for i := range 5 {
			entityCounter[i] += 1
			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)))

			if i+1 != playerPos {
				if r.Float64() < missRate[i] {
					entityCounter[i] = 0
					comboFax[i] = 1.0
				}
			}

			if entityCounter[i]%100 == 1 && entityCounter[i] > 1 {
				comboFax[i] += 0.01
			}
			if comboFax[i] > 1.1 {
				comboFax[i] = 1.1
			}

			score[i] += ((float64(power*powerWeight[i]) / weightedNotesCount) * // Team power / weighted notes count
				4 * // Constant
				weight * // Note weight
				1 * // Judge weight (Always 1)
				levelFax * // Level fax
				comboFax[i] * // Combo fax
				1) // Skill fax (Always 1)
		}

		beat, err := getValueFromData(entity.Data, "#BEAT")
		if err != nil {
			continue
		}
		pos := 0
		switch playerPos {
		case 1:
			pos = 0
		case 2:
			pos = 1
		case 3:
			pos = 2
		case 4:
			pos = 3
		case 5:
			pos = 4
		}

		frames = append(frames, PedFrame{
			Time:       getTimeFromBpmChanges(bpmChanges, beat) + levelData.BgmOffset,
			TotalScore: score[0] + score[1] + score[2] + score[3] + score[4],
			Score:      score[pos],
			Combo:      []int{entityCounter[0], entityCounter[1], entityCounter[2], entityCounter[3], entityCounter[4]},
		})
	}

	return frames
}

func WritePedFile(frames []PedFrame, assets string, path string, levelInfo sonolus.LevelInfo, levelData sonolus.LevelData, powerWeight []float64, enUI bool) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました (Failed to create file.) [%s]", err)
	}
	defer file.Close()

	writer := io.Writer(file)

	fmt.Fprintf(writer, "p|%s\n", assets)
	fmt.Fprintf(writer, "e|%s\n", strconv.FormatBool(enUI))
	fmt.Fprintf(writer, "v|%s\n", Version)
	fmt.Fprintf(writer, "u|%d\n", time.Now().Unix())

	lastScore := 0.0
	lastScore2 := 0.0
	rating := levelInfo.Rating
	for i, frame := range frames {
		// 2-variable scoring (supports accurate digits up to 1e+34)
		score := math.Mod(frame.TotalScore, 1e+17)
		score2 := math.Floor(frame.TotalScore / 1e+17)
		if score2 < 0 {
			score2 = math.Ceil(frame.TotalScore / 1e+17)
		}

		if math.Ceil(score) < 0 && math.Ceil(score2) < 0 {
			score = -score
		}

		frameScore := math.Mod(frame.Score-(lastScore+(lastScore2*1e+17)), 1e+17)
		frameScore2 := math.Floor(frame.Score/1e+17) - lastScore2

		lastScore = math.Mod(frame.Score, 1e+17)
		lastScore2 = math.Floor(frame.Score / 1e+17)
		if lastScore2 < 0 {
			lastScore2 = math.Ceil(frame.Score / 1e+17)
		}

		if math.Ceil(frameScore) < 0 && math.Floor(frameScore2) > 0 {
			frameScore = 1e+16 - frameScore
		}

		if math.Ceil(frameScore) < 0 && math.Ceil(frameScore2) < 0 {
			frameScore = -frameScore
		}

		rank := "n"
		scoreX := 0.0

		if rating < 5 {
			rating = 5
		} else if rating > 40 {
			rating = 40
		}

		rankBorder := float64(1200000+(rating-5)*4100) * 5
		rankS := float64(1040000+(rating-5)*5200) * 5
		rankA := float64(840000+(rating-5)*4200) * 5
		rankB := float64(400000+(rating-5)*2000) * 5
		rankC := float64(40000+(rating-5)*200) * 5

		// bar
		if math.Ceil(score2) < 0 || math.Ceil(score) < 0 {
			rank = "d"
			scoreX = 0
		} else if score >= rankBorder || math.Floor(score2) > 0 {
			rank = "s"
			scoreX = 372
		} else if score >= rankS {
			rank = "s"
			scoreX = (float64((score-rankS))/float64((rankBorder-rankS)))*36 + 335
		} else if score >= rankA {
			rank = "a"
			scoreX = (float64((score-rankA))/float64((rankS-rankA)))*55 + 280
		} else if score >= rankB {
			rank = "b"
			scoreX = (float64((score-rankB))/float64((rankA-rankB)))*55 + 225
		} else if score >= rankC {
			rank = "c"
			scoreX = (float64((score-rankC))/float64((rankB-rankC)))*55 + 170
		} else {
			rank = "d"
			scoreX = (float64(score) / float64(rankC)) * 168
		}

		var totalWeight float64 = 0.0
		for j := range powerWeight {
			totalWeight = totalWeight + powerWeight[j]
		}

		scoreX5 := (scoreX / 372) * ((powerWeight[0] + powerWeight[1] + powerWeight[2] + powerWeight[3] + powerWeight[4]) / totalWeight)
		scoreX4 := (scoreX / 372) * ((powerWeight[0] + powerWeight[1] + powerWeight[2] + powerWeight[3]) / totalWeight)
		scoreX3 := (scoreX / 372) * ((powerWeight[0] + powerWeight[1] + powerWeight[2]) / totalWeight)
		scoreX2 := (scoreX / 372) * ((powerWeight[0] + powerWeight[1]) / totalWeight)
		scoreX1 := (scoreX / 372) * ((powerWeight[0]) / totalWeight)

		time := frame.Time
		if time == 0 && i > 0 {
			time = frames[i-1].Time + 0.000001
		}

		writer.Write(fmt.Appendf(nil, "s|%f:%.0f:%.0f:%.0f:%.0f:%f:%f:%f:%f:%f:%s:%d\n", time, score2, score, frameScore2, frameScore, scoreX1, scoreX2, scoreX3, scoreX4, scoreX5, rank, i))
	}

	return nil
}

func WritePedMultiFile(frames []PedFrame, assets string, path string, levelData sonolus.LevelData, playerPos int) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました (Failed to create file.) [%s]", err)
	}
	defer file.Close()

	writer := io.Writer(file)

	fmt.Fprintf(writer, "p|%s\n", assets)
	fmt.Fprintf(writer, "l|%d\n", playerPos)
	fmt.Fprintf(writer, "v|%s\n", Version)
	fmt.Fprintf(writer, "u|%d\n", time.Now().Unix())

	lastTime := 0.0
	lastCombo := []int{0, 0, 0, 0, 0}
	for i, frame := range frames {
		score := math.Mod(frame.TotalScore, 1e+17)

		time := math.Ceil(frame.Time)
		if time == lastTime {
			continue
		}
		if time == 0 && i > 0 {
			time = frames[i-1].Time + 0.000001
		}
		lastTime = time

		combo1 := frame.Combo[0]
		combo2 := frame.Combo[1]
		combo3 := frame.Combo[2]
		combo4 := frame.Combo[3]
		combo5 := frame.Combo[4]
		if combo1 <= lastCombo[0] {
			combo1 = 0
		}
		if combo2 <= lastCombo[1] {
			combo2 = 0
		}
		if combo3 <= lastCombo[2] {
			combo3 = 0
		}
		if combo4 <= lastCombo[3] {
			combo4 = 0
		}
		if combo5 <= lastCombo[4] {
			combo5 = 0
		}
		lastCombo[0] = frame.Combo[0]
		lastCombo[1] = frame.Combo[1]
		lastCombo[2] = frame.Combo[2]
		lastCombo[3] = frame.Combo[3]
		lastCombo[4] = frame.Combo[4]

		writer.Write(fmt.Appendf(nil, "s|%f:%.0f:%d:%d:%d:%d:%d\n", time, score, combo1, combo2, combo3, combo4, combo5))
	}

	return nil
}
