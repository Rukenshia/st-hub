<script>
  import ShipBasicAverageStats from './ShipBasicAverageStats.svelte';
  import TeamSvg from './svg/team.svelte';
  import DivisionRate from './DivisionRate.svelte';

  export let battles;
  export let ship;
  let showInfo;
</script>

<style lang="scss">
@import '@material/card/mdc-card';
@import '@material/layout-grid/mdc-layout-grid';

.ship-card__title {
  margin: 0;

  font-weight: normal;
}

.ship-card {
  transition: height .2s;

  .ship-header {
    @apply pt-4;
  }
}

.sl-toggle {
  top: 8px;
  right: 8px;
  transform: translateY(50%);
  position: absolute;
}

</style>


<div class="mdc-card ship-card">
  <div on:click="{() => showInfo = !showInfo}" class="mdc-card__primary-action">
    <div class="ship-header mdc-layout-grid mdc-layout-grid--fixed-column-width mdc-layout-grid--align-left">
      <div class="mdc-layout-grid__inner">
        <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-12">
          <h2 class="ship-card__title mdc-typography--headline6">
            {ship.Name}
          </h2>
          <div class="-mt-1"><DivisionRate {battles} /></div>
        </div>
      </div>
    </div>
    <div class="sl-toggle">
      <i class="material-icons">
        {#if showInfo}
        arrow_drop_up
        {:else}
        arrow_drop_down
        {/if}
      </i>
    </div>

    {#if showInfo}
      <ShipBasicAverageStats {$battles} />
    {/if}
  </div>

  {#if showInfo}
    <div class="mb-4 ml-2 mt-4">
      <a href={`/details/${ship.ID}`} class="px-3 py-2 text-teal-500 hover:text-teal-600 font-medium">
        More details
        <!-- <span class="rounded-sm bg-teal-600 px-2 py-1 text-sm text-gray-200">new</span> -->
      </a>
    </div>
  {/if}
</div>
