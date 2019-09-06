<script>
  import ShipsList from '../ShipsList.svelte';
  import DivisionStatistics from '../DivisionStatistics.svelte';
  import Battles from '../Battles.svelte';

	import { iteration } from '../stores';
</script>

<DivisionStatistics/>

<ShipsList ships={$iteration.Ships} />

<Battles />