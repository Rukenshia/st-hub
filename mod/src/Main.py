API_VERSION = 'API_v1.0'
MOD_NAME = 'StHub'

devmenu.enable()

ev_log = open('events.log', 'w')
battle_log = open('battle.log', 'w')

class StHub:
    def __init__(self):
        self.battle_damage = 0

        events.onReceiveShellInfo(self.shell_info)
        events.onBattleEnd(self.battle_end)
        events.onBattleQuit(self.battle_quit)
        events.onBattleStart(self.battle_start)

    def shell_info(self, victim, shooter, ammo, mat, shoot, flags, damage, pos, yaw, hlinfo):
        if flags & 0b1:
            flash.call('sthub.LastEvent', ['StHub.Shell received damage %s'%(damage)])
            wl(battle_log, 'shell: %s %s %s %s %s %s %s %s %s %s'%(victim, shooter, ammo, mat, shoot, '{0:b}'.format(flags), damage, pos, yaw, hlinfo))
        else:
            flash.call('sthub.LastEvent', ['StHub.Shell dealt damage %s'%(damage)])
            wl(battle_log, 'shell: %s %s %s %s %s %s %s %s %s %s'%(victim, shooter, ammo, mat, shoot, '{0:b}'.format(flags), damage, pos, yaw, hlinfo))

            self.battle_damage = self.battle_damage + damage
        
        flash.call('sthub.LastEvent', ['StHub.AmmoInfo %s'%(battle.getAmmoParams(ammo))])
        wl(battle_log, 'ammo_info %s'%(battle.getAmmoParams(ammo)))

    def battle_start(self):
        wl(battle_log, '--- battle start')
        flash.call('sthub.LastEvent', ['StHub.BattleStart'])

        wl(open('battle.self.log', 'w'), '%s'%(battle.getSelfPlayerInfo()))
        wl(open('battle.players.log', 'w'), '%s'%(battle.getPlayersInfo()))

        self.battle_damage = 0

    def battle_quit(self, _m):
        flash.call('sthub.LastEvent', ['StHub.BattleQuit'])
        wl(battle_log, '--- battle quit')

        wl(open('battle.self_end.log', 'w'), '%s'%(battle.getSelfPlayerInfo()))
        wl(open('battle.players_end.log', 'w'), '%s'%(battle.getPlayersInfo()))


    def battle_end(self, a, b):
        flash.call('sthub.LastEvent', ['StHub.BattleEnd %s %s'%(a, b)])
        flash.call('sthub.LastEvent', ['StHub.AvgStats damage=%d'%(self.battle_damage)])
    
        wl(battle_log, '--- battle end (%s, %s)'%(a, b))


def wl(l, m):
    l.write('{}\n'.format(m))
    l.flush()


def set_last_event(name, data):
    flash.call('sthub.LastEvent', ['{}: {}'.format(name, utils.jsonEncode(data, ensure_ascii=True))])

    wl(ev_log, '{}\t{}'.format(name, utils.jsonEncode(data, ensure_ascii=True)))


def flash_ready(_m): 
    flash.call('sthub.CallMeBaby', [])


events.onSFMEvent(set_last_event)
events.onFlashReady(flash_ready)
StHub()
