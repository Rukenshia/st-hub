<script>
  import Swords from './svg/swords.svelte';
  import CorrectableInput from './CorrectableInput.svelte';
  import { createEventDispatcher } from 'svelte';

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
  class="mdc-card mdc-layout-grid__cell--span-6-desktop mdc-layout-grid__cell--span-8-tablet mdc-layout-grid__cell--span-4-phone">
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
                    class="battle-card__title mdc-typography--headline5">
                    {battle.ShipName}
                  </h2>
                </div>
                <div class="mdc-layout-grid__cell">
                  <div class="mdc-chip-set">
                    {#if battle.Status === 'active'}
                      <div class="mdc-chip">
                        <div class="mdc-chip__text">In Battle</div>
                      </div>
                      {:else if battle.Status === 'abandoned'}
                      <div class="mdc-chip abandoned">
                        <div class="mdc-chip__text">Abandoned</div>
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
      <div class="mdc-layout-grid">
        <div class="mdc-layout-grid__inner">
          <div class="mdc-layout-grid__cell">
            <CorrectableInput on:save={({ detail: {value} }) => updateValue('Damage', value)} noEdit={false} value={battle.Statistics.Damage} label="Damage" helptext="This value does not count fire and flooding damage by default." />
          </div>
          <div class="mdc-layout-grid__cell">
            <CorrectableInput on:save={({ detail: {value} }) => updateValue('Kills', value)} noEdit={false} value={battle.Statistics.Kills} label="Kills" helptext="" />
          </div>
          <div class="mdc-layout-grid__cell">
            <CorrectableInput value={battle.Statistics.InDivision} noEdit={true} label="In Division" helptext="" />
          </div>
        </div>
      </div>
    </div>
    {/if}
  </div>
</div>
