package examples

import (
	"fmt"
	hiscores "joey-colon/osrs-hiscores"
)

func main() {
	h := hiscores.NewHiscores()

	zezima, err := h.GetPlayer("zezima")
	if err != nil {
		panic(err)
	}

	stats, err := zezima.GetSkill("strength")
	if err != nil {
		panic(err)
	}

	fmt.Println(stats) // &{1142579 1271864 75}

	rangedLevel, err := h.GetPlayerSkillLevel("sudo2", "ranged")
	fmt.Println(rangedLevel) // 99
}
