<script>
  import ShipsList from '../ShipsList.svelte';
  import Battles from '../Battles.svelte';
  export let location;

  import { iteration } from '../stores';
</script>

<div class="mt-4 px-4">
  <div class="rounded-md bg-yellow-400 p-4 mx-auto max-w-2xl">
    <div class="flex">
      <div class="ml-3">
        <h3 class="text-sm leading-5 font-medium text-yellow-50">
          Important Information
        </h3>
        <div class="mt-2 text-sm leading-5 text-yellow-50">
          <p>
            Due to a bug in 0.9.5.0, it is currently not possible to retrieve
            battle results. I hope that Wargaming will address the issue in the
            next patch. Your battles will still be recorded, but the data will
            be inaccurate. You can edit battle data at the bottom of the page.
          </p>
        </div>
      </div>
    </div>
  </div>
</div>

<ShipsList ships={$iteration.Ships} />

<Battles />
