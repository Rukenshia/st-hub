<script>
  import Swords from './svg/swords.svelte';
  import { createEventDispatcher } from 'svelte';
  import ResultWithFallback from './ResultWithFallback.svelte';
  import Result from './Result.svelte';

	const dispatch = createEventDispatcher();

  export let battle;

  let toggle = false;

  function updateValue(field, value) {
    dispatch('update', { battle, field, value });
  }
</script>

<style lang="scss">
.toggle {
  top: 8px;
  right: 8px;
  transform: translateY(50%);
  position: absolute;
}

</style>

<div
  class="mdc-card mdc-layout-grid__cell--span-10-desktop mdc-layout-grid__cell--span-8-tablet mdc-layout-grid__cell--span-4-phone">
  <div class="mdc-card__primary-action battle-card">
    <div class="battle-card__primary" on:click="{() => toggle = !toggle}">
      <div class="mdc-layout-grid" >
        <div class="mdc-layout-grid__inner">
          <div
            class="mdc-layout-grid__cell
            mdc-layout-grid__cell--span-1-desktop
            mdc-layout-grid__cell--span-1-tablet
            mdc-layout-grid__cell--span-2-phone">
            {#if battle.Status === 'active'}
            <Swords />
            {/if}
          </div>
          <div
            class="mdc-layout-grid__cell
            mdc-layout-grid__cell--span-11-desktop
            mdc-layout-grid__cell--span-7-tablet
            mdc-layout-grid__cell--span-2-phone">
            <div class="mdc-layout-grid">
              <div class="mdc-layout-grid__inner">
                <div class="mdc-layout-grid__cell">
                  <h2
                    class="battle-card__title text-lg">
                    {battle.ShipName}
                  </h2>
                </div>
                <div class="mdc-layout-grid__cell">
                  <div class="mdc-chip-set">
                    {#if !battle.Results && battle.Status !== 'active'}
                      <div class="mdc-chip abandoned">
                        <div class="mdc-chip__text">Result screen missing</div>
                      </div>
                    {/if}
                    {#if battle.Status === 'active'}
                      <div class="mdc-chip">
                        <div class="mdc-chip__text">In Battle</div>
                      </div>
                      {:else if battle.Status === 'abandoned'}
                      <div class="mdc-chip">
                        <div class="mdc-chip__text">Left</div>
                      </div>
                      {:else}
                      <div
                        class="mdc-chip"
                        class:loss={!battle.Statistics.Win}>
                        <div class="mdc-chip__text">
                          {battle.Statistics.Win ? 'Win' : 'Loss'}
                        </div>
                      </div>
                      {#if battle.Statistics.InDivision.Value}
                      <div class="mdc-chip">
                        <div class="mdc-chip__text">Division</div>
                      </div>
                      {/if}
                    {/if}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="toggle">
      <i class="material-icons">
        {#if toggle}
        arrow_drop_up
        {:else}
        arrow_drop_down
        {/if}
      </i>
    </div>

    {#if toggle}
    <div class="battle-card__content">
      {#if !battle.Results}
      <div class="mdc-layout-grid">
        <div class="mdc-layout-grid__inner">
          <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-12 bg-yellow-400 text-yellow-900 p-4 rounded-sm">
            This battle does not contain accurate data because you left it before the game ended. To collect all data, please stay until you see the "Results screen" of a battle.
          </div>
        </div>
      </div>
      {/if}
      <div class="mdc-layout-grid">
        <div class="mdc-layout-grid__inner">
          {#if battle.Results}
          <div class="mdc-layout-grid__cell">
            <Result label="Base Experience" value={battle.Results.Economics.BaseExp} />
          </div>
          <div class="mdc-layout-grid__cell">
            <Result label="Credits" value={battle.Results.Economics.Credits} />
          </div>
          {/if}
          <div class="mdc-layout-grid__cell">
            <ResultWithFallback on:save={({ detail: {value} }) => updateValue('Damage', value)} label="Damage" value={battle.Results ? battle.Results.Damage.Sum : undefined} correctable={battle.Statistics.Damage} />
          </div>
          <div class="mdc-layout-grid__cell">
            <ResultWithFallback on:save={({ detail: {value} }) => updateValue('Kills', value)} label="Kills" value={undefined} correctable={battle.Statistics.Kills} />
          </div>
          <div class="mdc-layout-grid__cell">
            <ResultWithFallback on:save={({ detail: {value} }) => updateValue('InDivision', value)} label="In Division" value={undefined} correctable={battle.Statistics.InDivision} />
          </div>
        </div>
      </div>
    </div>
    {/if}
  </div>
</div>
