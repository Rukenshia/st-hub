<script>
  import {
    activeBattle,
    battles,
    darkMode,
    iteration,
    shipInfo
  } from './stores';
  import AppBar from './AppBar.svelte';
  import VersionNotice from './VersionNotice.svelte';
  import axios from 'axios';
  import semver from 'semver';
  import { onMount } from 'svelte';
  import { Router, Link, Route } from 'svelte-routing';

  // Views
  import Dashboard from './routes/Dashboard.svelte';
  import ShipDetails from './routes/ShipDetails.svelte';

  // SVG
  import SpaceSvg from './svg/space.svelte';

  // Utility
  import { getShipInfo } from './wows';

  let apiError = false;
  let version;
  let url = window.location.pathname;

  const availableVersion = '0.7.1';
  let loaded = false;

  const fetchIntegration = () => {
    return axios
      .get(`${ENDPOINT}/iterations/current`)
      .then(r => {
        apiError = false;
        return r;
      })
      .catch(err => {
        apiError = true;

        return {
          data: { ClientVersion: 'n/a', IterationName: 'n/a', Ships: [] }
        };
      });
  };

  const fetchBattles = () => {
    return axios
      .get(
        `${ENDPOINT}/iterations/${$iteration.ClientVersion}/${$iteration.IterationName}/battles`
      )
      .then(r => {
        apiError = false;
        setTimeout(() => (loaded = true), 500);
        return r;
      })
      .catch(err => {
        apiError = true;

        return { data: [] };
      });
  };

  const fetchVersion = () => {
    return axios
      .get(`${ENDPOINT}/version`)
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
        return false;
      }

      const res = await fetchIntegration();
      $iteration = res.data;

      if (apiError) {
        return false;
      }

      $shipInfo = await Promise.all(
        $iteration.Ships.map(({ ID }) => getShipInfo(ID))
      ).then(ships =>
        ships.reduce(
          (prev, cur) => ({
            ...prev,
            [`${cur.ship_id}`]: cur
          }),
          {}
        )
      );

      const resBattles = await fetchBattles();
      $battles = resBattles.data === null ? [] : resBattles.data.reverse();
      $activeBattle = $battles.find(b => b.Status === 'active');

      if (apiError) {
        return false;
      }

      setInterval(async () => {
        const resBattles = await fetchBattles();

        if (resBattles.data === null) {
          $battles = [];
          return;
        }

        // check if there are any new battles
        const reversedBattles = resBattles.data.reverse();

        if (JSON.stringify(reversedBattles) === JSON.stringify($battles)) {
          // Don't set same data again to prevent UI reloads
          return;
        }

        $battles = reversedBattles;
        $activeBattle = $battles.find(b => b.Status === 'active');
      }, 2500);

      return true;
    };

    if (!(await tryConnect())) {
      const loop = () => {
        setTimeout(async () => {
          if (!(await tryConnect())) {
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
  @import 'tailwindcss/base';
  @import 'tailwindcss/components';
  @import 'tailwindcss/utilities';

  body {
    width: 100%;

    &.dark {
      @apply bg-cool-gray-800 text-gray-50;

      .mdc-card,
      .mdc-card.ship-card {
        @apply bg-cool-gray-900;
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

      .input-dialog.mdc-dialog .mdc-dialog__surface {
        @apply bg-cool-gray-900;
        color: #cecece;

        .mdc-dialog__title,
        .mdc-dialog__content {
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

{#if version && semver.gt(availableVersion, version)}
  <VersionNotice {availableVersion} {version} />
{/if}

{#if !loaded}
  <div class="mx-auto mt-32 w-1/2 mb-64">
    <h1 class="text-2xl text-center pb-4">Loading your data</h1>
    <SpaceSvg />
  </div>
{:else}
  <Router {url}>
    <Route path="/" component={Dashboard} />
    <Route path="/details/:id" component={ShipDetails} />
  </Router>
{/if}

<footer>
  <div class="text-sm" style="text-align: center;">
    <div>
      Made with ❤️ by Rukenshia#4396(Discord), Email:
      <a href="mailto:svc-sthub@ruken.pw">svc-sthub@ruken.pw</a>
    </div>
    <br />
    <br />
    <div>
      Icons made by
      <a href="https://www.freepik.com/" title="Freepik">Freepik</a>
      from
      <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a>
      is licensed by
      <a
        href="http://creativecommons.org/licenses/by/3.0/"
        title="Creative Commons BY 3.0"
        target="_blank">
        CC 3.0 BY
      </a>
    </div>
  </div>
</footer>
