<script>
	import { activeBattle, battles, darkMode, iteration } from './stores';
	import AppBar from './AppBar.svelte';
	import VersionNotice from './VersionNotice.svelte';
	import axios from 'axios';
  import { onMount } from 'svelte';
  import { Router, Link, Route } from "svelte-routing";

  // Views
	import Dashboard from './routes/Dashboard.svelte';
	import ShipDetails from './routes/ShipDetails.svelte';

  let apiError = false;
  let version;
  let url = window.location.pathname;

  const availableVersion = '0.5.0';

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

	const fetchVersion = () => {
    return axios.get(`http://localhost:1323/version`)
      .then(r => {
        apiError = false;
        return r.data;
      })
      .catch(err => {
        apiError = true;

        return null;
      });
	};

	onMount(async () => {
    const tryConnect = async () => {
      version = await fetchVersion();

      if (apiError) {
        return false
      }

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
@tailwind base;
@tailwind utilities;
@tailwind components;

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

    #field-helper-text {
      color: #cecece;
    }

    .input-dialog.mdc-dialog .mdc-dialog__surface  {
      background-color: lighten(#121212, 15%);
      color: #cecece;

      .mdc-dialog__title,.mdc-dialog__content {
        color: #cecece;
      }
    }
	}
}

footer {
	a {
		color: #cecece;
	}
}
</style>

<AppBar iteration={$iteration} {version} {apiError} />

{#if version && version !== availableVersion}
<VersionNotice {availableVersion} {version} />
{/if}

<Router {url}>
  <Route path="/" component="{Dashboard}"></Route>
  <Route path="/details/:id" component="{ShipDetails}"></Route>
</Router>

<footer>
	<div class="mdc-typography--subtitle2" style="text-align: center;">
		<div>Made with ❤️ by Rukenshia#4396(Discord), Email: <a href="mailto:svc-sthub@ruken.pw">svc-sthub@ruken.pw</a></div>
		<br />
		<br />
		<div>Icons made by <a href="https://www.freepik.com/" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/"                 title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>
	</div>
</footer>