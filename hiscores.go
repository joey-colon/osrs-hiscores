package hiscores

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const VERSION = "0.0.1"
const HISCORES_ENDPOINT = "https://secure.runescape.com/m=hiscore_oldschool/index_lite.ws?player={rsn}"
const CACHE_TTL = 3600

type IHiscores interface {
	GetPlayer(rsn string) (*Player, error)
	GetPlayerSkillLevel(rsn string, skill string) (int64, error)
	GetPlayerSkillXp(rsn string, skill string) (int64, error)
	GetPlayerSkillRank(rsn string, skill string) (int64, error)
}

type Hiscores struct {
	cache   map[string]*Player
	skills  []string
	columns []string
}

type Player struct {
	updatedAt time.Time
	rsn       string
	stats     map[string]*Stats
}

type Stats struct {
	xp    int64
	rank  int64
	level int64
}

func NewHiscores() IHiscores {
	hiscores := &Hiscores{
		cache:   make(map[string]*Player),
		skills:  []string{"overall", "attack", "defence", "strength", "hitpoints", "ranged", "prayer", "magic", "cooking", "woodcutting", "fletching", "fishing", "firemaking", "crafting", "smithing", "mining", "herblore", "agility", "thieving", "slayer", "farming", "runecraft", "hunter", "construction"},
		columns: []string{"rank", "level", "xp"},
	}

	return hiscores
}

func isValidSkill(skill string, skills []string) bool {
	for _, s := range skills {
		if skill == s {
			return true
		}
	}
	return false
}

// GetPlayer takes in RuneScape name and responds with all stats
func (h *Hiscores) GetPlayer(rsn string) (*Player, error) {
	if player, ok := h.cache[rsn]; ok {
		duration := time.Since(player.updatedAt)
		if duration.Seconds() < CACHE_TTL {
			return player, nil
		} else {
			// evict from cache
			delete(h.cache, rsn)
		}
	}

	resp, err := http.Get(strings.Replace(HISCORES_ENDPOINT, "{rsn}", rsn, 1))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	payload := string(body)
	stats := make(map[string]*Stats)
	player := &Player{updatedAt: time.Now().UTC(), rsn: rsn, stats: stats}
	rows := strings.Split(payload, "\n")

	for idx, skill := range h.skills {
		tokens := strings.Split(rows[idx], ",")
		if len(tokens) != 3 {
			return nil, fmt.Errorf("Unable to parse response row")
		}
		xp, _ := strconv.ParseInt(tokens[0], 10, 32)
		level, _ := strconv.ParseInt(tokens[1], 10, 32)
		rank, _ := strconv.ParseInt(tokens[2], 10, 32)
		stats[skill] = &Stats{xp: xp, level: level, rank: rank}
	}

	h.cache[rsn] = player
	return h.cache[rsn], nil
}

// GetPlayerSkillLevel takes in RuneScape name / skill name and returns the skill level
func (h *Hiscores) GetPlayerSkillLevel(rsn string, skill string) (int64, error) {
	if !isValidSkill(skill, h.skills) {
		return -1, fmt.Errorf("Invalid skill.\nValid skills: " + strings.Join(h.skills, ","))
	}

	if player, ok := h.cache[rsn]; ok {
		duration := time.Since(player.updatedAt)
		if duration.Seconds() < CACHE_TTL {
			return player.stats[skill].level, nil
		} else {
			// evict from cache
			delete(h.cache, rsn)
		}
	}

	_, err := h.GetPlayer(rsn)
	if err != nil {
		return -1, err
	}

	return h.GetPlayerSkillLevel(rsn, skill)
}

// GetPlayerSkillRank takes in RuneScape name / skill name and returns the skill rank
func (h *Hiscores) GetPlayerSkillRank(rsn string, skill string) (int64, error) {
	if !isValidSkill(skill, h.skills) {
		return -1, fmt.Errorf("Invalid skill.\nValid skills: " + strings.Join(h.skills, ","))
	}

	if player, ok := h.cache[rsn]; ok {
		duration := time.Since(player.updatedAt)
		if duration.Seconds() < CACHE_TTL {
			return player.stats[skill].rank, nil
		} else {
			// evict from cache
			delete(h.cache, rsn)
		}
	}

	_, err := h.GetPlayer(rsn)
	if err != nil {
		return -1, err
	}

	return h.GetPlayerSkillRank(rsn, skill)
}

// GetPlayerSkillXp takes in RuneScape name / skill name and returns the skill xp
func (h *Hiscores) GetPlayerSkillXp(rsn string, skill string) (int64, error) {
	if !isValidSkill(skill, h.skills) {
		return -1, fmt.Errorf("Invalid skill.\nValid skills: " + strings.Join(h.skills, ","))
	}

	if player, ok := h.cache[rsn]; ok {
		duration := time.Since(player.updatedAt)
		if duration.Seconds() < CACHE_TTL {
			return player.stats[skill].xp, nil
		} else {
			// evict from cache
			delete(h.cache, rsn)
		}
	}

	_, err := h.GetPlayer(rsn)
	if err != nil {
		return -1, err
	}

	return h.GetPlayerSkillXp(rsn, skill)
}

// GetSkill returns the Stats for a given player
func (p *Player) GetSkill(skill string) (*Stats, error) {
	skills := make([]string, len(p.stats))
	for s := range p.stats {
		skills = append(skills, s)
	}

	if !isValidSkill(skill, skills) {
		return nil, fmt.Errorf("Invalid skill.\nValid skills: " + strings.Join(skills, ","))
	}

	return p.stats[skill], nil
}
