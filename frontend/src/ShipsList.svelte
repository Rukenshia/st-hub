<script>
  import { derived } from 'svelte/store';
  import { clientId, activeBattle, battles } from './stores';
  import Swords from './svg/swords.svelte';
  import ShipStatistics from './ShipStatistics.svelte';

  export let ships = [];

  export let showInfo = {};

  function toggle(shipName) {
    showInfo[shipName] = !(showInfo[shipName] || false);
  }

  function showDetails(shipId) {
  }
</script>

<style lang="scss">
@import '@material/card/mdc-card';
@import '@material/layout-grid/mdc-layout-grid';
@import '@material/typography/mdc-typography';

.ship-card__title {
  margin: 0;

  font-weight: normal;
}

.ship-card {
  transition: height .2s;
}

.sl-toggle {
  top: 8px;
  right: 8px;
  transform: translateY(50%);
  position: absolute;
}

.sl-header {
  margin-bottom: 0;
  padding-bottom: 0;
}
</style>

<div class="mdc-layout-grid sl-header">
  <div class="mdc-layout-grid__inner">
    <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-12">
      <h2 class="mdc-typography--headline4">Ships in this test iteration</h2>
    </div>
  </div>
</div>

<div class="mdc-layout-grid">
  <div class="mdc-layout-grid__inner">
    {#each ships as ship}
    <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-4">
      <div class="mdc-card ship-card">
        <div on:click="{() => toggle(ship.Name)}" class="mdc-card__primary-action">
          <div class="mdc-layout-grid mdc-layout-grid--fixed-column-width mdc-layout-grid--align-left">
            <div class="mdc-layout-grid__inner">
              <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-12">
                <h2 class="ship-card__title mdc-typography--headline6">
                  {ship.Name}
                </h2>
              </div>
            </div>
          </div>
          <div class="sl-toggle">
            <i class="material-icons">
              {#if showInfo[ship.Name]}
              arrow_drop_up
              {:else}
              arrow_drop_down
              {/if}
            </i>
          </div>

          {#if showInfo[ship.Name]}
            <ShipStatistics battles={$battles.filter(b => b.ShipID === ship.ID)} />
          {/if}
        </div>

        {#if showInfo[ship.Name]}
          <div class="mb-4 ml-2 mt-4">
            <a on:click={() => showDetails(ship.ID)} href={`/details/${ship.ID}`} class="px-3 py-2 text-teal-500 hover:text-teal-600 font-medium">
              More details
              <!-- <span class="rounded-sm bg-teal-600 px-2 py-1 text-sm text-gray-200">new</span> -->
            </a>
          </div>
        {/if}
      </div>
    </div>
    {/each}
  </div>
</div>
