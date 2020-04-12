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
      main_ap: obj.main_ap + battle.Results.Ammo.MainBatteryAP.Damage,
      main_sap: obj.main_sap + battle.Results.Ammo.MainBatterySAP.Damage,
      main_he: obj.main_he + battle.Results.Ammo.MainBatteryHE.Damage,
      secondary: obj.secondary + battle.Results.Ammo.SecondaryHE.Damage + battle.Results.Ammo.SecondarySAP.Damage + battle.Results.Ammo.SecondaryHE.Damage,
      fire: obj.fire + battle.Results.Damage.Fire,
      flood: obj.flood + battle.Results.Damage.Flooding,
      ram: obj.ram + battle.Results.Damage.Ramming,
    };
  }, { total: 0, main_ap: 0, main_sap: 0, main_he: 0, secondary: 0, fire: 0, flood: 0, ram: 0, torpedo: 0 }));

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
          labels: ['AP', 'SAP', 'HE', 'Secondary', 'Torpedo', 'Fire', 'Flooding', 'Ramming'],

        },
      });

    damageBreakdown.subscribe(breakdown => {
      // get percentages
      const main_ap = Math.round(breakdown.main_ap / breakdown.total * 100);
      const main_sap = Math.round(breakdown.main_sap / breakdown.total * 100);
      const main_he = Math.round(breakdown.main_he / breakdown.total * 100);
      const secondary = Math.round(breakdown.secondary / breakdown.total * 100);
      const torpedo = Math.round(breakdown.torpedo / breakdown.total * 100);
      const fire = Math.round(breakdown.fire / breakdown.total * 100);
      const flood = Math.round(breakdown.flood / breakdown.total * 100);
      const ram = Math.round(breakdown.ram / breakdown.total * 100);

      const data = [main_ap, main_sap, main_he, secondary, torpedo, fire, flood, ram];
      console.log(data);
      const colors = [
        'rgb(210, 77, 87)',
        'rgb(260, 57, 57)',
        'rgb(244, 179, 80)',
        'rgb(236, 236, 236)',
        'rgb(103, 128, 159)',
        'rgb(226, 106, 106)',
        'rgb(0, 181, 204)',
        'rgb(36, 37, 42)'
      ];
      const labels = ['AP', 'SAP', 'HE', 'Secondary', 'Torpedo', 'Fire', 'Flooding', 'Ramming'];

      const usedDataIdx = [main_ap, main_sap, main_he, secondary, torpedo, fire, flood, ram].map((v, i) => v > 0 ? i : null).filter(v => v !== null);
      console.log(usedDataIdx);

      chart.data.datasets = [{
        data: usedDataIdx.map(i => data[i]),
        backgroundColor: usedDataIdx.map(i => colors[i]),
      }];
      chart.data.labels = usedDataIdx.map(i => labels[i]);

      chart.update(0);
    });

  });
</script>

<div class="">
  <h1 class="text-xl">Damage breakdown</h1>
  <canvas with="100%" height="100%" id="damageBreakdown"></canvas>
</div>
