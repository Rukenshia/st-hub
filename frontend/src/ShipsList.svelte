<script>
  import { clientId, activeBattle } from './stores';
  import Swords from './svg/swords.svelte';

  export let ships = [];
</script>

<style lang="scss">
@import '@material/card/mdc-card';
@import '@material/layout-grid/mdc-layout-grid';
@import '@material/typography/mdc-typography';

.ship-card__title {
  margin: 0;

  font-weight: normal;
}
</style>

<h2 class="mdc-typography--headline4">Ships in this test iteration</h2>

<div class="mdc-layout-grid">
  <div class="mdc-layout-grid__inner">
    {#each ships as ship}
    <div class="mdc-layout-grid__cell mdc-card">
      <div class="mdc-card__primary-action">
        <div class="mdc-layout-grid mdc-layout-grid--align-left">
          <div class="mdc-layout-grid__inner">
            {#if $activeBattle && $activeBattle.Ship === ship.ID}
            <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-2-desktop mdc-layout-grid__cell--span-2-tablet mdc-layout-grid__cell--span-1-phone">
                <Swords />
            </div>
            {/if}
            <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-10-desktop mdc-layout-grid__cell--span-6-tablet mdc-layout-grid__cell--span-2-phone">
              <h2 class="ship-card__title mdc-typography--headline6">
                {ship.Name}
              </h2>
            </div>
          </div>
        </div>
      </div>
    </div>
    {/each}
  </div>
</div>
