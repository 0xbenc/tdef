Cinematic copy

The rest of the game's authored copy is editable in JSON under
internal/copytext/data/. See internal/copytext/README.txt and LORE-SOURCES.md
for the source map. The two cinematic files stay here.

Edit opening.json or ending.json in VS Code. These are the actual copies used
by the game. Each currently has 17 scenes.

Each entry matches one existing cinematic shot, in order:
  title   The small heading shown above the scene.
  speaker Leave empty for narration, or use GRAK, MALGRATH, GUILDMASTER, or HEALER.
  text    The narration or spoken dialogue. It wraps automatically on screen.

Keep the 17 entries in their current order while reviewing wording. Adding,
removing, or rearranging scenes also requires changing the artwork sequence
in render/cutscene.go.

Keep text concise: long captions reduce the space available for artwork.
Use JSON strings; escape a quotation mark inside text as \".

After editing, run go test ./... to check the copy and terminal layouts.
Rebuild with go build -o termtd.exe . to play your changes. The copy is embedded
in the executable; saving the file does not change an already running game.

Preview with .\termtd.exe intro or .\termtd.exe ending.

Checklist items 01 and 02 use opening.json; item 03 uses ending.json.
Empty speakers are narration; named speakers are dialogue. Artwork and scene
order remain in render/cutscene.go; camera timing lives in cutscene_panorama.go.
