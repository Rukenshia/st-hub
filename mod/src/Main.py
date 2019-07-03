API_VERSION = 'API_v1.0'
MOD_NAME = 'StHub'

devmenu.enable()

ev_log = open('events.log', 'w')
battle_log = open('battle.log', 'w')


class StHub:
    def __init__(self):
        self.battle_damage = 0
        self.in_division = False

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
        if flags & 0b1:
            flash.call('sthub.LastEvent', [
                       'StHub.Shell received damage %s' % (damage)])
            wl(battle_log, 'shell: %s %s %s %s %s %s %s %s %s %s' % (
                victim, shooter, ammo, mat, shoot, '{0:b}'.format(flags), damage, pos, yaw, hlinfo))
        else:
            flash.call('sthub.LastEvent', [
                       'StHub.Shell dealt damage %s' % (damage)])
            wl(battle_log, 'shell: %s %s %s %s %s %s %s %s %s %s' % (
                victim, shooter, ammo, mat, shoot, '{0:b}'.format(flags), damage, pos, yaw, hlinfo))

            self.battle_damage = self.battle_damage + damage

        flash.call('sthub.LastEvent', [
                   'StHub.AmmoInfo %s' % (battle.getAmmoParams(ammo))])
        wl(battle_log, 'ammo_info %s' % (battle.getAmmoParams(ammo)))

    def battle_start(self):
        wl(battle_log, '--- battle start')
        flash.call('sthub.LastEvent', ['StHub.BattleStart'])

        selfInfo = battle.getSelfPlayerInfo()

        wl(open('battle.self.log', 'w'), '%s' % (selfInfo))
        wl(open('battle.players.log', 'w'), '%s' % (battle.getPlayersInfo()))

        wl(battle_log, '> shipId: %s' % (selfInfo.shipId))
        wl(battle_log, '> shipInfo: %s' % (selfInfo.shipInfo))
        wl(battle_log, '> shipName: %s' % (selfInfo.shipInfo.name))

        data = {
            'ShipID': selfInfo.shipId,
            'ShipName': selfInfo.shipInfo.name,
            'InDivision': self.in_division,
        }

        wl(battle_log, '! allowed_urls: %s' % (web.getAllowedUrls()))

        res = web.openUrl(
            'http://localhost:1323/iteration/current/battles', data=utils.jsonEncode(data))
        wl(battle_log, '> server res: %s' % (res))
        wl(battle_log, '> server: %s' % (res.read()))

        self.battle_damage = 0

    def battle_quit(self, _m):
        flash.call('sthub.LastEvent', ['StHub.BattleQuit'])
        wl(battle_log, '--- battle quit')

        wl(open('battle.self_end.log', 'w'), '%s' %
           (battle.getSelfPlayerInfo()))
        wl(open('battle.players_end.log', 'w'), '%s' %
           (battle.getPlayersInfo()))

    def battle_end(self, a, b):
        flash.call('sthub.LastEvent', ['StHub.BattleEnd %s %s' % (a, b)])
        flash.call('sthub.LastEvent', [
                   'StHub.AvgStats damage=%d' % (self.battle_damage)])

        wl(battle_log, '--- battle end (%s, %s)' % (a, b))


def wl(l, m):
    l.write('{}\n'.format(m))
    l.flush()


def set_last_event(name, data):
    flash.call('sthub.LastEvent', ['{}: {}'.format(
        name, utils.jsonEncode(data, ensure_ascii=True))])

    wl(ev_log, '{}\t{}'.format(name, utils.jsonEncode(data, ensure_ascii=True)))


def flash_ready(_m):
    flash.call('sthub.CallMeBaby', [])


events.onSFMEvent(set_last_event)
events.onFlashReady(flash_ready)
StHub()
