API_VERSION = 'API_v1.0'
MOD_NAME = 'StHub'

devmenu.enable()

class StHub:
    def __init__(self):
        self.battle_timestamp = None
        self.battle_damage = 0
        self.in_division = False
        self.kills = 0
        self.alive = False
        self.sent_battle_end = False
        self.start_data = {}

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
        self.battle_timestamp = utils.timeNowUTC().strftime("%s")

        selfInfo = battle.getSelfPlayerInfo()
        data = {
            'Status': 'active',
            'Timestamp': self.battle_timestamp
            'ShipID': selfInfo.shipInfo.id,
            'InDivision': self.in_division,
        }

        with open('api/battle.%s'%(self.battle_timestamp), 'w') as f:
            f.write(utils.jsonEncode(data))

        self.battle_damage = 0
        self.kills = 0
        self.alive = True
        self.sent_battle_end = False
        self.start_data = data

    def battle_quit(self, _m):
        if self.sent_battle_end:
            return
           
        data = {
            'Status': 'abandoned',
            'Timestamp': self.start_data['Timestamp']
            'ShipID': self.start_data['ShipID']
            'InDivision': self.start_data['InDivision'],
            'Survived': True if self.alive == 1 else False,
            'Damage': self.battle_damage,
        }

        if self.start_data['InDivision'] == false and self.in_division:
            data['InDivision'] = self.in_division

        with open('api/battle.%s'%(self.battle_timestamp), 'w') as f:
            f.write(utils.jsonEncode(battle))

    def battle_end(self, winLoss, _unknown):
        self.sent_battle_end = True

        data = {
            'Status': 'finished',
            'Timestamp': self.start_data['Timestamp']
            'ShipID': self.start_data['ShipID']
            'InDivision': self.start_data['InDivision'],
            'Survived': True if self.alive == 1 else False,
            'Win': winLoss,
            'Damage': self.battle_damage,
        }

        if self.start_data['InDivision'] == false and self.in_division:
            data['InDivision'] = self.in_division

        with open('api/battle.%s'%(self.battle_timestamp), 'w') as f:
            f.write(utils.jsonEncode(battle))

StHub()
