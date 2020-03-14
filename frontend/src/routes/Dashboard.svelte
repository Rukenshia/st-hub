<script>
  import ShipsList from '../ShipsList.svelte';
  import Battles from '../Battles.svelte';
  export let location;

	import { iteration } from '../stores';

</script>

<div class="mx-full bg-orange-600 text-orange-100 p-2">
  Does anyone still use this? If yes, please just send me an empty email
  or write in Discord

  <a href="mailto:svc-sthub@ruken.pw" class="border-none bg-gray-700 rounded-sm hover:bg-gray-800 px-2 py-1">Mail me</a>
</div>
<ShipsList ships={$iteration.Ships} />

<Battles />
