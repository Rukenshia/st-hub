<script>
  import { derived } from 'svelte/store';
  import { battles } from './stores';
  import DivisionRate from './DivisionRate.svelte';
  import TeamSvg from './svg/team.svelte';

  let ships = derived(battles, $battles =>
    Object.values(
      $battles.reduce((p, c) => {
        if (p[c.ShipID]) {
          p[c.ShipID].battles++;

          if (c.Statistics.InDivision.Value) {
            p[c.ShipID].division++;
          }

          p[c.ShipID].rate = Math.round(
            (p[c.ShipID].division / p[c.ShipID].battles) * 100
          );
        } else {
          p[c.ShipID] = {
            name: c.ShipName,
            battles: 1,
            division: c.Statistics.InDivision.Value ? 1 : 0,
            rate: c.Statistics.InDivision.Value ? 100 : 0,
          };
        }

        return p;
      }, {})
    ).sort((a, b) => {
      if (a.rate < b.rate) {
        return 1;
      } else if (a.rate == b.rate) {
        return 0;
      } else {
        return -1;
      }
    })
  );
</script>

<style lang="scss">
  @import "@material/layout-grid/mdc-layout-grid";
  @import "@material/card/mdc-card";

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
    <div
      class="mdc-layout-grid__cell mdc-layout-grid__cell--span-12 mdc-card
      div-card">
      <div class="div-card__primary mdc-card__primary-action">
        <div class="mdc-layout-grid mdc-layout-grid--align-left">
          <div class="mdc-layout-grid__inner">
            <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-1">
              <TeamSvg />
            </div>
            <div class="mdc-layout-grid__cell mdc-layout-grid__cell--span-11">
              <h2 class="div-card__title mdc-typography--headline6">
                Division Statistics
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
                  {#each $ships as { name, battles, division, rate }}
                    <DivisionRate {name} {battles} {division} {rate} />
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
