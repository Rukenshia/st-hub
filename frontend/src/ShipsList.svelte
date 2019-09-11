<script>
  import { derived } from 'svelte/store';
  import { clientId, activeBattle, battles } from './stores';
  import ShipCard from './ShipCard.svelte';

  export let ships = [];
</script>

<style lang="scss">
@import '@material/layout-grid/mdc-layout-grid';
@import '@material/typography/mdc-typography';

.sl-header {
  margin-bottom: 0;
  padding-bottom: 0;
}
</style>

<div class="mdc-layout-grid sl-header">
  <div class="mdc-layout-grid__inner">
    <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-12">
      <div class="flex">
        <h2 class="mdc-typography--headline4">Ships in this test iteration</h2>
        <!--<a href="/ships/compare" class="ml-4 rounded-sm border border-teal-600 px-3 py-2 text-teal-600 hover:border-teal-400 hover:text-teal-400 uppercase">Compare all testships</a>-->
      </div>
    </div>
  </div>
</div>

<div class="mdc-layout-grid">
  <div class="mdc-layout-grid__inner">
    {#each ships as ship}
    <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-4">
      <ShipCard {ship} battles={derived(battles, nb => nb.filter(b => b.ShipID === ship.ID))} />
    </div>
    {/each}
  </div>
</div>
