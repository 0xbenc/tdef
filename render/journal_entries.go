package render

import (
	"fmt"
	"strings"

	"github.com/0xbenc/tdef/game"
)

var EnemyJournalEntries = []TowerJournalEntry{
	{"squire", -1, "Squire", "the first expedition", "The helm slips when he runs. He lifts it with one hand and remembers too late that the sword needs both. Someone taught him that dragons hoard gold. No one told him what else might be living beside it.", "Fast feet, ordinary armor. Do not let sympathy open a hole in the line.", 223, drawSquirePortrait},
	{"rogue", -1, "Rogue", "a blade ahead of its shadow", "I found the scarf caught on a thorn before I heard the footsteps. They send the quick ones where they think we are looking elsewhere. The rogue never looks toward the dragon. Only toward the next place to hide.", "Cover the turns. Speed is dangerous where your towers cannot keep sight.", 109, drawRoguePortrait},
	{"mercenary", -1, "Mercenary", "paid to be here", "He rests his hands on the cleaver like a man waiting for a cart. His armor has belonged to several people, and none of them had an easy life. I wonder what the guild promised him. I wonder whether they intend to pay.", "He takes more punishment than a recruit. Give him more than a passing shot.", 180, drawMercenaryPortrait},
	{"paladin", -1, "Paladin", "certain of the wrong thing", "The shield shines even underground. He has polished the guild mark until he can see himself in it. When he calls us monsters, he sounds relieved. It must be easier to raise a mace when the answer has already been given to you.", "A great deal of steel to break. Concentrate your strongest blows before he reaches the heart.", 223, drawPaladinPortrait},
	{"necromancer", -1, "Necromancer", "the guild spends its dead twice", "We buried the first two squires outside the garden. The next evening they came back wearing the same mud. The necromancer walked behind them, counting. I thought the guild had run out of mercy. It turns out they have found a use for that too.", "His death leaves two more fighters. Keep something ready to catch them.", 147, drawNecromancerPortrait},
	{"player", -1, "The Player", "the story they came to finish", "Every banner turns toward this one. Every frightened recruit stands a little straighter. I know that kind of certainty: the belief that the world has arranged itself around your arrival. I hope, when the gate stays shut, they learn that someone was already here.", "Save strength for the champion. A breach costs six lives, not one.", 180, drawPlayerPortrait},
	{"wizard", -1, "Wizard", "the shape of a spell", "His book has more bookmarks than pages. He mutters as though the stone is an inattentive pupil, then bends the air into a shape it was never meant to hold. Even the other guild fighters give him room. I do the same, for different reasons.", "Frail, but quick. Catch him before he slips through an uncovered stretch.", 147, drawWizardPortrait},
	{"centurion", -1, "Centurion", "hold the line", "The crest shows above the others like a strip of sunset. He does not hurry. His shield arrives with a blue light that turns aside blows, and the recruits gather behind it. The guild calls this protection. They have brought it here to take ours away.", "His ward shrugs off part of every hit. Sustained, concentrated fire wears him down.", 153, drawCenturionPortrait},
}

func EnemyJournalID(k game.EnemyKind) string {
	if k < 0 || int(k) >= len(EnemyJournalEntries) {
		return ""
	}
	return EnemyJournalEntries[k].ID
}

var JournalFloorIDs = []string{"rotunda", "rift", "halls", "garden", "heart", "depths"}
var JournalFloorNames = []string{"Rotunda", "Rift", "Long Halls", "Sunken Garden", "Heart", "Unmapped Depths"}

var PlaceJournalEntries = []TowerJournalEntry{
	{"rotunda", -1, "The Rotunda", "a roof around a promise", "The builders meant this chamber for an audience. Malgrath made it a place to sleep. The gold is softer than the stone, he told me once, though I think he only liked the sound it made when he settled. I still sweep a path to the water.", "Inner courts cover several passes with one defense. Holding the Rotunda opens the way to the Rift.", 180, drawJournalRotunda},
	{"rift", -1, "The Rift", "a wound that learned to glow", "There was a bridge here before either of us arrived. The stone remembers its ends; the middle remembers nothing. Heat comes through the crack at night. Malgrath used to warm his wings above it. Now I keep watch from the cooler bank.", "Separate banks need separate defenses. Holding the Rift opens the Long Halls.", 173, drawJournalRift},
	{"halls", -1, "The Long Halls", "footsteps before faces", "Every doorway is smaller than the one before it, until the far end seems built for someone else entirely. Malgrath disliked the echoes. I liked knowing who was coming. These days I can count a guild formation without seeing a single boot.", "The long road rewards reach and repeated coverage. Holding the Halls opens the Garden.", 180, drawJournalHalls},
	{"garden", -1, "The Sunken Garden", "something growing in the dark", "The roots found the water before we did. I thought the pale leaves meant the tree was dying. Malgrath said I was judging it by the wrong sky. There are flowers here that have never seen daylight. They open anyway.", "Slow the front ranks and catch the crowd at the bends. Holding every fixed floor unseals the Heart.", 151, drawJournalGarden},
	{"heart", -1, "The Heart", "what the walls are for", "There is no gold in the deepest chamber. Only warmth, and a hollow worn into the stone by a body that once filled the sky. I used to call every room in this place a defense. Here I remember what I am defending.", "Thin the outer ranks; finish survivors near the heart. Win all four fixed floors on this renown to enter. Holding the Heart opens the Depths.", 215, drawJournalHeart},
	{"depths", -1, "The Unmapped Depths", "the stone has another answer", "My first map was wrong by morning. My second was wrong before I finished it. The lair continues beneath the places we know, folding its roads around something I have not found. I take a torch, and leave the old map where Malgrath can see it.", "Routes change with the seed. The Depths open after the four floors and Heart are held on this renown.", 117, drawJournalDepths},
}

type journalMoment struct{ title, lore string }

var mainJournalMoments = [][3]journalMoment{
	{
		{"A Bowl by the Hoard", "After the last footsteps faded, I carried water through the Rotunda. Malgrath drank without lifting his head. Nothing in the guild's songs mentions the bowl. I think it may be the most important thing in this room."},
		{"The Place We Kept", "The gunner picked spent cases out of the gold while I reset the door. Malgrath asked whether anyone had broken the little garden pot. I told him no. Only then did he close his eyes."},
		{"Still Here", "The twentieth banner fell where the first had stood. For a while none of us moved. Then Malgrath shifted one wing, and the whole room breathed again. We had held a place. We had kept a life inside it."},
	},
	{
		{"Across the Crack", "I cut the guild banner free from the old bridge. It vanished into the glow before the cloth had finished turning. The stone stayed warm under my boots long after the fighting stopped."},
		{"A Wing Remembered", "The Rift lit the marks Malgrath's claws had left on the bank. They were higher than my head. I put my hand in one and tried to imagine the weight of him standing there, impatient to fly."},
		{"The Cooler Bank", "We held the narrow ledge until the last armored fighter fell. Across the crack, the far road was empty. I sat where Malgrath used to land and let the heat dry the shaking out of my hands."},
	},
	{
		{"The Last Echo", "The halls repeated our victory long after we stopped calling to each other. Somewhere near the far door, an echo sounded like a laugh. I let it. It had been too long since the stone had heard one."},
		{"Borrowed Helm", "A squire had left a helm beside the lowest stair. I set it upright instead of kicking it away. It was much too large for whoever had worn it. There are things a guild should not ask a child to carry."},
		{"No More Footsteps", "We waited for another formation. The echoes gave us our own breathing, a settling beam, the scrape of a sling stone put back in a pocket. At last I turned toward the garden. No one followed us from the far end."},
	},
	{
		{"Pale Flowers", "A flower opened beside a scorch mark while we were clearing the path. The frost mage stepped around it. I did too. We left the mark where it was and carried the broken stone elsewhere."},
		{"The Wrong Sky", "Malgrath had once told me the tree did not need our sky. After the defense, I found a new leaf folded beneath a root. There was no sun to thank. There was water, and room, and another night without the guild."},
		{"Room to Grow", "We held every entrance, and the garden kept its water. When the quiet came, the gunner sat under the pale branches. None of us spoke. The flowers opened around our boots as though we belonged there."},
	},
	{
		{"An Ember Kept", "Malgrath's breathing was shallow when I reached him. I counted until it steadied, then forgot the count. Behind us, the last gate was still standing. For tonight, that was enough."},
		{"His Name, Quietly", "The guild champion had called him a beast. When the chamber was ours again, I said his name. Malgrath opened one eye. There are words that take something away, and words that give it back."},
		{"The Heart Held", "The last blow struck the gate and did not open it. I stayed beside Malgrath until the heat returned to the stone. He asked whether it was morning. I said it could be, if he wanted. We had earned that much."},
	},
	{
		{"A Map Left Open", "I came back with a map that no longer agreed with the road. Malgrath studied it as if the mistake were interesting. Perhaps it was. We had found another way home."},
		{"The Road That Moved", "The path folded behind us, but every ally made it through the last bend. I marked the exit with a stone from the Rotunda. By morning it was still there. For once, the Depths had kept something I gave them."},
		{"Home by Another Way", "No banner reached the far end of the shifting road. I brought the empty map back through every chamber we had held. Malgrath listened while I described them. The lair was larger than the siege now."},
	},
}

var StoryJournalEntries = buildStoryJournalEntries()

func buildStoryJournalEntries() []TowerJournalEntry {
	var entries []TowerJournalEntry
	for floor, moments := range mainJournalMoments {
		for diff, m := range moments {
			fi, di := floor, diff
			entries = append(entries, TowerJournalEntry{fmt.Sprintf("%s:%d", JournalFloorIDs[floor], diff), -1, m.title, JournalFloorNames[floor] + " · " + DiffName(diff), m.lore, "Earned by holding " + JournalFloorNames[floor] + " on " + DiffName(diff) + ".", PlaceJournalEntries[floor].Accent, func(f *Frame, x, y, w, h int) { drawJournalMoment(f, x, y, w, h, fi, di) }})
		}
	}
	bonus := []journalMoment{
		{"A Good Morning", "Malgrath lifted his head before I brought the water. He asked whether there was anything left of the garden fruit. I had been saving it without knowing what for. We shared the smallest one."},
		{"The Spare Chair", "The slingers stayed after their stones were counted. Someone brought another chair to the fire. No one asked whose it was. By the next evening, it had become a place someone expected to find."},
		{"A Different Banner", "The guild sent a new banner, with the old tear carefully sewn shut. We found it outside the held gate. I took the cloth to patch a draft beside Malgrath's bed. It is useful now."},
		{"Another Ordinary Night", "I checked the doors, filled the bowl, and set tomorrow's stones beside the road. Malgrath was asleep before I sat down. Once I thought victory would sound like thunder. Tonight it sounds like someone resting."},
	}
	for i, m := range bonus {
		bi := i
		n := []int{1, 3, 5, 10}[i]
		entries = append(entries, TowerJournalEntry{fmt.Sprintf("afterward:%d", n), -1, m.title, "Afterward", m.lore, fmt.Sprintf("Earned after %d victories in the Heart or Unmapped Depths. A bonus memory beyond the main story.", n), 180, func(f *Frame, x, y, w, h int) { drawJournalAfterward(f, x, y, w, h, bi) }})
	}
	return entries
}

func JournalRequirement(section, cursor int) string {
	switch section {
	case JournalEnemies:
		return "Encounter this guild fighter to write its page."
	case JournalPlaces:
		return "Visit this part of the lair to write its page."
	case JournalStories:
		e := StoryJournalEntries[cursor]
		if strings.HasPrefix(e.ID, "afterward:") {
			return strings.TrimPrefix(e.ID, "afterward:") + " Heart / Depths wins · Afterward"
		}
		return "Win " + JournalFloorNames[cursor/3] + " · " + DiffName(cursor%3)
	default:
		return "Place a tower to write its page."
	}
}
