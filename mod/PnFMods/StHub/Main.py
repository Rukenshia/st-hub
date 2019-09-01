API_VERSION = 'API_v1.0'
MOD_NAME = 'StHub'

# devmenu.enable()


class StHub:
    def __init__(self):
        self.in_battle = False
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
        if name == "sfm.showResultScreen":
            with open('api/results.%s' % (self.battle_timestamp), 'w') as f:
                # Save the result screen information

                results = dock.getBattleResultInfo()

                data = {
                    'Timestamp': self.battle_timestamp,
                    'TeamID': results['team_id'],
                    'WinnerTeamID': results['winner_team_id'],
                    'BattleType': results['battle_type'],
                    'Duration': results['duration_sec'],
                    'PlaceInTeam': results['player_rank_exp'],

                    'Damage': {
                        'Sum': results['damage_sum'],
                        'Fire': results['damage_fire'],
                        'Flooding': results['damage_flood'],
                        'Ramming': results['damage_ram'],
                    },

                    'Ammo': {
                        'Torpedo': {'Damage': results['damage_tpd'], 'Shots': results['shots_tpd'], 'Hits': results['hits_tpd']},
                        'PlaneBomb': {'Damage': results['damage_bomb'], 'Shots': results['shots_bomb'], 'Hits': results['hits_bomb']},
                        'PlaneRocket': {'Damage': results['damage_rocket'], 'Shots': results['shots_rocket'], 'Hits': results['hits_rocket']},
                        'MainBatteryAP': {'Damage': results['damage_main_ap'], 'Shots': results['shots_main_ap'], 'Hits': results['hits_main_ap']},
                        'MainBatterySAP': {'Damage': results['damage_main_cs'], 'Shots': results['shots_main_cs'], 'Hits': results['hits_main_cs']},
                        'MainBatteryHE': {'Damage': results['damage_main_he'], 'Shots': results['shots_main_he'], 'Hits': results['hits_main_he']},
                        'SecondaryAP': {'Damage': results['damage_atba_ap'], 'Shots': results['shots_atba_ap'], 'Hits': results['hits_atba_ap']},
                        'SecondarySAP': {'Damage': results['damage_atba_cs'], 'Shots': results['shots_atba_cs'], 'Hits': results['hits_atba_cs']},
                        'SecondaryHE': {'Damage': results['damage_atba_he'], 'Shots': results['shots_atba_he'], 'Hits': results['hits_atba_he']},
                    },

                    'FloodsCaused': results['hits_flood'],
                    'ShipsDetected': results['detected'],
                    'LifeTime': results['life_time_sec'],
                    'PlanesKilled': results['killed_plane'],
                    'DistanceCovered': results['distance'],

                    'Economics': {
                        'Credits': results['credits'],
                        'BaseExp': results['exp'],
                    },
                }

                f.write(utils.jsonEncode(data))
        elif name == 'action.onEnterDivision':
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
        self.battle_timestamp = utils.timeNowUTC().strftime("%Y%m%d%H%M%S")

        selfInfo = battle.getSelfPlayerInfo()
        data = {
            'Status': 'active',
            'Timestamp': self.battle_timestamp,
            'ShipID': selfInfo.shipInfo.id,
            'InDivision': self.in_division,
        }

        with open('api/battle.%s' % (self.battle_timestamp), 'w') as f:
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
            'Timestamp': self.start_data['Timestamp'],
            'ShipID': self.start_data['ShipID'],
            'InDivision': self.start_data['InDivision'],
            'Survived': True if self.alive == 1 else False,
            'Kills': self.kills,
            'Damage': self.battle_damage,
        }

        if self.start_data['InDivision'] == False and self.in_division:
            data['InDivision'] = self.in_division

        with open('api/battle.%s' % (self.battle_timestamp), 'w') as f:
            f.write(utils.jsonEncode(data))

    def battle_end(self, winLoss, _unknown):
        self.sent_battle_end = True

        data = {
            'Status': 'finished',
            'Timestamp': self.start_data['Timestamp'],
            'ShipID': self.start_data['ShipID'],
            'InDivision': self.start_data['InDivision'],
            'Survived': True if self.alive == 1 else False,
            'Win': True if winLoss == 1 else False,
            'Kills': self.kills,
            'Damage': self.battle_damage,
        }

        if self.start_data['InDivision'] == False and self.in_division:
            data['InDivision'] = self.in_division

        with open('api/battle.%s' % (self.battle_timestamp), 'w') as f:
            f.write(utils.jsonEncode(data))


StHub()
