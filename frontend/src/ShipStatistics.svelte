<script>
    export let ship;
    export let battles;

    $: averageDamage = Math.round(battles.reduce((p, b) => p + b.Statistics.Damage.Value, 0) / battles.length);
    $: averageKills = Math.round(battles.reduce((p, b) => p + b.Statistics.Kills.Value, 0) / battles.length * 100) / 100;
    $: survivalRate = Math.round(battles.reduce((p, b) => p + b.Statistics.Survived, 0) / battles.length * 100);
    $: winRate = Math.round(battles.reduce((p, b) => p + b.Statistics.Win, 0) / battles.length * 100);
</script>
<!-- 
<div>
    {ship}
    avgDmg: {averageDamage}
    avgKills: {averageKills}
    surv: {survivalRate}%
    win: {winRate}%
</div> -->
