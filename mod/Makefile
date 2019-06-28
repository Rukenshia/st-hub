.PHONY: build
build: # sthub/bin/StHub.swf
	# 2. Create directory structure
	mkdir -p /mnt/g/World_Of_Warships_EU/res_mods/0.8.5.0/gui/flash
	mkdir -p /mnt/g/World_Of_Warships_EU/res_mods/0.8.5.0/gui/unbound/flash
	mkdir -p /mnt/g/World_Of_Warships_EU/res_mods/0.8.5.0/PnFMods/StHub
	touch /mnt/g/World_Of_Warships_EU/res_mods/0.8.5.0/PnFModsLoader.py


	# 3. Move required assets
	cp src/Main.py /mnt/g/World_Of_Warships_EU/res_mods/0.8.5.0/PnFMods/StHub
	cp src/USSExpressionsLoader.xml /mnt/g/World_Of_Warships_EU/res_mods/0.8.5.0/gui/flash
	cp src/dock.xml /mnt/g/World_Of_Warships_EU/res_mods/0.8.5.0/gui/unbound
	cp StHub.swf /mnt/g/World_Of_Warships_EU/res_mods/0.8.5.0/gui/unbound/flash/StHub.swf

StHub.swf: StHub.as dock.xml
	./tools/dist/xml2as.exe -i src/dock.xml -n StHub -f "S:/flex_sdk"
