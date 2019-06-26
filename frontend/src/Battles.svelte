<script>
    import { derived, writable } from 'svelte/store';
    import { battles } from './stores';
    import ShipStatistics from './ShipStatistics.svelte';

    const shipNames = derived(battles,
        $battles => [...new Set($battles.map(b => b.statistics.ship))]);

    const selectedShip = writable('all');

    const filteredBattles = derived([battles, selectedShip],
        ([b, s]) => b.filter(b => b.statistics.ship === s || s === 'all'));
</script>

<style lang="scss">
@import '@material/card/mdc-card';
@import '@material/chips/mdc-chips';
@import '@material/layout-grid/mdc-layout-grid';

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

                &.loss {
                    @include mdc-chip-fill-color(#fedede);
                }
            }
        }
    }
}
</style>

<select bind:value={$selectedShip}>
    <option value='all'>all</option>
    {$shipNames}
    {#each $shipNames as name}
        <option value={name}>{name}</option>
    {/each}
</select>

{#if $selectedShip !== 'all'}
<ShipStatistics ship={$selectedShip} battles={$filteredBattles} />
{/if}

<div class="mdc-layout-grid">
    <div class="mdc-layout-grid__inner">
        {#each $filteredBattles as battle}
        <div class="mdc-card mdc-layout-grid__cell mdc-layout-grid__cell--span-12">
            <div class="mdc-card__primary-action battle-card">
                <div class="battle-card__primary">
                    <div class="mdc-layout-grid">
                        <div class="mdc-layout-grid__inner">
                            <div class="mdc-layout-grid__cell mdc-layout-grid__cell">
                                <h2 class="battle-card__title mdc-typography--headline6">
                                    {battle.statistics.ship}
                                </h2>
                            </div>
                            <div class="mdc-layout-grid__cell">
                                <div class="mdc-chip-set">
                                    <div class="mdc-chip" class:loss={!battle.statistics.win}>
                                        <div class="mdc-chip__text">{battle.statistics.win ? 'Win' : 'Loss'}</div>
                                    </div>
                                    {#if battle.statistics.division}
                                    <div class="mdc-chip">
                                        <div class="mdc-chip__text">Division</div>
                                    </div>
                                    {/if}
                                </div>
                            </div>
                        </div>
                    </div>
                    
                    

                    <p>
                        Damage: {battle.statistics.damage}
                    </p>
                </div>
            </div>
        </div>
        {/each}
    </div>
</div>