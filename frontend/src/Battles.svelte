<script>
  import { onMount } from 'svelte';
  import { derived, writable } from 'svelte/store';
  import { battles, iteration } from './stores';
  import Battle from './Battle.svelte';
  import axios from 'axios';

  const shipNames = derived(battles, $battles => [
    ...new Set($battles.map(b => b.ShipName))
  ]);

  const selectedShip = writable('all');

  const filteredBattles = derived(
    [battles, selectedShip],
    ([$battles, $selectedShip]) =>
      $battles.filter(
        b => b.ShipName === $selectedShip || $selectedShip === 'all'
      )
  );

  function updateField({detail: {battle, field, value}}) {
    const idx = $battles.findIndex(b => b.ID === battle.ID);
    if (idx === -1) {
      return;
    }

    battle.Statistics[field].Corrected = parseInt(value, 10);
    $battles[idx].Statistics[field].Corrected = parseInt(value, 10);

    return axios.put(`${ENDPOINT}/iterations/${$iteration.ClientVersion}/battles/${battle.ID}`, battle)
      .catch(err => {
        alert(`Please go talk to rukenshia: ${err}`)
      });
  }
</script>

<style global lang="scss">
  @import '@material/card/mdc-card';
  @import '@material/chips/mdc-chips';
  @import '@material/layout-grid/mdc-layout-grid';

  body {
    .battle-card {
      .battle-card__primary {
        padding: 1rem;

        .battle-card__title {
          margin: 0;
        }

        .mdc-layout-grid {
          padding: 0;
        }

        .mdc-chip-set {
          padding-left: 0;
          padding-top: 0;
          .mdc-chip {
            @include mdc-chip-height(24px);
            @apply bg-gray-700;
            color: #cecece;
            font-size: 12px;


            &.abandoned {
              @apply bg-yellow-300;
              color: #121212;
            }
          }
        }
      }
    }
  }

  select {
    -webkit-appearance: none;
    -moz-appearance: none;
    text-indent: 1px;
    text-overflow: '';
  }

  .form-select.ship-filter {
    @apply border-gray-600 bg-gray-900 text-gray-50;
  }

  .battles-header {
    margin-bottom: 0;
    padding-bottom: 0;
  }
</style>

<div class="mdc-layout-grid mdc-layout-grid--align-left battles-header">
  <div class="mdc-layout-grid__inner">
    <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-2 mdc-layout-grid__cell--align-bottom">
      <h2 class="text-3xl">Battles</h2>
    </div>
    <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-2 mdc-layout-grid__cell--align-middle">
      <div>
        <div class="mt-1 relative rounded-md shadow-sm">
          <select class="form-select ship-filter block w-full sm:text-sm sm:leading-5" bind:value={$selectedShip}>
            <option value="all">All ships</option>
            {$shipNames}
            {#each $shipNames as name}
              <option value={name}>{name}</option>
            {/each}
          </select>
        </div>
      </div>
    </div>
  </div>
</div>

<div class="mdc-layout-grid battles">
  <div class="mdc-layout-grid__inner">
    {#if $filteredBattles.length === 0}
    <div class="mdc-layout-grid__cell" style="padding: 16px">
      <span class="text-md">No battles played</span>
    </div>
    {:else}
    {#each $filteredBattles as battle}
      <Battle {battle} on:update={updateField} />
      <div class="mdc-layout-grid__cell--span-2-desktop mdc-layout-grid__cell--span-1-tablet mdc-layout-grid__cell--span-4-phone"></div>
    {/each}
    {/if}
  </div>
</div>
