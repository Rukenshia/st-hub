<script>
	import { battles, iteration, darkMode, activeBattle } from './stores';
	import DivisionStatistics from './DivisionStatistics.svelte';
	import ShipsList from './ShipsList.svelte';
	import AppBar from './AppBar.svelte';
	import Battles from './Battles.svelte';
	import axios from 'axios';
  import { onMount } from 'svelte';

  let apiError = false;

	const fetchIntegration = () => {
    return axios.get('http://localhost:1323/iterations/current')
      .then(r => {
        apiError = false;
        return r;
      })
      .catch(err => {
        apiError = true;

        return { data: { ClientVersion: 'n/a', IterationName: 'n/a', Ships: [] }};
      });
	};

	const fetchBattles = () => {
    return axios.get(`http://localhost:1323/iterations/${$iteration.ClientVersion}/${$iteration.IterationName}/battles`)
      .then(r => {
        apiError = false;
        return r;
      })
      .catch(err => {
        apiError = true;

        return { data: [] };
      });
	};

	onMount(async () => {
    const tryConnect = async () => {
      const res = await fetchIntegration();
      $iteration = res.data;

      if (apiError) {
        return false;
      }

      const resBattles = await fetchBattles();
      $battles = resBattles.data === null ? [] : resBattles.data.reverse();
      $activeBattle = $battles.find(b => b.Status === 'active');

      if (apiError) {
        return false;
      }

      setInterval(async () => {
        const resBattles = await fetchBattles();
        $battles = resBattles.data === null ? [] : resBattles.data.reverse();
        $activeBattle = $battles.find(b => b.Status === 'active');
      }, 2500);

      return true;
    };

    if (!await tryConnect()) {
      const loop = () => {
        setTimeout(async () => {
          if (!await tryConnect()) {
            loop();
          }
        }, 2000);
      };

      loop();
    }


	});

	darkMode.subscribe(v => {
		if (v) {
			document.body.classList.add('dark');
		} else {
			document.body.classList.remove('dark');
		}
	});
</script>

<style global lang="scss">
@import '@material/typography/mdc-typography';

body {
	@include mdc-typography-base();
	width: 100%;


	&.dark {
		color: #cecece;
		background-color: #121212;

		.mdc-card {
			background-color: lighten(#121212, 5%);
		}

		svg {
			stroke: #cecece;
			fill: #cecece;
    }

    .mdc-text-field input:disabled {
      color: #cecece;
    }
	}
}

header {
	margin: -8px;
}

footer {
	a {
		color: #cecece;
	}
}
</style>

<AppBar iteration={$iteration} {apiError} />

<DivisionStatistics/>

<ShipsList ships={$iteration.Ships} />

<Battles />

<footer>
	<div class="mdc-typography--subtitle2" style="text-align: center;">
		<div>Made with ❤️ by Rukenshia#4396(Discord), Email: <a href="mailto:svc-sthub@ruken.pw">svc-sthub@ruken.pw</a></div>
		<br />
		<br />
		<div>Icons made by <a href="https://www.freepik.com/" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/"                 title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>
	</div>
</footer>
