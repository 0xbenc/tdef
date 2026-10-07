package game

const AllTowersMask uint32 = (1 << TowerCount) - 1
const MainTowersMask uint32 = (1 << TowerRuneforge) - 1

type Recruitment struct {
	Kind      TowerKind
	AfterWave int
}

// A recruit arrives during the preceding break, before its first useful fight.
func RecruitmentOrder(stage int) []Recruitment {
	if stage == 1 {
		return []Recruitment{{TowerFrost, 1}, {TowerCannon, 2}, {TowerSniper, 4}, {TowerFlak, 5}, {TowerTesla, 7}, {TowerMortar, 9}}
	}
	if stage == 2 {
		return []Recruitment{{TowerRuneforge, 0}, {TowerHookmaster, 2}, {TowerSappers, 4}, {TowerWitch, 6}}
	}
	return nil
}

func (s *State) TowerAvailable(k TowerKind) bool {
	return k.Valid() && (s.TrainingStage == 0 || s.UnlockedTowers&(1<<k) != 0)
}

func (s *State) ConfigureTraining(stage int, unlocked uint32) {
	s.TrainingStage = stage
	s.UnlockedTowers = unlocked | 1<<TowerGunner
	if stage == 2 {
		s.UnlockedTowers |= MainTowersMask
	}
	if stage == 1 && s.UnlockedTowers == 1<<TowerGunner {
		s.Lesson = TowerGunner
		s.LessonPending = true
	}
	s.recruitDue()
}

func (s *State) recruitDue() {
	if s.TrainingStage == 0 || s.WaveActive || s.LessonPending || s.Status != StatusRunning {
		return
	}
	for _, r := range RecruitmentOrder(s.TrainingStage) {
		if s.Wave >= r.AfterWave && !s.TowerAvailable(r.Kind) {
			s.UnlockedTowers |= 1 << r.Kind
			s.Gold += TowerSpecs[r.Kind].Cost[0]
			s.Lesson, s.LessonPending = r.Kind, true
			return
		}
	}
}
