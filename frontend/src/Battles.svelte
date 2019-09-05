<script>
  import { onMount } from 'svelte';
  import { derived, writable } from 'svelte/store';
  import { battles, iteration } from './stores';
  import { MDCSelect } from '@material/select';
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

  onMount(() => {
    new MDCSelect(document.querySelector('.mdc-select'));
  });

  function updateField({detail: {battle, field, value}}) {
    const idx = $battles.findIndex(b => b.ID === battle.ID);
    if (idx === -1) {
      return;
    }

    battle.Statistics[field].Corrected = parseInt(value, 10);
    $battles[idx].Statistics[field].Corrected = parseInt(value, 10);

    return axios.put(`http://localhost:1323/iterations/${$iteration.ClientVersion}/${$iteration.IterationName}/battles/${battle.ID}`, battle)
      .catch(err => {
        alert(`Please go talk to rukenshia: ${err}`)
      });
  }
</script>

<style global lang="scss">
  @import '@material/card/mdc-card';
  @import '@material/chips/mdc-chips';
  @import '@material/select/mdc-select';
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
            font-size: 12px;

            &.abandoned {
              @apply bg-yellow-500;
              color: #121212;
            }
          }
        }
      }
    }

    &.dark {
      .battle-card .battle-card__primary .mdc-chip-set {
        .mdc-chip {
          @include mdc-chip-fill-color(lighten(#121212, 11%));
          color: #cecece;

          &.abandoned {
            @apply bg-yellow-500;
            color: #121212;
          }
        }
      }

      .mdc-select {
        @include mdc-select-container-fill-color(lighten(#121212, 5%));
        @include mdc-select-focused-label-color(
          lighten(rgba(98, 0, 238, 0.87), 25%)
        );
        @include mdc-select-focused-bottom-line-color(
          lighten(rgba(98, 0, 238, 0.87), 25%)
        );

        .mdc-floating-label,
        .mdc-select__native-control {
          color: #cecece;

          option {
            background-color: lighten(#121212, 5%);
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

  .battles-header {
    margin-bottom: 0;
    padding-bottom: 0;
  }
</style>

<div class="mdc-layout-grid mdc-layout-grid--align-left battles-header">
  <div class="mdc-layout-grid__inner">
    <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-2 mdc-layout-grid__cell--align-bottom">
      <h2 class="mdc-typography--headline4">Battles</h2>
    </div>
    <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-2 mdc-layout-grid__cell--align-middle">
      <div class="mdc-select">
        <i class="mdc-select__dropdown-icon" />
        <select class="mdc-select__native-control" bind:value={$selectedShip}>
          <option value="all">all</option>
           {$shipNames}
          {#each $shipNames as name}
            <option value={name}>{name}</option>
          {/each}
        </select>
        <label class="mdc-floating-label">Pick a ship</label>
        <div class="mdc-line-ripple" />
      </div>
    </div>
  </div>
</div>

<div class="mdc-layout-grid battles">
  <div class="mdc-layout-grid__inner">
    {#if $filteredBattles.length === 0}
    <div class="mdc-layout-grid__cell" style="padding: 16px">
      <span class="mdc-typography--subtitle1">No battles played</span>
    </div>
    {:else}
    {#each $filteredBattles as battle}
      <Battle {battle} on:update={updateField} />
      <div class="mdc-layout-grid__cell--span-2-desktop mdc-layout-grid__cell--span-1-tablet mdc-layout-grid__cell--span-4-phone"></div>
    {/each}
    {/if}
  </div>
</div>
