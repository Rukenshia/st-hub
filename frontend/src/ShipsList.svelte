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
    <div class="mdc-layout-grid__cell">
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
            <ShipStatistics ship={ship.Name} battles={$battles.filter(b => b.ShipID === ship.ID)} />
          {/if}
        </div>
      </div>
    </div>
    {/each}
  </div>
</div>
