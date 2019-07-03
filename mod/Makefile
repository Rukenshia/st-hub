.PHONY: build
build: # sthub/bin/StHub.swf
	# 2. Create directory structure
	mkdir -p /mnt/c/Games/World_Of_Warships/res_mods/0.8.5.1
	mkdir -p /mnt/c/Games/World_Of_Warships/res_mods/0.8.5.1/gui/flash
	mkdir -p /mnt/c/Games/World_Of_Warships/res_mods/0.8.5.1/gui/unbound/flash
	mkdir -p /mnt/c/Games/World_Of_Warships/res_mods/0.8.5.1/PnFMods/StHub
	touch /mnt/c/Games/World_Of_Warships/res_mods/0.8.5.1/PnFModsLoader.py


	# 3. Move required assets
	cp src/Main.py /mnt/c/Games/World_Of_Warships/res_mods/0.8.5.1/PnFMods/StHub
	# cp StHubBridge/bin/Main.swf /mnt/c/Games/World_Of_Warships/res_mods/0.8.5.1/PnFMods/StHub
	# cp src/USSExpressionsLoader.xml /mnt/g/World_Of_Warships_EU/res_mods/0.8.5.1/gui/flash
	# cp src/dock.xml /mnt/g/World_Of_Warships_EU/res_mods/0.8.5.1/gui/unbound
	# cp StHub.swf /mnt/g/World_Of_Warships_EU/res_mods/0.8.5.1/gui/unbound/flash/StHub.swf

# StHub.swf: StHub.as dock.xml
# 	./tools/dist/xml2as.exe -i src/dock.xml -n StHub -f "S:/flex_sdk"

SOURCE_DIR=examples

.PHONY: process_logs
process_logs:
	mkdir -p $(SOURCE_DIR)/parsed
	mkdir -p examples/parsed

	cat $(SOURCE_DIR)/battle.players.log | sed "s/'/\"/g" | sed 's/<Dummy>: //g' | sed 's/L,/,/g' | sed 's/L]/]/g' \
		| sed 's/u"/"/g' | sed -e 's/\([0-9]\+\):/"\1":/g' | sed 's/None/null/g' | sed 's/True/true/g' | sed 's/False/false/g' \
		> examples/parsed/battle.players.raw.json
	
	jq -r '.' examples/parsed/battle.players.raw.json > examples/parsed/battle.players.json

	cat $(SOURCE_DIR)/battle.players_end.log | sed "s/'/\"/g" | sed 's/<Dummy>: //g' | sed 's/L,/,/g' | sed 's/L]/]/g' \
		| sed 's/u"/"/g' | sed -e 's/\([0-9]\+\):/"\1":/g' | sed 's/None/null/g' | sed 's/True/true/g' | sed 's/False/false/g' \
		> examples/parsed/battle.players_end.raw.json
	
	jq -r '.' examples/parsed/battle.players_end.raw.json > examples/parsed/battle.players_end.json

	cat $(SOURCE_DIR)/battle.self.log | sed "s/'/\"/g" | sed 's/<Dummy>: //g' | sed 's/L,/,/g' | sed 's/L]/]/g' \
		| sed 's/u"/"/g' | sed -e 's/\([0-9]\+\):/"\1":/g' | sed 's/None/null/g' | sed 's/True/true/g' | sed 's/False/false/g' \
		| sed 's/(<.*>)/"<removed>"/g' \
		> examples/parsed/battle.self.raw.json
	
	jq -r '.' examples/parsed/battle.self.raw.json > examples/parsed/battle.self.json

	cat $(SOURCE_DIR)/battle.self_end.log | sed "s/'/\"/g" | sed 's/<Dummy>: //g' | sed 's/L,/,/g' | sed 's/L]/]/g' \
		| sed 's/u"/"/g' | sed -e 's/\([0-9]\+\):/"\1":/g' | sed 's/None/null/g' | sed 's/True/true/g' | sed 's/False/false/g' \
		| sed 's/(<.*>)/"<removed>"/g' \
		> examples/parsed/battle.self_end.raw.json
	
	jq -r '.' examples/parsed/battle.self_end.raw.json > examples/parsed/battle.self_end.json