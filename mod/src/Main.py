API_VERSION = 'API_v1.0'
MOD_NAME = 'StHub'

def test(*args, **kwargs):
    print 'STHUB TEST CALLED'

flash.addExternalCallback('stHubTest', test)