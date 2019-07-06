API_VERSION = 'API_v1.0'
MOD_NAME = 'StHub'

devmenu.enable()

class StHub:
    def __init__(self):
        self.battle_damage = 0
        self.in_division = False
        self.kills = 0
        self.alive = False
        self.sent_battle_end = False

        events.onReceiveShellInfo(self.shell_info)
        events.onBattleEnd(self.battle_end)
        events.onBattleQuit(self.battle_quit)
        events.onBattleStart(self.battle_start)
        events.onSFMEvent(self.sfm_event)

    def sfm_event(self, name, data):
        if name == 'action.onEnterDivision':
            self.in_division = True
        elif name == 'action.leaveDivision':
            self.in_division = False

    def shell_info(self, victim, shooter, ammo, mat, shoot, flags, damage, pos, yaw, hlinfo):
        if (flags & 0b1) == 0:
            self.battle_damage = self.battle_damage + damage
            if flags & 0b1000:
                self.kills = self.kills + 1
        else:
            if flags & 0b1000:
                self.alive = False

    def battle_start(self):
        selfInfo = battle.getSelfPlayerInfo()
        data = {
            'ShipID': selfInfo.shipInfo.id,
            'ShipName': selfInfo.shipInfo.name,
            'InDivision': self.in_division,
        }

        with open('api/battle.start', 'w') as f:
            f.write(utils.jsonEncode(data))

        self.battle_damage = 0
        self.kills = 0
        self.alive = True
        self.sent_battle_end = False

    def battle_quit(self, _m):
        if self.sent_battle_end:
            return

        with open('api/battle.response', 'r') as f:
            data = f.read()

            if data == 'ERR_NOT_IN_TESTING':
                return

            battle = utils.jsonDecode(data)

            battle['Status'] = 'abandoned'
            battle['Statistics']['Survived'] = True if self.alive == 1 else False
            battle['Statistics']['Damage']['Value'] = self.battle_damage
            battle['Statistics']['Kills']['Value'] = self.kills

            if battle['Statistics']['InDivision']['Value'] == False and self.in_division:
                ## Joined division mid-game (probably)
                battle['Statistics']['InDivision']['Value'] = self.in_division

            with open('api/battle.end', 'w') as wf:
                wf.write(utils.jsonEncode(battle))
                wf.close()

            f.close()

    def battle_end(self, winLoss, _unknown):
        self.sent_battle_end = True
        with open('api/battle.response', 'r') as f:
            data = f.read()

            if data == 'ERR_NOT_IN_TESTING':
                return

            battle = utils.jsonDecode(data)

            battle['Status'] = 'finished'
            battle['Statistics']['Win'] = True if winLoss == 1 else False
            battle['Statistics']['Survived'] = True if self.alive == 1 else False
            battle['Statistics']['Damage']['Value'] = self.battle_damage
            battle['Statistics']['Kills']['Value'] = self.kills

            if battle['Statistics']['InDivision']['Value'] == False and self.in_division:
                ## Joined division mid-game (probably)
                battle['Statistics']['InDivision']['Value'] = self.in_division

            with open('api/battle.end', 'w') as wf:
                wf.write(utils.jsonEncode(battle))
                wf.close()

            f.close()


def wl(l, m):
    l.write('{}\n'.format(m))
    l.flush()


StHub()
