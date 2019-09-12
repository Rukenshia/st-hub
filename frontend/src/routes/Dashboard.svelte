<script>
  import ShipsList from '../ShipsList.svelte';
  import Battles from '../Battles.svelte';

	import { iteration } from '../stores';

</script>

<ShipsList ships={$iteration.Ships} />

<Battles />
