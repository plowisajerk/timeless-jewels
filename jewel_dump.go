//go:build tools

package main

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"strconv"

	"github.com/Vilsol/timeless-jewels/calculator"
	"github.com/Vilsol/timeless-jewels/data"
	"github.com/Vilsol/timeless-jewels/wasm/exposition"
)

// Uses separate steps so finder step has the new data loaded by data package

//go:generate go run tools.go re-zip
//go:generate go run tools.go find
//go:generate go run tools.go types

func main() {
	if len(os.Args) < 2 {
		return
	}

	switch os.Args[1] {
	case "re-zip":
		reZip()
	case "find":
		findAll()
	case "types":
		generateTypes()
	case "dump":
		// go run jewel_dump.go dump 2000 GloriousVanity Xibaqua
		dump_init(os.Args[2], os.Args[3], os.Args[4])
	}
}

func reZip() {
	reMarshalZip[[]*data.AlternatePassiveAddition]("alternate_passive_additions.json")
	reMarshalZip[[]*data.AlternatePassiveSkill]("alternate_passive_skills.json")
	reMarshalZip[[]*data.AlternateTreeVersion]("alternate_tree_versions.json")
	reMarshalZip[[]*data.PassiveSkill]("passive_skills.json")
	reMarshalZip[[]*data.Stat]("stats.json")
	reMarshalZip[map[string]interface{}]("SkillTree.json")
	reMarshalZip[[]interface{}]("passive_skill.min.json")
}

func reMarshalZip[T any](name string) {
	in, err := os.ReadFile("./source_data/" + name)
	if err != nil {
		panic(err)
	}

	var blob = new(T)
	if err := json.Unmarshal(in, &blob); err != nil {
		panic(err)
	}

	writeZipped("./data/"+name+".gz", blob)
}

func writeZipped(path string, data interface{}) {
	b, err := json.MarshalIndent(data,"","  ")
	if err != nil {
		panic(err)
	}

	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}

	writer := gzip.NewWriter(out)
	if _, err := writer.Write(b); err != nil {
		panic(err)
	}

	if err := writer.Close(); err != nil {
		panic(err)
	}
}

func findAll() {
	foundStats := make(map[data.JewelType]map[uint32]int)
	applicable := data.GetApplicablePassives()

	for jewelType := range data.TimelessJewelConquerors {
		foundStats[jewelType] = make(map[uint32]int)

		var firstConqueror data.Conqueror
		for conqueror := range data.TimelessJewelConquerors[jewelType] {
			firstConqueror = conqueror
			break
		}

		println(jewelType.String())

		min := data.TimelessJewelSeedRanges[jewelType].Min
		max := data.TimelessJewelSeedRanges[jewelType].Max

		if data.TimelessJewelSeedRanges[jewelType].Special {
			min /= 20
			max /= 20
		}

		for seed := min; seed <= max; seed++ {
			realSeed := seed
			if data.TimelessJewelSeedRanges[jewelType].Special {
				realSeed *= 20
			}

			if realSeed%500 == 0 {
				println(realSeed)
			}

			for _, skill := range applicable {
				if skill.IsKeystone {
					continue
				}

				results := calculator.Calculate(skill.Index, realSeed, jewelType, firstConqueror)
				if results.AlternatePassiveSkill != nil {
					for _, key := range results.AlternatePassiveSkill.StatsKeys {
						foundStats[jewelType][key]++
					}
				}

				for _, info := range results.AlternatePassiveAdditionInformations {
					if info.AlternatePassiveAddition != nil {
						for _, key := range info.AlternatePassiveAddition.StatsKeys {
							foundStats[jewelType][key]++
						}
					}
				}
			}
		}
	}

	writeZipped("./data/possible_stats.json.gz", foundStats)
}

func generateTypes() {
	e := exposition.Expose()
	tsFile, jsFile, err := e.Build()
	if err != nil {
		panic(err)
	}

	tsFile = "/* eslint-disable */\n" + tsFile
	jsFile = "/* eslint-disable */\n" + jsFile

	if err := os.MkdirAll("./frontend/src/lib/types", 0777); err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}

	if err := os.WriteFile("./frontend/src/lib/types/index.js", []byte(jsFile), 0777); err != nil {
		panic(err)
	}

	if err := os.WriteFile("./frontend/src/lib/types/index.d.ts", []byte(tsFile), 0777); err != nil {
		panic(err)
	}
}

func dump_init(seedInt string, jewelTypeStr string, conquerorStr string) {
	// var seedInt = 2000
	// var jewelTypeStr = "GloriousVanity"
	// var conquerorStr = "Xibaqua"

	//var seed = uint32(seedInt)
	seed64, err := strconv.ParseUint(seedInt, 10, 32)
  	if err != nil {
   		panic(err)
  	}
  	var seed = uint32(seed64)

	var jewelType = strToJewelType(jewelTypeStr)
	var conqueror = strToConqueror(conquerorStr)



	dump(seed, jewelType, conqueror)

	//dump(2000, data.GloriousVanity, data.Xibaqua)
	//dump(12000, data.LethalPride, data.Kaom)
	//dump(2000, data.BrutalRestraint, data.Deshret)
	//dump(4433, data.MilitantFaith, data.Maxarius)
	//dump(53740, data.ElegantHubris, data.Caspiro)


}

func dump(seed uint32, jewelType data.JewelType, timelessJewelConqueror data.Conqueror) {
	
	// testing data
	// steelwood stance
	// elegant hubris 53740
	// 80% phys 
	//var skillIndex uint32 = 2096
	//var skillIndex uint32 = 518

	//var seed uint32 = 53740
	//var seed uint32 = 2687
	//var jewelType = data.ElegantHubris
	//var timelessJewelConqueror = data.Caspiro

	var applicable = data.GetApplicablePassives()

	var alternatePassiveSkillMap = make(map[uint32]data.AlternatePassiveSkill)
	var alternatePassiveAdditionMap = make(map[uint32][]data.AlternatePassiveAdditionInformation)

	var resultMap = make(map[string]interface{})

	//var alternateSkill = false
	//var additionSkill = false

	var alternate = 0
	var addition = 0

	for _, skill := range applicable {
		if skill.IsKeystone {
			continue
		}
	
		var results = calculator.Calculate(skill.Index, seed, jewelType, timelessJewelConqueror)
		//calculator.Calculate(skillIndex, seed, jewelType, timelessJewelConqueror)


		if results.AlternatePassiveSkill != nil {
			//printAlternatePassiveSkill(results.AlternatePassiveSkill)
			alternate++
			alternatePassiveSkillMap[skill.Index] = *results.AlternatePassiveSkill
		}

		// good
		// if results.AlternatePassiveAdditionInformations != nil {
		if len(results.AlternatePassiveAdditionInformations) != 0 {
			addition++
			alternatePassiveAdditionMap[skill.Index] = results.AlternatePassiveAdditionInformations
		}
		
	}

	// println("Found " + strconv.Itoa(alternate) + " alternate passives")
	// println("Found " + strconv.Itoa(addition) + " additive passives")

	if alternate > 0 {
		// var json, err = json.Marshal(alternatePassiveSkillMap)
		// if err != nil {
		// 	panic(err)
		// }
		// print(string(json))
		// write(seed, jewelType, timelessJewelConqueror, json)
		resultMap["alternate"] = alternatePassiveSkillMap
	}
	if addition > 0 {
		// var json, err = json.Marshal(alternatePassiveAdditionMap)
		// if err != nil {
		// 	panic(err)
		// }
		// print(string(json))	
		// write(seed, jewelType, timelessJewelConqueror, json)
		resultMap["additive"] = alternatePassiveAdditionMap
	}

	// var json, err = json.Marshal(resultMap)
	// if err != nil {
	// 	panic(err)
	// }

	// write(seed, jewelType, timelessJewelConqueror, resultMap)
	write(seed, jewelType, timelessJewelConqueror, resultMap)

	//var passiveSkill1 = alternatePassiveSkillMap[2096]
	//var passiveSkill2 = alternatePassiveSkillMap[518]

	//printAlternatePassiveSkill(&passiveSkill1)
	//printAlternatePassiveSkill(&passiveSkill2)

}

func write(seed uint32, jewelType data.JewelType, timelessJewelConqueror data.Conqueror, json interface{}) {
	var location = "../plow-5way/jewels/"
	// var name = jewelTypeString(jewelType) + "_" + strconv.Itoa(int(seed)) + "_" + conquerorString(timelessJewelConqueror) + ".json.gz"
    var name = jewelTypeString(jewelType) + "_" + strconv.Itoa(int(seed)) + ".json.gz"

	//println(location + name)

	writeZipped(location + name, json)

	println("Successfully wrote " + location + name)
}

func printAlternatePassiveSkill(skill *data.AlternatePassiveSkill) {
	println(skill.Index)
	println(skill.ID)
	println(skill.AlternateTreeVersionsKey)
	println(skill.Name)
	println(skill.PassiveType)
	println(skill.StatsKeys)
	println(skill.Stat1Min)
	println(skill.Stat1Max)
	println(skill.Stat2Min)
	println(skill.Stat2Max)
	println(skill.Stat2Max)
	println(skill.Stat2Max)
	println(skill.Stat3Min)
	println(skill.Stat3Max)
	println(skill.Stat4Min)
	println(skill.Stat4Max)
	println(skill.SpawnWeight)
	println(skill.ConquerorIndex)
	println(skill.RandomMin)
	println(skill.ConquerorVersion)
}

func printAlternatePassiveAddition(skill *data.AlternatePassiveAddition) {
	println(skill.Index)
	println(skill.ID)
	println(skill.AlternateTreeVersionsKey)
	println(skill.SpawnWeight)
	println(skill.StatsKeys)
	println(skill.Stat1Min)
	println(skill.Stat1Max)
	println(skill.Stat2Min)
	println(skill.Stat2Max)
	println(skill.PassiveType)
}

func jewelTypeString(jewelType data.JewelType) string {
	switch jewelType {
	case data.GloriousVanity:  return "GloriousVanity"
	case data.LethalPride:     return "LethalPride"
	case data.BrutalRestraint: return "BrutalRestraint"
	case data.MilitantFaith:   return "MilitantFaith"
	case data.ElegantHubris:   return "ElegantHubris"
	default:                   return "NULL"
	}
}

func conquerorString(conqueror data.Conqueror) string {
	switch conqueror {
	case data.Xibaqua:  return "Xibaqua"
	case data.Zerphi:   return "Zerphi"
	case data.Ahuana:   return "Ahuana"
	case data.Doryani:  return "Doryani"

	case data.Kaom:     return "Kaom"
	case data.Rakiata:  return "Rakiata"
	case data.Kiloava:  return "Kiloava"
	case data.Akoya:    return "Akoya"

	case data.Deshret:  return "Deshret"
	case data.Balbala:  return "Balbala"
	case data.Asenath:  return "Asenath"
	case data.Nasima:   return "Nasima"

	case data.Venarius: return "Venarius"
	case data.Maxarius: return "Maxarius"
	case data.Dominus:  return "Dominus"
	case data.Avarius:  return "Avarius"

	case data.Cadiro:   return "Cadiro"
	case data.Victario: return "Victario"
	case data.Chitus:   return "Chitus"
	case data.Caspiro:  return "Caspiro"
    default: return "NULL"
	}
}

func strToJewelType(jewelTypeStr string) data.JewelType {
	var answer = data.GloriousVanity
	switch jewelTypeStr {
	case "GloriousVanity": 		answer = data.GloriousVanity
	case "LethalPride":			answer = data.LethalPride
	case "BrutalRestraint":		answer = data.BrutalRestraint
	case "MilitantFaith":		answer = data.MilitantFaith
	case "ElegantHubris": 		answer = data.ElegantHubris
	default:                   	panic("Invalid Jewel Type")
	}
	return answer
}

func strToConqueror(conquerorStr string) data.Conqueror {
	var answer = data.Xibaqua
	switch conquerorStr {
	case "Xibaqua":  answer = data.Xibaqua
	case "Zerphi":   answer = data.Zerphi
	case "Ahuana":   answer = data.Ahuana
	case "Doryani":  answer = data.Doryani

	case "Kaom":     answer = data.Kaom
	case "Rakiata":  answer = data.Rakiata
	case "Kiloava":  answer = data.Kiloava
	case "Akoya":    answer = data.Akoya

	case "Deshret":  answer = data.Deshret
	case "Balbala":  answer = data.Balbala
	case "Asenath":  answer = data.Asenath
	case "Nasima":   answer = data.Nasima

	case "Venarius": answer = data.Venarius
	case "Maxarius": answer = data.Maxarius
	case "Dominus":  answer = data.Dominus
	case "Avarius":  answer = data.Avarius

	case "Cadiro":   answer = data.Cadiro
	case "Victario": answer = data.Victario
	case "Chitus":   answer = data.Chitus
	case "Caspiro":  answer = data.Caspiro
    default: panic("Invalid Conqueror")
	}
	return answer
}

