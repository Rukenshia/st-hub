<script>
    import { derived, writable } from 'svelte/store';
    import { battles } from './stores';
    import ShipStatistics from './ShipStatistics.svelte';

    const shipNames = derived(battles,
        $battles => [...new Set($battles.map(b => b.ShipName))]);

    const selectedShip = writable('all');

    const filteredBattles = derived([battles, selectedShip],
        ([b, s]) => b.filter(b => b.ShipName  === s || s === 'all'));
</script>

<style global lang="scss">
@import '@material/card/mdc-card';
@import '@material/chips/mdc-chips';
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

                    &.loss {
                        @include mdc-chip-fill-color(#fedede);
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
                &.loss {
                    @include mdc-chip-fill-color(#ff574a);
                    color: rgba(0, 0, 0, 0.87);
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
                                    {battle.ShipName}
                                </h2>
                            </div>
                            <div class="mdc-layout-grid__cell">
                                <div class="mdc-chip-set">
                                    {#if battle.Status === 'active'}
                                    <div class="mdc-chip">
                                        <div class="mdc-chip__text">In Battle</div>
                                    </div>
                                    {:else}
                                        <div class="mdc-chip" class:loss={!battle.Statistics.Win}>
                                            <div class="mdc-chip__text">{battle.Statistics.Win ? 'Win' : 'Loss'}</div>
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
                    
                    

                    <p>
                        Damage (raw): {battle.Statistics.Damage.Value}
                    </p>
                </div>
            </div>
        </div>
        {/each}
    </div>
</div>