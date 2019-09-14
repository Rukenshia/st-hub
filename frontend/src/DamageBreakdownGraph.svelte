<script>
  import { onMount } from 'svelte';
  import { derived } from 'svelte/store';

  import Chart from 'chart.js';

  export let battles;

  const damageBreakdown = derived(battles, battles => battles.reduce((obj, battle) => {
    if (!battle.Results) {
      return obj;
    }

    return {
      total: obj.total + battle.Results.Damage.Sum,
      torpedo: obj.torpedo + battle.Results.Ammo.Torpedo.Damage,
      main_ap: obj.main_ap + battle.Results.Ammo.MainBatteryAP.Damage + battle.Results.Ammo.MainBatterySAP.Damage,
      main_he: obj.main_he + battle.Results.Ammo.MainBatteryHE.Damage,
      secondary: obj.secondary + battle.Results.Ammo.SecondaryHE.Damage + battle.Results.Ammo.SecondarySAP.Damage + battle.Results.Ammo.SecondaryHE.Damage,
      fire: obj.fire + battle.Results.Damage.Fire,
      flood: obj.flood + battle.Results.Damage.Flooding,
      ram: obj.ram + battle.Results.Damage.Ramming,
    };
  }, { total: 0, main_ap: 0, main_he: 0, secondary: 0, fire: 0, flood: 0, ram: 0, torpedo: 0 }));

  onMount(() => {
    const el = document.getElementById('damageBreakdown');

    const chart =  new Chart(el, {
        type: 'pie',
        options: {
          tooltips: {
            callbacks: {
              label: function(tooltipItem, data) {
                const dataset = data.datasets[tooltipItem.datasetIndex];
                const currentValue = dataset.data[tooltipItem.index];
                return `${currentValue}%`;
              },
              title: function(tooltipItem, data) {
                return data.labels[tooltipItem[0].index];
              }
            }
          },
        },
        data: {
          labels: ['Main battery AP', 'Main battery HE', 'Secondary battery', 'Torpedo', 'Fire', 'Flooding', 'Ramming'],

        },
      });

    damageBreakdown.subscribe(breakdown => {
      // get percentages
      const main_ap = Math.round(breakdown.main_ap / breakdown.total * 100);
      const main_he = Math.round(breakdown.main_he / breakdown.total * 100);
      const secondary = Math.round(breakdown.secondary / breakdown.total * 100);
      const torpedo = Math.round(breakdown.torpedo / breakdown.total * 100);
      const fire = Math.round(breakdown.fire / breakdown.total * 100);
      const flood = Math.round(breakdown.flood / breakdown.total * 100);
      const ram = Math.round(breakdown.ram / breakdown.total * 100);

      chart.data.datasets = [{
        data: [main_ap, main_he, secondary, torpedo, fire, flood, ram],
        backgroundColor: [
          'rgb(210, 77, 87)',
          'rgb(244, 179, 80)',
          'rgb(236, 236, 236)',
          'rgb(103, 128, 159)',
          'rgb(226, 106, 106)',
          'rgb(0, 181, 204)',
          'rgb(36, 37, 42)'
        ],
      }];

      chart.update(0);
    });

  });
</script>

<div class="border rounded-sm border-gray-900 bg-gray-900 p-4">
  <h1 class="text-xl">Damage breakdown</h1>
  <canvas with="100%" height="100%" id="damageBreakdown"></canvas>
</div>
