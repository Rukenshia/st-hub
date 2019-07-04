<script>
  import { clientId } from './stores';

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

.mdc-card__primary-action {
  padding: 8px;
}
</style>

<div class="mdc-layout-grid">
  <div class="mdc-layout-grid__inner">
    {#each ships as ship}
    <div class="mdc-layout-grid__cell mdc-card">
      <div class="mdc-card__primary-action">
        <h2 class="ship-card__title mdc-typography--headline6">
          {ship.Name}
        </h2>
      </div>
    </div>
    {/each}
  </div>
</div>
