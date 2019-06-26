<script>
    import { derived } from 'svelte/store';
    import { battles } from './stores';
    import DivisionRate from './DivisionRate.svelte';

    let ships = derived(
        battles,
        $battles => Object.values($battles.reduce((p, c) => {
                if (p[c.statistics.ship]) {
                    p[c.statistics.ship].battles++;

                    if (c.statistics.division) {
                        p[c.statistics.ship].division++;
                    }

                    p[c.statistics.ship].rate = Math.round(p[c.statistics.ship].division / p[c.statistics.ship].battles * 100);
                } else {
                    p[c.statistics.ship] = {
                        name: c.statistics.ship,
                        battles: 1,
                        division: c.statistics.division ? 1 : 0,
                        rate: Math.round(c.statistics.division / 1 * 100),
                    };
                };

                return p;
            }, {})).sort((a, b) => {
                    if (a.rate < b.rate) { return 1; }
                    else if (a.rate == b.rate) { return 0; }
                    else { return -1; }
                })
        );
</script>

<style lang="scss">
@import '@material/layout-grid/mdc-layout-grid';
@import '@material/card/mdc-card';

.div-card {
    .div-card__primary {
        padding: 1rem;

        .div-card__title {
            margin: 0;
        }

        .div-card__title-icon {
            height: 3rem;
        }

        table {
            padding-top: 8px;

            th {
                text-align: left;
            }
        }
    }
}
</style>

<div class="mdc-layout-grid">
    <div class="mdc-layout-grid__inner">
        <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-12 mdc-card div-card">
            <div class="div-card__primary mdc-card__primary-action">
                <div class="mdc-layout-grid">
                    <div class="mdc-layout-grid__inner">
                        <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-1">
                            <img alt="division" src="/team.svg" class="div-card__title-icon" />
                        </div>
                        <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-11">
                            <h2 class="div-card__title mdc-typography--headline6">
                                Division statistics
                            </h2>

                            <table>
                                <thead>
                                    <tr>
                                        <th>Ship</th>
                                        <th>Ratio</th>
                                        <th>Div / Total</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {#each $ships as {name, battles, division, rate}}
                                    <DivisionRate  {name} {battles} {division} {rate}></DivisionRate>
                                    {/each}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
